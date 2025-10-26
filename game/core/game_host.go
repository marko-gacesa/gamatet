// Copyright (c) 2020-2025 by Marko Gaćeša

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
	"github.com/marko-gacesa/udpstar/channel"
	"github.com/marko-gacesa/udpstar/udpstar/controller"
	"math/rand/v2"
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

	renderReqCh chan field.RenderRequest

	doneCh chan struct{}
}

type hostFieldData struct {
	Field    *field.Field
	Sweepers []sweeper.Sweeper
	OutCh    chan<- []byte

	events     event.List
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

	r := rand.New(rand.NewPCG(uint64(setup.Config.RandomSeed), 0))

	for i := range setup.Fields {
		players := setup.Fields[i].Players

		width := setup.Config.WidthPerPlayer * len(players)
		height := setup.Config.Height

		f := field.Make(width, height, len(players), field.WithRand(r))
		f.Idx = i
		f.Config = setup.Config.FieldConfig

		var sweepers []sweeper.Sweeper
		sweepers = append(sweepers, sweeper.NewRow(f))
		sweepers = append(sweepers, sweeper.NewShaker(f))

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

			inputs = append(inputs, hostPlayerData{
				Name:       players[j].Name,
				PiecePlace: pp,
				InCh:       players[j].InCh,
			})
		}

		if setup.Fields[i].InCh != nil {
			panic(fmt.Sprintf("field=%d should not have InCh", i))
		}

		fields[i] = hostFieldData{
			Field:    f,
			Sweepers: sweepers,
			OutCh:    setup.Fields[i].OutCh,
		}
	}

	for i := range fields {
		s := sweeper.NewGameOver(fields[i].Field)
		fields[i].Sweepers = append(fields[i].Sweepers, s)
	}

	if len(fields) > 1 {
		for i := range fields {
			s := sweeper.NewPunisher(fields[i].Field, getFieldPushers(fields, i))
			fields[i].Sweepers = append(fields[i].Sweepers, s)
		}
	}

	return &GameHost{
		fields:      fields,
		inputs:      inputs,
		suspendCh:   make(chan bool),
		renderReqCh: make(chan field.RenderRequest),
		doneCh:      make(chan struct{}),
	}
}

func (g *GameHost) Perform(ctx context.Context) {
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

	ctrlTimer := channel.Join(g.doneCh, func() <-chan channel.Input[time.Time, field.PiecePlace] {
		ch := make(chan channel.Input[time.Time, field.PiecePlace])
		go func() {
			defer close(ch)
			for fIdx, f := range g.fields {
				ctrlCount := byte(f.Field.Ctrls())
				for pIdx := byte(0); pIdx < ctrlCount; pIdx++ {
					ch <- channel.Input[time.Time, field.PiecePlace]{
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

	inputCh := channel.JoinSlicePtr(g.doneCh, g.inputs, func(p *hostPlayerData) <-chan []byte {
		return p.InCh
	})

	type sweeperPusher struct {
		sweeper sweeper.Sweeper
		pusher  event.Pusher
	}

	sweeperTimer := channel.Join(g.doneCh, func() <-chan channel.Input[time.Time, sweeperPusher] {
		ch := make(chan channel.Input[time.Time, sweeperPusher])
		go func() {
			defer close(ch)
			for idx := range g.fields {
				for _, s := range g.fields[idx].Sweepers {
					ch <- channel.Input[time.Time, sweeperPusher]{
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

	defer g.sendStop()

	defer close(g.doneCh)

	for {
		for i := range g.fields {
			g.fields[i].events.Clear()
		}

		select {
		case <-ctx.Done():
			return

		case suspend := <-g.suspendCh:
			if suspend {
				g.suspend()
			} else {
				g.unsuspend()
			}
			g.applyEvents()

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

			if g.paused && a == action.Drop {
				a = action.Pause
			} else if a == action.Abort {
				if ctrl.State.IsAbortable() {
					return
				}

				a = action.Pause
			}

			if a == action.Pause {
				g.pauseToggle(ctrl)
			}

			if g.paused {
				a = action.NoOp
			}

			machine.HandleActionInput(f, ctrl, events, a)
			g.applyEvents()

		case fc := <-ctrlTimer:
			f := g.fields[fc.ID.FieldIdx].Field
			ctrl := f.Ctrl(fc.ID.CtrlIdx)
			events := &g.fields[fc.ID.FieldIdx].events

			if isDone := machine.HandleTimeout(f, ctrl, events); isDone {
				g.checkWinner(int(fc.ID.FieldIdx))
			}

			g.applyEvents()

		case sw := <-sweeperTimer:
			sw.ID.sweeper.Sweep(sw.ID.pusher)
			g.applyEvents()

		case rr := <-g.renderReqCh:
			renderInfo := field.ObtainRenderInfo()
			f := g.fields[rr.FieldIdx].Field
			f.FillRenderInfo(renderInfo, rr.Time)
			rr.RenderInfo <- renderInfo
		}
	}
}

func (g *GameHost) Suspend() {
	select {
	case <-g.doneCh:
	case g.suspendCh <- true:
	}
}

func (g *GameHost) Resume() {
	select {
	case <-g.doneCh:
	case g.suspendCh <- false:
	}
}

func (g *GameHost) RenderRequest(fieldIdx int, t time.Time, ch chan<- *field.RenderInfo) {
	select {
	case <-g.doneCh:
		close(ch)
	case g.renderReqCh <- field.RenderRequest{FieldIdx: fieldIdx, Time: t, RenderInfo: ch}:
	}
}

func (g *GameHost) GetSize(idx int) (int, int, []piece.DisplayPosition) {
	f := g.fields[idx].Field
	return f.GetWidth(), f.GetHeight(), f.CtrlInfoPositions()
}

func (g *GameHost) sendStop() {
	for fIdx := range g.fields {
		g.fields[fIdx].OutCh <- op.FieldStopBytes
	}
}

func (g *GameHost) pauseToggle(ctrl *piece.Ctrl) {
	if g.suspended {
		return
	}

	if g.paused {
		g.unpauseAllFields()
		g.paused = false
	} else if ctrl.State.IsPausable() {
		g.pauseAllFields(field.ModePause)
		g.paused = true
	}
}

func (g *GameHost) suspend() {
	if g.suspended {
		return
	}

	g.suspended = true
	g.paused = true
	g.pauseAllFields(field.ModeSuspended)
}

func (g *GameHost) unsuspend() {
	if !g.suspended {
		return
	}

	g.suspended = false
	g.unsuspendAllFields()
}

func (g *GameHost) unsuspendAllFields() {
	for fIdx := 0; fIdx < len(g.fields); fIdx++ {
		f := g.fields[fIdx].Field
		oldMode := f.GetMode()
		if oldMode == field.ModeSuspended {
			g.fields[fIdx].events.Push(op.NewFieldMode(f, field.ModePause, false))
		}
	}
}

func (g *GameHost) pauseAllFields(newMode field.Mode) {
	for fIdx := 0; fIdx < len(g.fields); fIdx++ {
		f := g.fields[fIdx].Field
		oldMode := f.GetMode()
		if oldMode == field.ModeNormal {
			g.fields[fIdx].events.Push(op.NewFieldMode(f, newMode, false))
			g.pauseField(fIdx)
		}
	}
}

func (g *GameHost) unpauseAllFields() {
	for fIdx := 0; fIdx < len(g.fields); fIdx++ {
		f := g.fields[fIdx].Field
		oldMode := f.GetMode()
		if oldMode == field.ModePause {
			g.fields[fIdx].events.Push(op.NewFieldMode(f, field.ModeNormal, false))
			g.unpauseField(fIdx)
		}
	}
}

func (g *GameHost) pauseField(fIdx int) {
	g.fields[fIdx].Field.Pause()
	for _, s := range g.fields[fIdx].Sweepers {
		s.Pause()
	}
}

func (g *GameHost) unpauseField(fIdx int) {
	g.fields[fIdx].Field.Unpause()
	for _, s := range g.fields[fIdx].Sweepers {
		s.Unpause()
	}
}

func (g *GameHost) applyEvents() {
	for fIdx := range g.fields {
		fd := &g.fields[fIdx]

		if fd.events.IsEmpty() {
			continue
		}

		f := g.fields[fIdx].Field
		analyzer := &sweeper.Analyzer{Field: f}

		fd.events.Range(func(e event.Event) {
			analyzer.Analyze(e)
			e.Do(f)
		})

		fd.OutCh <- fd.serializer.Serialize(&fd.events)

		for _, s := range g.fields[fIdx].Sweepers {
			s.Start(analyzer)
		}
	}
}

func (g *GameHost) checkWinner(loserIdx int) {
	var (
		playingLastIdx int
		playingCount   int
	)

	if len(g.fields) == 1 {
		g.fields[0].events.Push(op.NewFieldMode(g.fields[0].Field, field.ModeGameOver, true))
		return
	}

	playingLastIdx = -1
	for fIdx := range g.fields {
		f := g.fields[fIdx].Field

		if fIdx == loserIdx {
			g.fields[loserIdx].events.Push(op.NewFieldMode(f, field.ModeDefeat, true))
			continue
		}

		if !f.IsFinished() {
			playingCount++
			playingLastIdx = fIdx
			continue
		}
	}

	if playingCount == 1 {
		f := g.fields[playingLastIdx].Field
		g.fields[playingLastIdx].events.Push(op.NewFieldMode(f, field.ModeVictory, true))
	}
}

func getFieldPushers(fields []hostFieldData, exceptIdx int) []sweeper.FieldPusher {
	list := make([]sweeper.FieldPusher, 0, len(fields)-1)
	for i := range fields {
		if i == exceptIdx {
			continue
		}

		list = append(list, sweeper.FieldPusher{
			Field:  fields[i].Field,
			Pusher: &fields[i].events,
		})
	}

	return list
}
