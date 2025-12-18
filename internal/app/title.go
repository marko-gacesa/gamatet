// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	"math"
	"time"

	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/core"
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/op"
	"github.com/marko-gacesa/gamatet/game/piece"
	"github.com/marko-gacesa/gamatet/internal/types"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

var titleFirst bool

func (app *App) title(ctx screen.Context) types.DemoParams {
	const (
		fullW  = 39
		fullH  = 18
		fieldW = 37
		fieldH = 5
		offsX  = 0
		offsY  = 5
	)

	programCh := make(chan []byte)
	go func() {
		t := time.NewTimer(50 * time.Millisecond)

		defer close(programCh)

		var serializer core.Serializer
		program := _titleProgram(0, fieldH-1, op.TypeSet, field.AnimNo)

		var delay time.Duration
		if !titleFirst {
			delay = 40 * time.Millisecond
			titleFirst = true
		} else {
			delay = 8 * time.Millisecond
		}

		var idx int
		var list event.List

		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if idx == len(program) {
					return
				}

				list.Push(program[idx])
				programCh <- serializer.Serialize(&list)
				list.Clear()

				idx++

				t.Reset(delay)
			}
		}
	}()

	demo := core.MakeInterpreter(core.Setup{
		Name: "",
		Config: core.GameConfig{
			WidthPerPlayer: fieldW,
			Height:         fieldH,
			Level:          0,
			PlayerZones:    false,
			FieldConfig: field.Config{
				PieceCollision: false,
				Anim:           true,
			},
			RandomSeed: 0,
			PieceFeed:  nil,
			SamePieces: true,
			Shooters:   false,
		},
		Fields: []core.FieldSetup{
			{InCh: programCh},
		},
		ActionCh: nil,
	}, core.InterpreterOptions{
		RenderOptions: field.RenderOptions{
			HideFrame:   true,
			HideBack:    true,
			HideShadows: true,
		},
	})

	go demo.Perform(ctx)

	return types.DemoParams{
		FullW: fullW,
		FullH: fullH,
		OffsX: offsX,
		OffsY: offsY,
		RotX:  -math.Pi / 16,
		Demo:  demo,
		Done:  ctx.Done(),
	}
}

func _titleProgram(posX, posY int, opType op.Type, anim int) event.Slice {
	events := event.Slice(make([]event.Event, 0, 86))
	var curX, curY int

	setXY := func(x, y int) {
		curX, curY = x, y
	}

	putXY := func(x, y, k int) {
		events.Push(op.NewFieldBlockSet(
			x, y, opType, anim, 0, block.Block{
				Type:  block.TypeRock,
				Color: piece.DefaultColor{}.Color(uint(k), 0),
			},
		))
	}

	drawXY := func(n, k, deltaX, deltaY int) {
		for range n {
			putXY(curX, curY, k)
			curX += deltaX
			curY += deltaY
		}
	}

	drawL := func(n, k int) { drawXY(n, k, -1, 0) }
	drawR := func(n, k int) { drawXY(n, k, 1, 0) }
	drawU := func(n, k int) { drawXY(n, k, 0, 1) }
	drawD := func(n, k int) { drawXY(n, k, 0, -1) }

	const (
		colorG = 1
		colorA = 5
		colorM = 6
		colorT = 2
		colorE = 0
	)

	// G

	setXY(posX+3, posY)
	drawL(3, colorG)
	drawD(4, colorG)
	drawR(3, colorG)
	drawU(2, colorG)
	drawL(2, colorG)

	// A

	setXY(posX+5, posY-4)
	drawU(4, colorA)
	drawR(3, colorA)
	drawD(5, colorA)
	setXY(posX+6, posY-2)
	drawR(2, colorA)

	// M

	setXY(posX+10, posY-4)
	drawU(5, colorM)
	setXY(posX+11, posY-1)
	drawR(1, colorM)
	setXY(posX+12, posY-2)
	drawR(1, colorM)
	setXY(posX+13, posY-1)
	drawR(1, colorM)
	setXY(posX+14, posY)
	drawD(5, colorM)

	// A

	setXY(posX+16, posY-4)
	drawU(4, colorA)
	drawR(3, colorA)
	drawD(5, colorA)
	setXY(posX+17, posY-2)
	drawR(2, colorA)

	// T

	setXY(posX+21, posY)
	drawR(5, colorT)
	setXY(posX+23, posY-1)
	drawD(4, colorT)

	// E

	setXY(posX+30, posY)
	drawL(3, colorE)
	drawD(4, colorE)
	drawR(4, colorE)
	setXY(posX+28, posY-2)
	drawR(2, colorE)

	// T

	setXY(posX+32, posY)
	drawR(5, colorT)
	setXY(posX+34, posY-1)
	drawD(4, colorT)

	return events
}
