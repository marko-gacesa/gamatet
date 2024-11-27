// Copyright (c) 2020-2024 by Marko Gaćeša

package core

import (
	"context"
	"fmt"
	"gamatet/game/action"
	"gamatet/game/event"
	"gamatet/game/field"
	"gamatet/game/machine"
	"gamatet/game/op"
	"gamatet/game/piece"
	"gamatet/game/sweeper"
	"github.com/marko-gacesa/udpstar/joinchannel"
	"github.com/marko-gacesa/udpstar/udpstar/controller"
	"time"
)

var _ interface {
	Performer
	RenderRequester
	controller.Controller
} = (*GameHost)(nil)

type GameHost struct {
	// fixed setup
	fields    []hostFieldData
	inputs    []hostPlayerData
	suspendCh chan bool

	// state
	suspended bool
	paused    bool
}

type hostFieldData struct {
	Field       *field.Field
	Sweepers    []sweeper.Sweeper
	OutCh       chan<- []byte
	RenderReqCh chan field.RenderRequest

	events     event.List
	analyzer   op.Analyzer
	serializer serializer
}

type hostPlayerData struct {
	Name string
	field.PiecePlace
	InCh <-chan []byte // player actions, either direct local or from remote players
}

func MakeHost(setup Setup) *GameHost {
	var inputs []hostPlayerData
	fields := make([]hostFieldData, len(setup.Fields))

	for i := range setup.Fields {
		players := setup.Fields[i].Players

		width := setup.Config.WidthPerPlayer * len(players)
		height := setup.Config.Height

		f := field.Make(width, height, len(players))
		f.Idx = i
		f.Config = setup.Config.FieldConfig

		var sweepers []sweeper.Sweeper
		sweepers = append(sweepers, sweeper.NewRow(f))

		for j := range players {
			ctrl := f.Ctrl(byte(j))

			ctrl.Name = players[j].Name
			ctrl.Feed = setup.Config.PieceFeed
			ctrl.Config = players[j].Config
			ctrl.Level = setup.Config.Level
			ctrl.IsShadowShown = true

			ctrl.IsColumnLimited = setup.Config.PlayerZones
			ctrl.ColumnLimit = piece.ColumnLimit{
				Min: j * setup.Config.WidthPerPlayer,
				Max: (j+1)*setup.Config.WidthPerPlayer - 1,
			}

			pp := field.PiecePlace{
				FieldIdx: byte(i),
				CtrlIdx:  byte(j),
			}

			if players[j].OutCh != nil {
				panic(fmt.Sprintf("player %d in field %d should not have OutCh", j, i))
			}

			inputs = append(inputs, hostPlayerData{
				Name:       players[j].Name,
				PiecePlace: pp,
				InCh:       players[j].InCh,
			})
		}

		if setup.Fields[i].InCh != nil {
			panic(fmt.Sprintf("field %d should not have InCh", i))
		}

		fields[i] = hostFieldData{
			Field:       f,
			Sweepers:    sweepers,
			OutCh:       setup.Fields[i].OutCh,
			RenderReqCh: make(chan field.RenderRequest),
		}
	}

	return &GameHost{
		fields:    fields,
		inputs:    inputs,
		suspendCh: make(chan bool),
	}
}

func (g *GameHost) Perform(
	ctx context.Context,
) {
	ctx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()

	for _, f := range g.fields {
		f.Field.StartTimers()
	}
	defer func() {
		for _, f := range g.fields {
			f.Field.StopTimers()
			for _, s := range f.Sweepers {
				s.Pause()
			}
		}
	}()

	defer g.stop(ctx)

	stopCh := ctx.Done()

	ctrlTimer := joinchannel.Channel(ctx, func() <-chan joinchannel.Input[time.Time, field.PiecePlace] {
		ch := make(chan joinchannel.Input[time.Time, field.PiecePlace])
		go func() {
			defer close(ch)
			for fIdx, f := range g.fields {
				ctrlCount := byte(f.Field.Ctrls())
				for pIdx := byte(0); pIdx < ctrlCount; pIdx++ {
					ch <- joinchannel.Input[time.Time, field.PiecePlace]{
						ID: field.PiecePlace{
							FieldIdx: byte(fIdx),
							CtrlIdx:  pIdx,
						},
						Ch: f.Field.Ctrl(pIdx).Timer.C,
					}
				}
			}
		}()
		return ch
	}())

	inputCh := joinchannel.SlicePtr(ctx, g.inputs, func(p *hostPlayerData) <-chan []byte {
		return p.InCh
	})

	type sweeperPusher struct {
		sweeper sweeper.Sweeper
		pusher  event.Pusher
	}

	sweeperTimer := joinchannel.Channel(ctx, func() <-chan joinchannel.Input[time.Time, sweeperPusher] {
		ch := make(chan joinchannel.Input[time.Time, sweeperPusher])
		go func() {
			defer close(ch)
			for idx := range g.fields {
				for _, s := range g.fields[idx].Sweepers {
					ch <- joinchannel.Input[time.Time, sweeperPusher]{
						ID: sweeperPusher{
							sweeper: s,
							pusher:  &g.fields[idx].events,
						},
						Ch: s.Timer(),
					}
				}
			}
		}()
		return ch
	}())

	renderReqCh := joinchannel.SlicePtr(ctx, g.fields, func(fd *hostFieldData) <-chan field.RenderRequest {
		return fd.RenderReqCh
	})

	/////////////////////////////
	/*
		func(f *field.Field, events *event.List) {
			w := f.GetWidth()
			for i := 0; i <= 2; i++ {
				for d := 0; d <= i; d++ {
					putBlock(events, d, i-d, block.Wall)
					putBlock(events, w-1-d, i-d, block.Wall)
				}
				putBlock(events, i, 4-i, block.Block{
					Type:     block.TypeRock,
					Hardness: byte(1 + i),
					Color:    0x00FFFFFF,
				})
				putBlock(events, w-1-i, 4-i, block.Block{
					Type:     block.TypeRuby,
					Hardness: byte(1 + i),
					Color:    0xFFFF00FF,
				})
			}
			//for i := 3; i < 7; i++ {
			//	for j := 0; j < 18; j++ {
			//		putBlock(events, i, j, block.Block{Type: block.TypeRock, Hardness: byte(i - 3), Color: 0x90FF80FF})
			//		//putBlock(events, i, j, block.Iron)
			//	}
			//}
			conjureBlock(&g.fields[0].events, 0, 6, block.Goal)
			conjureBlock(&g.fields[0].events, 1, 5, block.Block{Type: block.TypeGoal, Hardness: 0, Color: 0x0000FFFF})
			conjureBlock(&g.fields[0].events, 7, 4, block.Iron)
			g.applyEvents(ctx)
		}(g.fields[0].Field, &g.fields[0].events)
	*/
	////////////////////////////

	for {
		for i := range g.fields {
			g.fields[i].events.Clear()
		}

		select {
		case <-stopCh:
			return

		case suspend := <-g.suspendCh:
			g.suspended = suspend
			if suspend {
				g.pause(ctx)
			}

		case inputData := <-inputCh:
			data := inputData.Data
			if len(data) != 1 {
				continue
			}

			pp := &g.inputs[inputData.ID]
			fIdx := pp.FieldIdx
			pIdx := pp.CtrlIdx

			f := g.fields[fIdx].Field
			ctrl := f.Ctrl(pIdx)
			events := &g.fields[fIdx].events

			a := action.Action(data[0])

			if a == action.Abort {
				if ctrl.State.IsAbortable() {
					return
				} else {
					a = action.Pause
				}
			}

			if a == action.Pause {
				if ctrl.State.IsPausable() {
					if g.paused {
						g.unpause(ctx)
					} else {
						g.pause(ctx)
					}
				}
				continue
			}

			if g.paused {
				continue
			}

			machine.HandleActionInput(f, ctrl, events, a)
			g.applyEvents(ctx)

		case fc := <-ctrlTimer:
			f := g.fields[fc.ID.FieldIdx].Field
			ctrl := f.Ctrl(fc.ID.CtrlIdx)
			events := &g.fields[fc.ID.FieldIdx].events

			machine.HandleTimeout(f, ctrl, events)
			g.applyEvents(ctx)

		case sw := <-sweeperTimer:
			sw.ID.sweeper.Sweep(sw.ID.pusher)
			g.applyEvents(ctx)

		case rr := <-renderReqCh:
			renderInfo := field.ObtainRenderInfo()
			f := g.fields[rr.ID].Field
			f.FillRenderInfo(renderInfo, field.GameInfo{
				Paused: g.paused,
			}, rr.Data.Time)
			go func(ctx context.Context, ch chan<- *field.RenderInfo) {
				select {
				case <-ctx.Done():
				case ch <- renderInfo:
				}
			}(ctx, rr.Data.RenderInfo)
		}
	}
}

func (g *GameHost) Suspend(ctx context.Context) {
	select {
	case <-ctx.Done():
	case g.suspendCh <- true:
	}
}

func (g *GameHost) Resume(ctx context.Context) {
	select {
	case <-ctx.Done():
	case g.suspendCh <- false:
	}
}

func (g *GameHost) RenderRequest(ctx context.Context, fieldIdx int, t time.Time, ch chan<- *field.RenderInfo) {
	select {
	case <-ctx.Done():
	case g.fields[fieldIdx].RenderReqCh <- field.RenderRequest{
		FieldIdx:   fieldIdx,
		Time:       t,
		RenderInfo: ch,
	}:
	}
}

func (g *GameHost) GetSize(idx int) (int, int) {
	f := g.fields[idx].Field
	return f.GetWidth(), f.GetHeight()
}

func (g *GameHost) stop(ctx context.Context) {
	for fIdx := range g.fields {
		select {
		case <-ctx.Done():
			return
		case g.fields[fIdx].OutCh <- op.FieldStopBytes:
		}
	}
}

func (g *GameHost) pause(ctx context.Context) {
	if g.paused {
		return
	}

	g.paused = true

	for fIdx := 0; fIdx < len(g.fields); fIdx++ {
		g.fields[fIdx].Field.Pause()
		for _, s := range g.fields[fIdx].Sweepers {
			s.Pause()
		}
		select {
		case <-ctx.Done():
			return
		case g.fields[fIdx].OutCh <- op.FieldPauseBytes:
		}
	}
}

func (g *GameHost) unpause(ctx context.Context) {
	if g.suspended || !g.paused {
		return
	}

	g.paused = false

	for fIdx := 0; fIdx < len(g.fields); fIdx++ {
		g.fields[fIdx].Field.Unpause()
		for _, s := range g.fields[fIdx].Sweepers {
			s.Unpause()
		}
		select {
		case <-ctx.Done():
			return
		case g.fields[fIdx].OutCh <- op.FieldUnpauseBytes:
		}
	}
}

func (g *GameHost) applyEvents(ctx context.Context) {
	for fIdx := range g.fields {
		fd := &g.fields[fIdx]

		if fd.events.IsEmpty() {
			continue
		}

		f := g.fields[fIdx].Field
		fd.analyzer.Reset()
		fd.events.Range(func(e event.Event) {
			fd.analyzer.Analyze(e)
			e.Do(f)
		})

		select {
		case <-ctx.Done():
		case fd.OutCh <- fd.serializer.Serialize(&fd.events):
		}

		for _, s := range g.fields[fIdx].Sweepers {
			s.Start(g.fields[fIdx].analyzer)
		}
	}
}
