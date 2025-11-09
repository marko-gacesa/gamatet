// Copyright (c) 2020, 2025 by Marko Gaćeša

package field

import (
	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/piece"
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		w, h, p int
	}{
		{MinWidth, MinHeight, 1},
		{MaxWidth, MaxHeight, 2},
	}

	for _, test := range tests {
		f := Make(test.w, test.h, test.p)

		if len(f.pieces) != test.p {
			t.Errorf("test %dx%d %dP failed: len(f.pieces)=%d", test.w, test.h, test.p, len(f.pieces))
			continue
		}

		if f.w != test.w || f.h != test.h {
			t.Errorf("test %dx%d %dP failed: f.w=%d f.h=%d", test.w, test.h, test.p, f.w, f.h)
			continue
		}

		if len(f.blocks) != test.w*test.h {
			t.Errorf("test %dx%d %dP failed: len(f.blocks)=%d", test.w, test.h, test.p, len(f.blocks))
			continue
		}

		for i := 0; i < test.p; i++ {
			if f.Ctrl(byte(i)).Idx != i {
				t.Errorf("test %dx%d %dP failed: Ctrl(%d).Idx=%d", test.w, test.h, test.p, i, f.Ctrl(byte(i)).Idx)
				continue
			}
		}
	}
}

func TestField_fill(t *testing.T) {
	f := Make(6, 6, 0)
	f.fill(block.Block{Type: block.TypeRock})

	for y := 0; y < f.h; y++ {
		for x := 0; x < f.w; x++ {
			if f._isXYEmpty(x, y, 0, f.w-1, false, -1) {
				t.Errorf("test failed for (x, y)=(%d, %d)", x, y)
				return
			}
		}
	}
}

func TestField_clear(t *testing.T) {
	f := Make(6, 6, 0)
	f.fill(block.Block{Type: block.TypeRock})
	f.clear()

	for y := 0; y < f.h; y++ {
		for x := 0; x < f.w; x++ {
			if !f._isXYEmpty(x, y, 0, f.w-1, false, -1) {
				t.Errorf("test failed for (x, y)=(%d, %d)", x, y)
				return
			}
		}
	}
}

func TestField_getXYPieceIdx(t *testing.T) {
	f := Make(6, 6, 4)

	p0 := piece.NewStandardTetromino(piece.TetrominoI)
	p1 := p0.Clone()
	p1.UndoActivate() // rotate CW

	// calculate empty x, y offsets for
	p0x := int(p0.LeftEmptyColumns())
	p0y := int(p0.TopEmptyRows())
	p1x := int(p1.LeftEmptyColumns())
	p1y := int(p1.TopEmptyRows())

	f.Ctrl(0).SetXYP(0-p0x, 0+p0y, p0)
	f.Ctrl(1).SetXYP(2-p0x, 5+p0y, p0)
	f.Ctrl(2).SetXYP(3-p1x, 4+p1y, p1)
	f.Ctrl(3).SetXYP(0-p1x, 5+p1y, p1)

	// 5 3 . 1 1 1 1
	// 4 3 . . 2 . .
	// 3 3 . . 2 . .
	// 2 3 . . 2 . .
	// 1 . . . 2 . .
	// 0 0 0 0 0 . .
	//   0 1 2 3 4 5

	tests := []struct {
		x, y, exp int
	}{
		{0, 0, 0}, {3, 0, 0}, {4, 0, -1},
		{1, 5, -1}, {2, 5, 1}, {5, 5, 1},
		{3, 1, 2}, {3, 4, 2},
		{0, 1, -1}, {0, 2, 3}, {0, 5, 3},
	}

	for _, test := range tests {
		r := f._getXYPieceIdx(test.x, test.y)
		if r != test.exp {
			t.Errorf("test failed for (%d, %d); expected %d got %d", test.x, test.y, test.exp, r)
		}
	}
}

func TestField_isXYEmpty(t *testing.T) {
	f := Make(6, 6, 2)

	b := block.Block{Type: block.TypeRock}
	f.setXY(0, 0, b)
	f.setXY(1, 1, b)

	p0 := piece.NewStandardTetromino(piece.TetrominoI)
	p0x := int(p0.LeftEmptyColumns())
	p0y := int(p0.TopEmptyRows())
	f.Ctrl(0).SetXYP(2-p0x, 2+p0y, p0)
	f.Ctrl(1).SetXYP(2-p0x, 1+p0y, p0)

	// 5 . . . . . .
	// 4 . . . . . .
	// 3 . . . . . .
	// 2 . . 0 0 0 0
	// 1 . # 1 1 1 1
	// 0 # . . . . .
	//   0 1 2 3 4 5

	tests := []struct {
		x, y           int
		colMin, colMax int
		liftAll        bool
		liftIdx        int
		expIsEmpty     bool
	}{
		// row 0: block at 0, 0 is in the way
		{0, 0, 0, f.w - 1, false, -1, false},
		{1, 0, 0, f.w - 1, false, -1, true},
		// row 1: block at 1, 1 and piece 1 are in the way
		{2, 1, 0, f.w - 1, false, -1, false},
		{2, 1, 0, f.w - 1, false, 1, true},
		{2, 1, 0, f.w - 1, true, -1, true},
		// row 2: piece 0 is in the way
		{2, 2, 0, f.w - 1, false, -1, false},
		{2, 2, 0, f.w - 1, false, 0, true},
		{2, 2, 0, f.w - 1, true, -1, true},
		// row 3: limited by col min and max
		{1, 3, 2, 4, false, -1, false},
		{2, 3, 2, 4, false, -1, true},
		{5, 3, 2, 4, false, -1, false},
		{4, 3, 2, 4, false, -1, true},
	}

	for _, test := range tests {
		r := f._isXYEmpty(test.x, test.y, test.colMin, test.colMax, test.liftAll, test.liftIdx)
		if r != test.expIsEmpty {
			t.Errorf("test failed for (%d, %d, %t, %d); expected %t got %t", test.x, test.y, test.liftAll, test.liftIdx, test.expIsEmpty, r)
		}
	}
}

func TestField_canPlacePiece(t *testing.T) {
	f := Make(6, 6, 1)

	b := block.Block{Type: block.TypeRock}
	f.setXY(2, 2, b)
	f.setXY(5, 5, b)

	p0 := piece.NewStandardTetromino(piece.TetrominoT)
	p0x := int(p0.LeftEmptyColumns())
	p0y := int(p0.TopEmptyRows())
	f.Ctrl(0).SetXYP(2-p0x, 2+p0y, p0)

	// 5 . . . . . #
	// 4 . . . . . .
	// 3 . . . . . .
	// 2 . . # 0 . .
	// 1 . . 0 0 0 .
	// 0 . . . . . .
	//   0 1 2 3 4 5

	tests := []struct {
		name    string
		px, py  int
		liftAll bool
		liftIdx int
		exp     bool
	}{
		{"can fit in the lower left", 0, 1, false, -1, true},
		{"can fit in the upper right", 3, 5, false, -1, true},
		{"blocked by the block at (2, 2)", 0, 3, false, -1, false},
		{"blocked by the right edge", 4, 1, false, -1, false},
		{"blocked by the bottom edge", 2, 0, false, -1, false},
		{"blocked by the top edge", 1, 6, false, -1, false},
		{"blocked by the left edge", -1, 4, false, -1, false},
		{"blocked by the piece 0", 0, 2, false, -1, false},
		{"can fit when the piece 0 is up", 0, 2, false, 0, true},
		{"can fit when all pieces are up", 0, 2, true, -1, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := f._canPlacePiece(test.px, test.py, 0, f.w-1, p0, test.liftAll, test.liftIdx)
			if r != test.exp {
				t.Errorf("test '%s' failed, for (%d, %d, %t, %d); expected %t got %t", test.name, test.px, test.py, test.liftAll, test.liftIdx, test.exp, r)
			}
		})
	}
}
