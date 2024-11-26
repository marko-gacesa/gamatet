// Copyright (c) 2020-2024 by Marko Gaćeša

package core

import (
	"bytes"
	"context"
	"fmt"
	"gamatet/game/event"
	"gamatet/game/field"
	"gamatet/game/piece"
	"github.com/marko-gacesa/udpstar/joinchannel"
	"log"
	"time"
)

var _ interface {
	Performer
	RenderRequester
} = (*GameInterpreter)(nil)

type GameInterpreter struct {
	// fixed setup
	fields []interpreterFieldData
	inputs []interpreterPlayerData

	// state
	paused bool
}

type interpreterFieldData struct {
	Field       *field.Field
	InCh        <-chan []byte
	RenderReqCh chan field.RenderRequest

	// internal caches
	buffer bytes.Buffer
}

type interpreterPlayerData struct {
	Name string
	field.PiecePlace
}

func MakeInterpreter(setup Setup) *GameInterpreter {
	var inputs []interpreterPlayerData
	fields := make([]interpreterFieldData, len(setup.Fields))

	for i := range setup.Fields {
		players := setup.Fields[i].Players

		width := setup.Config.WidthPerPlayer * len(players)
		height := setup.Config.Height

		f := field.Make(width, height, len(players))
		f.Idx = i
		f.Config = setup.Config.FieldConfig

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

			if players[j].InCh != nil {
				panic(fmt.Sprintf("player %d i field %d should not have InCh", j, i))
			}

			if players[j].OutCh != nil {
				panic(fmt.Sprintf("player %d i field %d should not have OutCh", j, i))
			}

			inputs = append(inputs, interpreterPlayerData{
				Name:       players[j].Name,
				PiecePlace: pp,
			})
		}

		if setup.Fields[i].OutCh != nil {
			panic(fmt.Sprintf("field %d should not have OutCh", i))
		}

		fields[i] = interpreterFieldData{
			Field:       f,
			InCh:        setup.Fields[i].InCh,
			RenderReqCh: make(chan field.RenderRequest),
		}
	}

	return &GameInterpreter{
		fields: fields,
		inputs: inputs,
	}
}

func (g *GameInterpreter) Perform(ctx context.Context) {
	ctx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()

	fieldEventCh := joinchannel.SlicePtr(ctx, g.fields, func(fd *interpreterFieldData) <-chan []byte {
		return fd.InCh
	})

	renderReqCh := joinchannel.SlicePtr(ctx, g.fields, func(fd *interpreterFieldData) <-chan field.RenderRequest {
		return fd.RenderReqCh
	})

	var s serializer

	for {
		select {
		case <-ctx.Done():
			return

		case fieldEventData := <-fieldEventCh:
			var events event.List
			if err := s.Deserialize(fieldEventData.Data, &events); err != nil {
				log.Println(err.Error())
				continue
			}

			f := g.fields[fieldEventData.ID].Field
			events.Range(func(e event.Event) {
				e.Do(f)
			})
			events.Clear()

		case rr := <-renderReqCh:
			f := g.fields[rr.ID].Field
			renderInfo := f.GetRenderInfo(rr.Data.Time)
			go func(ctx context.Context, ch chan<- *field.RenderInfo) {
				select {
				case <-ctx.Done():
				case ch <- renderInfo:
				}
			}(ctx, rr.Data.RenderInfo)
		}
	}
}

func (g *GameInterpreter) RenderRequest(ctx context.Context, fieldIdx int, t time.Time, ch chan<- *field.RenderInfo) {
	select {
	case <-ctx.Done():
	case g.fields[fieldIdx].RenderReqCh <- field.RenderRequest{
		FieldIdx:   fieldIdx,
		Time:       t,
		RenderInfo: ch,
	}:
	}
}

func (g *GameInterpreter) GetSize(idx int) (int, int) {
	f := g.fields[idx].Field
	return f.GetWidth(), f.GetHeight()
}
