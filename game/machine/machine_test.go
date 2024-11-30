// Copyright (c) 2020-2024 by Marko Gaćeša

package machine

import (
	"gamatet/game/block"
	"gamatet/game/event"
	"gamatet/game/field"
	"gamatet/game/op"
	"gamatet/game/piece"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
)

func TestPieceMeld(t *testing.T) {
	tests := []struct {
		name   string
		x, y   int
		p      piece.Piece
		blocks []block.XYB
		events []event.Event
	}{
		{
			name: "rock: simple test",
			x:    0, y: 3, p: piece.NewTetromino(piece.TetrominoO, block.Block{Type: block.TypeRock}),
			blocks: []block.XYB{
				{XY: block.XY{X: 1, Y: 1}, Block: block.Wall},
			},
			events: []event.Event{
				&op.FieldBlockSet{Col: 1, Row: 2, Op: op.TypeSet, Block: block.Rock},
				&op.FieldBlockSet{Col: 0, Row: 2, Op: op.TypeSet, Block: block.Rock},
				&op.FieldBlockSet{Col: 1, Row: 3, Op: op.TypeSet, Block: block.Rock},
				&op.FieldBlockSet{Col: 0, Row: 3, Op: op.TypeSet, Block: block.Rock},
			},
		},
		// liquid piece tests
		{
			name: "lava: empty and wall",
			x:    0, y: 3, p: piece.NewTetromino(piece.TetrominoO, block.Block{Type: block.TypeLava}),
			blocks: []block.XYB{
				{XY: block.XY{X: 1, Y: 1}, Block: block.Wall},
			},
			events: []event.Event{
				&op.FieldBlockSet{Col: 1, Row: 2, Op: op.TypeSet, Block: block.Rock},
				&op.FieldBlockSet{Col: 0, Row: 0, Op: op.TypeSet, Block: block.Rock},
				&op.FieldBlockSet{Col: 1, Row: 3, Op: op.TypeSet, Block: block.Rock},
				&op.FieldBlockSet{Col: 0, Row: 1, Op: op.TypeSet, Block: block.Rock},
			},
		},
		{
			name: "lava: height 0, 1 and 2",
			x:    0, y: 3, p: piece.NewTetromino(piece.TetrominoI, block.Block{Type: block.TypeLava}),
			blocks: []block.XYB{
				{XY: block.XY{X: 0, Y: 0}, Block: block.Rock},
				{XY: block.XY{X: 1, Y: 1}, Block: block.Rock},
			},
			events: []event.Event{
				&op.FieldBlockSet{Col: 3, Row: 0, Op: op.TypeSet, Block: block.Rock},
				&op.FieldBlockSet{Col: 2, Row: 0, Op: op.TypeSet, Block: block.Rock},
				&op.FieldBlockSet{Col: 1, Row: 2, Op: op.TypeSet, Block: block.Rock},
				&op.FieldBlockSet{Col: 0, Row: 1, Op: op.TypeSet, Block: block.Rock},
			},
		},
		{
			name: "acid: empty and wall",
			x:    0, y: 3, p: piece.NewTetromino(piece.TetrominoO, block.Block{Type: block.TypeAcid}),
			blocks: []block.XYB{
				{XY: block.XY{X: 1, Y: 1}, Block: block.Wall},
			},
			events: []event.Event{
				&op.FieldExBlock{Col: 1, Row: 1, AnimType: field.AnimFall, AnimParam: 1, Block: block.Acid},
				&op.FieldExBlock{Col: 0, Row: 0, AnimType: field.AnimFall, AnimParam: 2, Block: block.Acid},
				&op.FieldExBlock{Col: 1, Row: 1, AnimType: field.AnimFall, AnimParam: 2, Block: block.Acid},
				&op.FieldExBlock{Col: 0, Row: 0, AnimType: field.AnimFall, AnimParam: 3, Block: block.Acid},
			},
		},
		{
			name: "acid: single and double block",
			x:    0, y: 3, p: piece.NewTetromino(piece.TetrominoO, block.Block{Type: block.TypeAcid}),
			blocks: []block.XYB{
				{XY: block.XY{X: 0, Y: 0}, Block: block.Rock},
				{XY: block.XY{X: 1, Y: 0}, Block: block.Rock},
				{XY: block.XY{X: 1, Y: 1}, Block: block.Rock},
			},
			events: []event.Event{
				&op.FieldExBlock{Col: 1, Row: 1, AnimType: field.AnimFall, AnimParam: 1, Block: block.Acid},
				&op.FieldBlockSet{Col: 1, Row: 1, Op: op.TypeClear, Block: block.Rock},
				&op.FieldExBlock{Col: 0, Row: 0, AnimType: field.AnimFall, AnimParam: 2, Block: block.Acid},
				&op.FieldBlockSet{Col: 0, Row: 0, Op: op.TypeClear, Block: block.Rock},
				&op.FieldExBlock{Col: 1, Row: 0, AnimType: field.AnimFall, AnimParam: 3, Block: block.Acid},
				&op.FieldBlockSet{Col: 1, Row: 0, Op: op.TypeClear, Block: block.Rock},
				&op.FieldExBlock{Col: 0, Row: 0, AnimType: field.AnimFall, AnimParam: 3, Block: block.Acid},
			},
		},
		{
			name: "acid: hard on top and at bottom",
			x:    0, y: 3, p: piece.NewTetromino(piece.TetrominoO, block.Block{Type: block.TypeAcid}),
			blocks: []block.XYB{
				{XY: block.XY{X: 0, Y: 0}, Block: block.Rock},
				{XY: block.XY{X: 0, Y: 1}, Block: block.Hard},
				{XY: block.XY{X: 1, Y: 0}, Block: block.Hard},
				{XY: block.XY{X: 1, Y: 1}, Block: block.Rock},
			},
			events: []event.Event{
				&op.FieldExBlock{Col: 1, Row: 1, AnimType: field.AnimFall, AnimParam: 1, Block: block.Acid},
				&op.FieldBlockSet{Col: 1, Row: 1, Op: op.TypeClear, Block: block.Rock},
				&op.FieldExBlock{Col: 0, Row: 1, AnimType: field.AnimFall, AnimParam: 1, Block: block.Acid},
				&op.FieldBlockHardness{Col: 0, Row: 1, Hardness: -1},
				&op.FieldExBlock{Col: 1, Row: 0, AnimType: field.AnimFall, AnimParam: 3, Block: block.Acid},
				&op.FieldBlockHardness{Col: 1, Row: 0, Hardness: -1},
				&op.FieldExBlock{Col: 0, Row: 1, AnimType: field.AnimFall, AnimParam: 2, Block: block.Acid},
				&op.FieldBlockSet{Col: 0, Row: 1, Op: op.TypeClear, Block: block.Rock},
			},
		},
	}

	nextPieceEvents := []event.Event{
		&op.PieceSet{Op: op.TypeClear, X: 0, Y: 3},
		&op.PieceState{OldState: piece.StateSlide, NewState: piece.StateNew},
	}

	opts := cmp.Options{
		cmpopts.IgnoreFields(block.Block{}, "Color"),
		cmpopts.IgnoreFields(op.FieldBlockSet{}, "AnimType", "AnimParam"),
		cmpopts.IgnoreFields(op.FieldBlockHardness{}, "AnimType", "AnimParam"),
		cmpopts.IgnoreFields(op.PieceSet{}, "Piece"),
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.events = append(test.events, nextPieceEvents...)

			f := field.Make(6, 6, 1)
			ctrl := f.Ctrl(0)

			ctrl.SetXYP(test.x, test.y, test.p)
			ctrl.Blocks = piece.GetBlocks(test.p, nil)
			ctrl.State = piece.StateSlide

			for _, xyb := range test.blocks {
				f.SetXY(xyb.X, xyb.Y, field.AnimNo, 0, xyb.Block)
			}

			var events event.Slice

			HandleTimeout(f, ctrl, &events)

			if diff := cmp.Diff(test.events, ([]event.Event)(events), opts); diff != "" {
				t.Errorf("test '%s' failed. diff=%s",
					test.name, diff)
			}
		})
	}
}
