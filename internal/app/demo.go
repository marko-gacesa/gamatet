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

func (app *App) demo(ctx screen.Context, getProgram func() ([]programOp, piece.Feed)) types.DemoParams {
	const (
		fullW  = 30
		fullH  = 20
		fieldW = 5
		fieldH = 12
		offsX  = 12
		offsY  = 0
	)

	program, feed := getProgram()

	programCh := make(chan []byte)
	go func() {
		t := time.NewTimer(50 * time.Millisecond)
		defer t.Stop()

		defer close(programCh)

		var serializer core.Serializer

		var idx int
		var list event.List

		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				for _, o := range program[idx].Event {
					list.Push(o)
				}
				programCh <- serializer.Serialize(&list)
				list.Clear()

				t.Reset(program[idx].Delay)

				idx++
				if idx == len(program) {
					idx = 0
				}
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
			PieceFeed:  feed,
			SamePieces: true,
			Shooters:   false,
		},
		Fields: []core.FieldSetup{
			{
				InCh: programCh,
				Players: []core.PlayerSetup{
					{
						Config:  piece.Config{},
						IsLocal: false,
						Index:   0,
						InCh:    nil,
					},
				},
			},
		},
		ActionCh: nil,
	}, core.InterpreterOptions{
		RenderOptions: field.RenderOptions{
			HideFrame:   true,
			HideBack:    false,
			HideShadows: false,
		},
	})

	go demo.Perform(ctx)

	return types.DemoParams{
		FullW: fullW,
		FullH: fullH,
		OffsX: offsX,
		OffsY: offsY,
		RotY:  -math.Pi / 16,
		Demo:  demo,
		Done:  ctx.Done(),
	}
}

type programOp struct {
	Event []event.Event
	Delay time.Duration
}

const (
	programDelay     = 500 * time.Millisecond
	programDelayLong = 2500 * time.Millisecond
)

func _programInit(t op.Type) []event.Event {
	return append([]event.Event{},
		op.NewFieldBlockSet(0, 0, t, field.AnimNo, 0, block.Rock),
		op.NewFieldBlockSet(2, 0, t, field.AnimNo, 0, block.Rock),
		op.NewFieldBlockSet(4, 0, t, field.AnimNo, 0, block.Rock),
		op.NewFieldBlockSet(0, 1, t, field.AnimNo, 0, block.Rock),
		op.NewFieldBlockSet(4, 1, t, field.AnimNo, 0, block.Rock),
		op.NewFieldBlockSet(0, 2, t, field.AnimNo, 0, block.Rock),
		op.NewFieldBlockSet(1, 2, t, field.AnimNo, 0, block.Rock),
		op.NewFieldBlockSet(0, 3, t, field.AnimNo, 0, block.Rock),
		op.NewFieldBlockSet(2, 3, t, field.AnimNo, 0, block.Rock),
		op.NewFieldBlockSet(1, 4, t, field.AnimNo, 0, block.Rock),
		op.NewFieldBlockSet(2, 4, t, field.AnimNo, 0, block.Rock),
		op.NewFieldBlockSet(1, 5, t, field.AnimNo, 0, block.Rock),
	)
}

func _programDescend() []event.Event {
	return append([]event.Event{},
		op.NewPieceFall(0, 1),
		op.NewPieceState(0, piece.StateDescend, piece.StateDescend, 0, 0),
	)
}

func programLava() ([]programOp, piece.Feed) {
	b := block.Lava
	program := []programOp{
		{Event: _programInit(op.TypeSet), Delay: 0},
		{
			Event: append([]event.Event{},
				op.NewPieceState(0, piece.StateInit, piece.StateNew, 0, 0),
				op.NewPieceSet(0, op.TypeSet, 1, 11, piece.NewTetromino(piece.TetrominoJ, b), 1),
				op.NewPieceState(0, piece.StateNew, piece.StateDescend, 0, 0),
			),
			Delay: programDelay,
		},
		{Event: _programDescend(), Delay: programDelay},
		{Event: _programDescend(), Delay: programDelay},
		{Event: _programDescend(), Delay: programDelay},
		{Event: _programDescend(), Delay: programDelay},
		{
			Event: append([]event.Event{},
				op.NewPieceState(0, piece.StateDescend, piece.StateSlide, 0, 0),
				op.NewFieldBlockSet(3, 0, op.TypeSet, field.AnimFall, 6, block.Block{Type: block.TypeRock, Hardness: 0, Color: b.Color}),
				op.NewFieldBlockSet(2, 5, op.TypeSet, field.AnimFall, 1, block.Block{Type: block.TypeRock, Hardness: 0, Color: b.Color}),
				op.NewFieldBlockSet(1, 6, op.TypeSet, field.AnimFall, 0, block.Block{Type: block.TypeRock, Hardness: 0, Color: b.Color}),
				op.NewFieldBlockSet(1, 7, op.TypeSet, field.AnimFall, 0, block.Block{Type: block.TypeRock, Hardness: 0, Color: b.Color}),
				op.NewPieceSet(0, op.TypeClear, 1, 7, piece.NewTetromino(piece.TetrominoJ, b), 1),
			),
			Delay: programDelayLong,
		},
		{
			Event: append([]event.Event{},
				op.NewFieldBlockSet(3, 0, op.TypeClear, field.AnimNo, 0, block.Block{Type: block.TypeRock, Hardness: 0, Color: b.Color}),
				op.NewFieldBlockSet(2, 5, op.TypeClear, field.AnimNo, 0, block.Block{Type: block.TypeRock, Hardness: 0, Color: b.Color}),
				op.NewFieldBlockSet(1, 6, op.TypeClear, field.AnimNo, 0, block.Block{Type: block.TypeRock, Hardness: 0, Color: b.Color}),
				op.NewFieldBlockSet(1, 7, op.TypeClear, field.AnimNo, 0, block.Block{Type: block.TypeRock, Hardness: 0, Color: b.Color}),
			),
			Delay: 0,
		},
		{Event: _programInit(op.TypeClear), Delay: 0},
	}

	return program, piece.SamePieceFeed{Piece: piece.NewTetromino(piece.TetrominoL, b)}
}

func programAcid() ([]programOp, piece.Feed) {
	b := block.Acid
	program := []programOp{
		{Event: _programInit(op.TypeSet), Delay: 0},
		{
			Event: append([]event.Event{},
				op.NewPieceState(0, piece.StateInit, piece.StateNew, 0, 0),
				op.NewPieceSet(0, op.TypeSet, 1, 11, piece.NewTetromino(piece.TetrominoJ, b), 1),
				op.NewPieceState(0, piece.StateNew, piece.StateDescend, 0, 0),
			),
			Delay: programDelay,
		},
		{Event: _programDescend(), Delay: programDelay},
		{Event: _programDescend(), Delay: programDelay},
		{Event: _programDescend(), Delay: programDelay},
		{Event: _programDescend(), Delay: programDelay},
		{
			Event: append([]event.Event{},
				op.NewPieceState(0, piece.StateDescend, piece.StateSlide, 0, 0),
				op.NewFieldExBlock(3, 0, field.AnimFall, 6, b),
				op.NewFieldExBlock(2, 4, field.AnimFall, 2, b),
				op.NewFieldBlockSet(2, 4, op.TypeClear, field.AnimPop, 0, block.Rock),
				op.NewFieldExBlock(1, 5, field.AnimFall, 1, b),
				op.NewFieldBlockSet(1, 5, op.TypeClear, field.AnimPop, 0, block.Rock),
				op.NewFieldExBlock(1, 4, field.AnimFall, 3, b),
				op.NewFieldBlockSet(1, 4, op.TypeClear, field.AnimPop, 0, block.Rock),
				op.NewPieceSet(0, op.TypeClear, 1, 7, piece.NewTetromino(piece.TetrominoJ, b), 1),
			),
			Delay: programDelayLong,
		},
		{
			Event: append([]event.Event{},
				op.NewFieldBlockSet(2, 4, op.TypeSet, field.AnimNo, 0, block.Rock),
				op.NewFieldBlockSet(1, 5, op.TypeSet, field.AnimNo, 0, block.Rock),
				op.NewFieldBlockSet(1, 4, op.TypeSet, field.AnimNo, 0, block.Rock),
			),
			Delay: 0,
		},
		{Event: _programInit(op.TypeClear), Delay: 0},
	}

	return program, piece.SamePieceFeed{Piece: piece.NewTetromino(piece.TetrominoL, b)}
}

func programCurl() ([]programOp, piece.Feed) {
	b := block.Curl
	program := []programOp{
		{Event: _programInit(op.TypeSet), Delay: 0},
		{
			Event: append([]event.Event{},
				op.NewPieceState(0, piece.StateInit, piece.StateNew, 0, 0),
				op.NewPieceSet(0, op.TypeSet, 1, 11, piece.NewTetromino(piece.TetrominoJ, b), 1),
				op.NewPieceState(0, piece.StateNew, piece.StateDescend, 0, 0),
			),
			Delay: programDelay,
		},
		{Event: _programDescend(), Delay: programDelay},
		{Event: _programDescend(), Delay: programDelay},
		{Event: _programDescend(), Delay: programDelay},
		{Event: _programDescend(), Delay: programDelay},
		{
			Event: append([]event.Event{},
				op.NewPieceState(0, piece.StateDescend, piece.StateSlide, 0, 0),
				op.NewFieldExBlock(3, 6, field.AnimDestroy, 0, block.Block{Type: block.TypeRock, Color: b.Color}),
				op.NewFieldBlockSet(2, 2, op.TypeSet, field.AnimFall, 4, block.Block{Type: block.TypeRock, Color: b.Color}),
				op.NewFieldBlockSet(1, 3, op.TypeSet, field.AnimFall, 3, block.Block{Type: block.TypeRock, Color: b.Color}),
				op.NewFieldBlockSet(1, 1, op.TypeSet, field.AnimFall, 6, block.Block{Type: block.TypeRock, Color: b.Color}),
				op.NewPieceSet(0, op.TypeClear, 1, 7, piece.NewTetromino(piece.TetrominoJ, b), 1),
			),
			Delay: programDelayLong,
		},
		{
			Event: append([]event.Event{},
				op.NewFieldBlockSet(2, 2, op.TypeClear, field.AnimNo, 0, block.Block{Type: block.TypeRock, Color: b.Color}),
				op.NewFieldBlockSet(1, 3, op.TypeClear, field.AnimNo, 0, block.Block{Type: block.TypeRock, Color: b.Color}),
				op.NewFieldBlockSet(1, 1, op.TypeClear, field.AnimNo, 0, block.Block{Type: block.TypeRock, Color: b.Color}),
			),
			Delay: 0,
		},
		{Event: _programInit(op.TypeClear), Delay: 0},
	}

	return program, piece.SamePieceFeed{Piece: piece.NewTetromino(piece.TetrominoL, b)}
}

func programWave() ([]programOp, piece.Feed) {
	b := block.Wave
	program := []programOp{
		{Event: _programInit(op.TypeSet), Delay: 0},
		{
			Event: append([]event.Event{},
				op.NewPieceState(0, piece.StateInit, piece.StateNew, 0, 0),
				op.NewPieceSet(0, op.TypeSet, 1, 11, piece.NewTetromino(piece.TetrominoJ, b), 1),
				op.NewPieceState(0, piece.StateNew, piece.StateDescend, 0, 0),
			),
			Delay: programDelay,
		},
		{Event: _programDescend(), Delay: programDelay},
		{Event: _programDescend(), Delay: programDelay},
		{Event: _programDescend(), Delay: programDelay},
		{Event: _programDescend(), Delay: programDelay},
		{
			Event: append([]event.Event{},
				op.NewPieceState(0, piece.StateDescend, piece.StateSlide, 0, 0),
				op.NewFieldExBlock(3, 6, field.AnimDestroy, 0, block.Block{Type: block.TypeRock, Color: b.Color}),
				op.NewFieldBlockSet(2, 1, op.TypeSet, field.AnimFall, 5, block.Block{Type: block.TypeRock, Color: b.Color}),
				op.NewFieldBlockSet(1, 0, op.TypeSet, field.AnimFall, 6, block.Block{Type: block.TypeRock, Color: b.Color}),
				op.NewFieldBlockSet(1, 1, op.TypeSet, field.AnimFall, 6, block.Block{Type: block.TypeRock, Color: b.Color}),
				op.NewPieceSet(0, op.TypeClear, 1, 7, piece.NewTetromino(piece.TetrominoJ, b), 1),
			),
			Delay: programDelayLong,
		},
		{
			Event: append([]event.Event{},
				op.NewFieldBlockSet(2, 1, op.TypeClear, field.AnimNo, 0, block.Block{Type: block.TypeRock, Color: b.Color}),
				op.NewFieldBlockSet(1, 0, op.TypeClear, field.AnimNo, 0, block.Block{Type: block.TypeRock, Color: b.Color}),
				op.NewFieldBlockSet(1, 1, op.TypeClear, field.AnimNo, 0, block.Block{Type: block.TypeRock, Color: b.Color}),
			),
			Delay: 0,
		},
		{Event: _programInit(op.TypeClear), Delay: 0},
	}

	return program, piece.SamePieceFeed{Piece: piece.NewTetromino(piece.TetrominoL, b)}
}

func programPieceTypeRot() ([]programOp, piece.Feed) {
	p := piece.NewTetromino(piece.TetrominoJ, block.Block{Type: block.TypeRock, Color: piece.DefaultColor{}.Color(0, 0)})
	program := []programOp{
		{
			Event: append([]event.Event{},
				op.NewPieceState(0, piece.StateInit, piece.StateNew, 0, 0),
				op.NewPieceSet(0, op.TypeSet, 1, 11, p, 1),
				op.NewPieceState(0, piece.StateNew, piece.StateDescend, 0, 0),
			),
			Delay: programDelay,
		},
		{Event: _programDescend(), Delay: programDelay / 2},
		{Event: []event.Event{op.NewPieceRotate(0, false)}, Delay: programDelay / 2},
		{Event: _programDescend(), Delay: programDelay / 2},
		{Event: []event.Event{op.NewPieceRotate(0, false)}, Delay: programDelay / 2},
		{Event: _programDescend(), Delay: programDelay / 2},
		{Event: []event.Event{op.NewPieceRotate(0, false)}, Delay: programDelay / 2},
		{Event: _programDescend(), Delay: programDelay / 2},
		{Event: []event.Event{op.NewPieceRotate(0, false)}, Delay: programDelay / 2},
		{Event: _programDescend(), Delay: programDelay / 2},
		{Event: []event.Event{op.NewPieceRotate(0, false)}, Delay: programDelay / 2},
		{Event: _programDescend(), Delay: programDelay / 2},
		{Event: []event.Event{op.NewPieceFlip(0)}, Delay: programDelay / 2},
		{Event: []event.Event{op.NewPieceSet(0, op.TypeClear, 1, 5, p, 1)}, Delay: programDelay},
	}

	return program, piece.SamePieceFeed{Piece: p}
}

func programPieceTypeMirrorV() ([]programOp, piece.Feed) {
	p := piece.NewFlipVTetromino(8, block.Block{Type: block.TypeRock, Color: piece.DefaultColor{}.Color(0, 0)})
	return programPieceTypeMirror(p)
}

func programPieceTypeMirrorH() ([]programOp, piece.Feed) {
	p := piece.NewFlipHTetromino(8, block.Block{Type: block.TypeRock, Color: piece.DefaultColor{}.Color(0, 0)})
	return programPieceTypeMirror(p)
}

func programPieceTypeMirror(p piece.Piece) ([]programOp, piece.Feed) {
	program := []programOp{
		{
			Event: append([]event.Event{},
				op.NewPieceState(0, piece.StateInit, piece.StateNew, 0, 0),
				op.NewPieceSet(0, op.TypeSet, 1, 11, p, 1),
				op.NewPieceState(0, piece.StateNew, piece.StateDescend, 0, 0),
			),
			Delay: programDelay,
		},
		{Event: _programDescend(), Delay: programDelay / 2},
		{Event: []event.Event{op.NewPieceFlip(0)}, Delay: programDelay / 2},
		{Event: _programDescend(), Delay: programDelay / 2},
		{Event: []event.Event{op.NewPieceFlip(0)}, Delay: programDelay / 2},
		{Event: _programDescend(), Delay: programDelay / 2},
		{Event: []event.Event{op.NewPieceFlip(0)}, Delay: programDelay / 2},
		{Event: _programDescend(), Delay: programDelay / 2},
		{Event: []event.Event{op.NewPieceFlip(0)}, Delay: programDelay / 2},
		{Event: _programDescend(), Delay: programDelay / 2},
		{Event: []event.Event{op.NewPieceFlip(0)}, Delay: programDelay / 2},
		{Event: _programDescend(), Delay: programDelay / 2},
		{Event: []event.Event{op.NewPieceFlip(0)}, Delay: programDelay / 2},
		{Event: []event.Event{op.NewPieceSet(0, op.TypeClear, 1, 5, p, 1)}, Delay: programDelay},
	}

	return program, piece.SamePieceFeed{Piece: p}
}

func programPieceTypeShooter() ([]programOp, piece.Feed) {
	p := piece.Shooter(5, block.TypeRock)
	program := []programOp{
		{
			Event: append([]event.Event{},
				op.NewPieceState(0, piece.StateInit, piece.StateNew, 0, 0),
				op.NewPieceSet(0, op.TypeSet, 2, 11, p, 1),
				op.NewPieceState(0, piece.StateNew, piece.StateDescend, 0, 0),
			),
			Delay: programDelay,
		},
		{Event: _programDescend(), Delay: programDelay / 2},
		{
			Event: []event.Event{
				op.NewFieldBlockSet(2, 0, op.TypeSet, field.AnimFall, 10, block.Rock),
				op.NewPieceShoot(0, true, block.TypeRock),
			}, Delay: programDelay / 2,
		},
		{Event: _programDescend(), Delay: programDelay / 2},
		{
			Event: []event.Event{
				op.NewFieldBlockSet(2, 1, op.TypeSet, field.AnimFall, 8, block.Rock),
				op.NewPieceShoot(0, true, block.TypeRock),
			}, Delay: programDelay / 2,
		},
		{Event: _programDescend(), Delay: programDelay / 2},
		{
			Event: []event.Event{
				op.NewFieldBlockSet(2, 2, op.TypeSet, field.AnimFall, 6, block.Rock),
				op.NewPieceShoot(0, true, block.TypeRock),
			}, Delay: programDelay / 2,
		},
		{Event: _programDescend(), Delay: programDelay},
		{
			Event: []event.Event{
				op.NewPieceSet(0, op.TypeClear, 2, 5, p, 1),
				op.NewFieldBlockSet(2, 0, op.TypeClear, field.AnimNo, 0, block.Rock),
				op.NewFieldBlockSet(2, 1, op.TypeClear, field.AnimNo, 0, block.Rock),
				op.NewFieldBlockSet(2, 2, op.TypeClear, field.AnimNo, 0, block.Rock),
			},
			Delay: 0,
		},
	}

	return program, piece.SamePieceFeed{Piece: p}
}
