// Copyright (c) 2020-2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package core

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/marko-gacesa/channel"
	"github.com/marko-gacesa/gamatet/game/action"
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/piece"
	"github.com/marko-gacesa/gamatet/logic/anim"
	"github.com/marko-gacesa/gamatet/logic/latency"
)

var _ interface {
	Performer
	RenderRequester
} = (*GameInterpreter)(nil)

type GameInterpreter struct {
	// fixed setup
	fields   []interpreterFieldData
	inputs   []interpreterPlayerData
	actionCh <-chan action.Action

	renderReqCh chan field.RenderRequest

	doneCh chan struct{}

	options InterpreterOptions
}

type InterpreterOptions struct {
	field.RenderOptions
	LocalPlayerActionCh chan<- []byte
	SinceLastContactFn  func() time.Duration
	Latencies           *latency.List
	StartPaused         bool
	StartUpDuration     time.Duration
}

type interpreterFieldData struct {
	Field *field.Field
	InCh  <-chan []byte
}

type interpreterPlayerData struct {
	Name string
	PiecePlace
}

func MakeInterpreter(setup Setup, options InterpreterOptions) *GameInterpreter {
	var inputs []interpreterPlayerData
	fields := make([]interpreterFieldData, len(setup.Fields))

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
		f.RenderOptions = options.RenderOptions
		f.SetStartUpDuration(options.StartUpDuration)

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

			if players[j].InCh != nil {
				panic(fmt.Sprintf("player=%d@field=%d shouldn't have InCh", j, i))
			}

			inputs = append(inputs, interpreterPlayerData{
				Name:       players[j].Name,
				PiecePlace: pp,
			})
		}

		if setup.Fields[i].OutCh != nil {
			panic(fmt.Sprintf("field=%d should not have OutCh", i))
		}

		f.UpdateBlocksRemoved(0)

		fields[i] = interpreterFieldData{
			Field: f,
			InCh:  setup.Fields[i].InCh,
		}
	}

	return &GameInterpreter{
		fields:      fields,
		inputs:      inputs,
		actionCh:    setup.ActionCh,
		renderReqCh: make(chan field.RenderRequest),
		doneCh:      make(chan struct{}),

		options: options,
	}
}

func (g *GameInterpreter) Perform(ctx context.Context) (result []PlayerResult) {
	fieldEventCh := channel.Context(ctx, channel.JoinSlicePtr(g.doneCh, g.fields, func(fd *interpreterFieldData) <-chan []byte {
		return fd.InCh
	}))

	fieldsDoneCh := channel.JoinSlicePtr(g.doneCh, g.fields, func(fd *interpreterFieldData) <-chan struct{} {
		return fd.Field.GetDone()
	})

	if g.options.StartPaused {
		for _, f := range g.fields {
			f.Field.SetState(field.StateGetReady)
		}
	}

	defer close(g.doneCh)

	defer func() {
		for fieldIdx := range g.fields {
			f := g.fields[fieldIdx].Field
			for ctrlIdx := range byte(f.Ctrls()) {
				ctrl := f.Ctrl(ctrlIdx)
				result = append(result, PlayerResult{
					FieldIdx:      byte(fieldIdx),
					CtrlIdx:       ctrlIdx,
					PlayerIndex:   ctrl.PlayerIndex,
					Outcome:       f.GetOutcome(),
					BlocksRemoved: f.GetBlocksRemoved(),
					Score:         ctrl.Score,
					PieceCount:    ctrl.PieceCount,
					Level:         ctrl.Level,
				})
			}
		}
	}()

	var s Serializer

	const serverLostDuration = time.Second * 5

	for {
		select {
		case <-ctx.Done():
			return

		case <-fieldsDoneCh:
			return

		case a := <-g.actionCh:
			switch a {
			case action.Abort:
				if g.options.SinceLastContactFn != nil && g.options.SinceLastContactFn() > serverLostDuration {
					return
				}
				if m := g.fields[0].Field.GetState(); m != field.StateNormal {
					return
				}
				if g.options.LocalPlayerActionCh != nil {
					g.options.LocalPlayerActionCh <- []byte{byte(a)}
				}
			case action.Pause:
				if g.options.LocalPlayerActionCh != nil {
					g.options.LocalPlayerActionCh <- []byte{byte(a)}
				}
			}

		case fieldEventData := <-fieldEventCh:
			var events event.List
			if err := s.Deserialize(fieldEventData.Data, &events); err != nil {
				log.Printf("failed to deserialize events: %s\n", err.Error())
				continue
			}

			f := g.fields[fieldEventData.ID].Field
			events.Range(func(e event.Event) {
				e.Do(f)
			})
			events.Clear()

		case rr := <-g.renderReqCh:
			renderInfo := rr.RenderInfo
			f := g.fields[rr.FieldIdx].Field
			f.FillRenderInfo(renderInfo, rr.Time)
			if g.options.SinceLastContactFn != nil && g.options.SinceLastContactFn() > serverLostDuration {
				renderInfo.State = field.StateServerLost
			}
			if g.options.Latencies != nil {
				renderInfo.TextData.Latencies = g.options.Latencies.String()
			}
			rr.Done <- true
		}
	}
}

func (g *GameInterpreter) RenderRequest(fieldIdx int, t time.Time, info *field.RenderInfo, chDone chan<- bool) {
	select {
	case <-g.doneCh:
		chDone <- false
	case g.renderReqCh <- field.RenderRequest{FieldIdx: fieldIdx, Time: t, RenderInfo: info, Done: chDone}:
	}
}

func (g *GameInterpreter) GetSize(idx int) (int, int, int) {
	f := g.fields[idx].Field
	return f.GetWidth(), f.GetHeight(), f.Ctrls()
}

func (g *GameInterpreter) AddAnim(anim anim.Anim) {
	for i := range g.fields {
		g.fields[i].Field.Anim(anim)
	}
}
