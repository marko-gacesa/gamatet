// Copyright (c) 2020-2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package core

import (
	"context"
	"fmt"
	"time"

	"github.com/marko-gacesa/channel"
	"github.com/marko-gacesa/gamatet/game/action"
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/machine"
	"github.com/marko-gacesa/gamatet/game/op"
	"github.com/marko-gacesa/gamatet/game/piece"
	"github.com/marko-gacesa/gamatet/game/sweeper"
	"github.com/marko-gacesa/gamatet/logic/anim"
	"github.com/marko-gacesa/gamatet/logic/latency"
	"github.com/marko-gacesa/udpstar/udpstar/controller"
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
	actionCh  <-chan action.Action
	suspendCh chan bool

	renderReqCh chan field.RenderRequest

	doneCh chan struct{}

	options HostOptions
	state   hostState
}

type HostOptions struct {
	field.RenderOptions
	Latencies   *latency.List
	Init        func(f *field.Field, p event.Pusher)
	StartPaused bool
}

type hostFieldData struct {
	Field      *field.Field
	Sweepers   []sweeper.Sweeper
	GnawKeeper *GnawKeeper
	OutCh      chan<- []byte

	events     event.List
	serializer Serializer
}

type hostPlayerData struct {
	Name string
	PiecePlace
	InCh <-chan []byte // player actions, either direct local or from remote players
}

func MakeHost(setup Setup, options HostOptions) *GameHost {
	if setup.ActionCh == nil {
		panic("ActionCh must not be nil")
	}

	var inputs []hostPlayerData
	fields := make([]hostFieldData, len(setup.Fields))

	for i := range setup.Fields {
		players := setup.Fields[i].Players

		width := setup.Config.WidthPerPlayer
		height := setup.Config.Height
		if len(players) > 0 {
			width *= len(players)
		}

		f := field.Make(width, height, len(players))
		f.Idx = i
		f.Config = setup.Config.FieldConfig
		f.Seed(setup.Config.RandomSeed)
		f.RenderOptions = options.RenderOptions

		for j := range players {
			ctrl := f.Ctrl(byte(j))

			ctrl.PlayerIndex = byte(players[j].Index)
			ctrl.Name = players[j].Name
			ctrl.ControlsStr = players[j].ControlsStr
			ctrl.Feed = piece.NewCtrlFeed(setup.Config.PieceFeed, i, j, setup.Config.SamePieces)
			ctrl.Config = players[j].Config
			ctrl.SetLevel(uint(setup.Config.Level))
			ctrl.IsShadowShown = true

			ctrl.IsColumnLimited = setup.Config.PlayerZones
			ctrl.ColumnLimit = piece.ColumnLimit{
				Min: j * setup.Config.WidthPerPlayer,
				Max: (j+1)*setup.Config.WidthPerPlayer - 1,
			}

			pp := PiecePlace{
				FieldIdx: byte(i),
				CtrlIdx:  byte(j),
			}

			if players[j].InCh == nil {
				panic(fmt.Sprintf("player=%d@field=%d should have OutCh", j, i))
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

		f.UpdateBlocksRemoved(0)

		fields[i] = hostFieldData{
			Field:      f,
			Sweepers:   nil, // set below
			GnawKeeper: NewGnawKeeper(f, uint(setup.Config.RandomSeed)),
			OutCh:      setup.Fields[i].OutCh,
		}
	}

	for i := range fields {
		f := fields[i].Field
		fields[i].Sweepers = append(fields[i].Sweepers,
			sweeper.NewRow(f),
			sweeper.NewScore(f),
			sweeper.NewShaker(f),
			sweeper.NewGameOver(f),
			sweeper.NewSpeedUp(f),
			sweeper.NewLingering(f),
		)
	}

	if len(fields) > 1 {
		for i := range fields {
			f := fields[i].Field
			others := getFieldPunishers(fields, i)

			if setup.Config.Shooters {
				fields[i].Sweepers = append(fields[i].Sweepers,
					sweeper.NewShotTransfer(f, others),
				)
			}

			fields[i].Sweepers = append(fields[i].Sweepers,
				sweeper.NewSpeedUpOnDefeat(f, others),
				sweeper.NewBlizzard(f, others),
				sweeper.NewMagic(f, others, setup.Config.RandomSeed, sweeper.MagicTypeAll),
			)
		}
	}

	if options.Init != nil {
		for i := range fields {
			f := fields[i].Field
			p := &fields[i].events
			options.Init(f, p)
		}
	}

	return &GameHost{
		fields:      fields,
		inputs:      inputs,
		actionCh:    setup.ActionCh,
		suspendCh:   make(chan bool),
		renderReqCh: make(chan field.RenderRequest),
		doneCh:      make(chan struct{}),

		options: options,

		state: hostStateNormal,
	}
}

func (g *GameHost) Perform(ctx context.Context) {
	for _, f := range g.fields {
		f.Field.StartTimers()
	}

	defer g._pauseAllFields()
	defer g.sendStop()
	defer close(g.doneCh)

	ctrlTimer := channel.Join(g.doneCh, func() <-chan channel.Input[time.Time, PiecePlace] {
		ch := make(chan channel.Input[time.Time, PiecePlace])
		go func() {
			defer close(ch)
			for fIdx, f := range g.fields {
				ctrlCount := byte(f.Field.Ctrls())
				for pIdx := range ctrlCount {
					ch <- channel.Input[time.Time, PiecePlace]{
						ID: PiecePlace{
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

	gnawTimer := channel.JoinSlicePtr(g.doneCh, g.fields, func(p *hostFieldData) <-chan time.Time {
		return p.GnawKeeper.Chan()
	})

	startDeley := time.Nanosecond
	if g.options.StartPaused {
		g.stateTransitionGetReady()
		startDeley = field.StartupDuration
	}

	g.applyEvents()

	getReadyTimer := time.NewTimer(startDeley)
	defer getReadyTimer.Stop()

	for {
		for i := range g.fields {
			g.fields[i].events.Clear()
		}

		select {
		case <-ctx.Done():
			return

		case <-getReadyTimer.C:
			g.stateTransitionPlay()
			g.applyEvents()

		case a := <-g.actionCh:
			switch a {
			case action.Abort:
				if g.state == hostStateFinish || g.state == hostStatePause {
					return
				}
				g.stateTransitionPauseToggle()
			case action.Pause:
				g.stateTransitionPauseToggle()
			default:
				continue
			}
			g.applyEvents()

		case suspend := <-g.suspendCh:
			if suspend {
				g.stateTransitionSuspend()
			} else {
				g.stateTransitionUnsuspend()
			}
			g.applyEvents()

		case inputData := <-inputCh:
			data := inputData.Data
			if len(data) != 1 {
				continue
			}

			a := action.Action(data[0])

			if a == action.NoOp {
				continue
			}

			pp := &g.inputs[inputData.ID]
			fIdx := pp.FieldIdx
			pIdx := pp.CtrlIdx

			f := g.fields[fIdx].Field
			ctrl := f.Ctrl(pIdx)
			events := &g.fields[fIdx].events

			if g.state == hostStatePause {
				a = action.Pause
			} else if a == action.Abort {
				if ctrl.State.IsAbortable() {
					return
				}
				a = action.Pause
			}

			if a == action.Pause {
				g.stateTransitionPauseToggle()
			}

			if g.state == hostStatePause {
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

		case gn := <-gnawTimer:
			f := &g.fields[gn.ID]
			f.GnawKeeper.ProcessAll(&f.events)
			g.applyEvents()

		case rr := <-g.renderReqCh:
			renderInfo := field.ObtainRenderInfo()
			f := g.fields[rr.FieldIdx].Field
			f.FillRenderInfo(renderInfo, rr.Time)
			if g.options.Latencies != nil {
				renderInfo.TextData.Latencies = g.options.Latencies.String()
			}
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

func (g *GameHost) GetSize(idx int) (int, int, int) {
	f := g.fields[idx].Field
	return f.GetWidth(), f.GetHeight(), f.Ctrls()
}

func (g *GameHost) AddAnim(anim anim.Anim) {
	for i := range g.fields {
		g.fields[i].Field.Anim(anim)
	}
}

func (g *GameHost) sendStop() {
	for fIdx := range g.fields {
		g.fields[fIdx].OutCh <- op.FieldStopBytes
		close(g.fields[fIdx].OutCh)
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
		g.stateTransitionFinish()
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
		g.stateTransitionFinish()
		f := g.fields[playingLastIdx].Field
		g.fields[playingLastIdx].events.Push(op.NewFieldMode(f, field.ModeVictory, true))
	}
}

func getFieldPunishers(fields []hostFieldData, exceptIdx int) []sweeper.FieldPunisher {
	list := make([]sweeper.FieldPunisher, 0, len(fields)-1)
	for i := range fields {
		if i == exceptIdx {
			continue
		}

		list = append(list, sweeper.FieldPunisher{
			Field:   fields[i].Field,
			Pusher:  &fields[i].events,
			GnawAdd: fields[i].GnawKeeper.AddSmall,
		})
	}

	return list
}
