// Copyright (c) 2020-2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package field

import (
	"reflect"
	"testing"

	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/piece"
)

func TestField_GetRow(t *testing.T) {
	f := Make(6, 6, 0)

	for row := range 6 {
		for column := range 6 {
			b := block.Block{Type: block.TypeRock, Hardness: byte(row + column)}
			f.setXY(column, row, b)
		}
	}

	blocks := f.GetRow(2)

	for column := range 6 {
		b := blocks[column]
		if b.Type != block.TypeRock {
			t.Errorf("test failed, at blocks index %d; expected Type %d got %d", column, block.TypeRock, b.Type)
		}
		if int(b.Hardness) != 2+column {
			t.Errorf("test failed, at blocks index %d; expected Hardness %d got %d", column, 2+column, b.Hardness)
		}
	}
}

func TestField_CanMovePiece(t *testing.T) {
	f := Make(6, 6, 2)

	b := block.Block{Type: block.TypeRock}
	f.setXY(2, 2, b)

	f.Ctrl(0).SetXYP(2, 2, piece.NewStandardTetromino(piece.TetrominoT))
	f.Ctrl(1).SetXYP(0, 1, piece.NewStandardTetromino(piece.TetrominoO))

	// 5 . . . . . .
	// 4 . . . . . .
	// 3 . . . . . .
	// 2 . . # 0 . .
	// 1 1 1 0 0 0 .
	// 0 1 1 . . . .
	//   0 1 2 3 4 5

	tests := []struct {
		name         string
		dx, dy, pIdx int
		liftAll      bool
		exp          bool
	}{
		{"P0 can't move up", 0, 1, 0, false, false},
		{"P0 can't move left", -1, 0, 0, false, false},
		{"P0 can move down", 0, -1, 0, false, true},
		{"P0 can move right", 1, 0, 0, false, true},
		{"P1 can move up", 0, 1, 1, false, true},
		{"P1 can't move left", -1, 0, 1, false, false},
		{"P1 can't move down", 0, -1, 1, false, false},
		{"P1 can't move right", 1, 0, 1, false, false},
		{"P1 can move right with no collisions", 1, 0, 1, true, true},
	}

	for _, test := range tests {
		r := f.CanMovePiece(test.dx, test.dy, test.pIdx, test.liftAll)
		if r != test.exp {
			t.Errorf("test '%s' failed, for (%d, %d); expected %t got %t", test.name, test.dx, test.dy, test.exp, r)
		}
	}
}

func TestField_CanRotatePiece(t *testing.T) {
	f := Make(6, 6, 3)

	b := block.Block{Type: block.TypeRock}
	f.setXY(5, 2, b)

	p0 := piece.NewStandardTetromino(piece.TetrominoJ)
	p1 := piece.NewStandardTetromino(piece.TetrominoT)
	p0x := int(p0.LeftEmptyColumns())
	p0y := int(p0.TopEmptyRows())
	p1x := int(p1.LeftEmptyColumns())
	p1y := int(p1.TopEmptyRows())

	f.Ctrl(0).SetXYP(3-p0x, 2+p0y, p0)
	f.Ctrl(1).SetXYP(2-p1x, 4+p1y, p1)
	f.Ctrl(2).SetXYP(0-p1x, 1+p1y, p1)

	// 5 . . . . . .
	// 4 . . . 1 . .
	// 3 . . 1 1 1 .
	// 2 . . . 0 . #
	// 1 . 2 . 0 0 0
	// 0 2 2 2 . . .
	//   0 1 2 3 4 5

	tests := []struct {
		name    string
		cw      bool
		pIdx    int
		liftAll bool
		exp     bool
	}{
		{"P0 can't rotate CW", true, 0, false, false},
		{"P0 can rotate CCW", false, 0, false, true},
		{"P1 can't rotate CW", true, 1, false, false},
		{"P1 can't rotate CCW", false, 1, false, false},
		{"P1 can't rotate CW no collisions", true, 1, true, true},
		{"P1 can't rotate CCW no collisions", false, 1, true, true},
		{"P2 can't rotate CW", true, 2, false, false},
		{"P2 can't rotate CCW", false, 2, false, false},
	}

	for _, test := range tests {
		f.Ctrl(byte(test.pIdx)).Config.RotationDirectionCW = test.cw
		r, _, _, _, _ := f.CanRotatePiece(test.pIdx, test.liftAll)
		if r != test.exp {
			t.Errorf("test '%s' failed, for cw=%t; expected %t got %t", test.name, test.cw, test.exp, r)
		}
	}
}

func TestField_CanRotatePiece_WallKick(t *testing.T) {
	//   []
	// [][][]
	p0 := piece.NewStandardTetromino(piece.TetrominoT)

	// []
	// [][]
	// []
	p0CW := p0.Clone()
	p0CW.UndoActivate() // rotate CW

	//   []
	// [][]
	//   []
	p0CCW := p0.Clone()
	p0CCW.Activate() // rotate CCW

	// [][][][]
	p1 := piece.NewStandardTetromino(piece.TetrominoI)

	p1R := p1.Clone()
	p1R.UndoActivate() // rotate CW

	tests := []struct {
		name     string
		wall     byte // 'L'=left, 'R'=right
		p        piece.Piece
		colPad   int
		wallKick byte
		success  bool
		dx       int
	}{
		// the T piece
		{name: "T piece, left wall, no wall kick", wall: 'L', p: p0CW, wallKick: 0, success: false, dx: 0},
		{name: "T piece, left wall, wall kick", wall: 'L', p: p0CW, wallKick: 1, success: true, dx: 1},
		{name: "T piece, right wall, no wall kick", wall: 'R', p: p0CCW, wallKick: 0, success: false, dx: 0},
		{name: "T piece, right wall, wall kick", wall: 'R', p: p0CCW, wallKick: 1, success: true, dx: -1},
		// the I piece
		{name: "I piece, left wall, no wall kick", wall: 'L', p: p1R, wallKick: 0, success: false, dx: 0},
		{name: "I piece, left wall, wall kick", wall: 'L', p: p1R, wallKick: 1, success: true, dx: 1},
		{name: "I piece, right wall, no wall kick", wall: 'R', p: p1R, wallKick: 0, success: false, dx: 0},
		{name: "I piece, right wall, wall kick 1", wall: 'R', p: p1R, wallKick: 1, success: false, dx: 0},
		{name: "I piece, right wall, wall kick 2", wall: 'R', p: p1R, wallKick: 2, success: true, dx: -2},
		// the I piece with column limit
		{name: "I piece, left wall, no wall kick; col lim", wall: 'L', p: p1R, colPad: 1, wallKick: 0, success: false, dx: 0},
		{name: "I piece, left wall, wall kick; col lim", wall: 'L', p: p1R, colPad: 1, wallKick: 1, success: true, dx: 1},
		{name: "I piece, right wall, no wall kick; col lim", wall: 'R', p: p1R, colPad: 1, wallKick: 0, success: false, dx: 0},
		{name: "I piece, right wall, wall kick 1; col lim", wall: 'R', p: p1R, colPad: 1, wallKick: 1, success: false, dx: 0},
		{name: "I piece, right wall, wall kick 2; col lim", wall: 'R', p: p1R, colPad: 1, wallKick: 2, success: true, dx: -2},
	}

	for _, test := range tests {
		const dimW = 6
		const dimH = 6

		f := Make(dimW, dimH, 1)

		ctrl := f.Ctrl(0)
		ctrl.Config.WallKick = test.wallKick
		ctrl.Config.RotationDirectionCW = true

		var colMin, colMax int

		if test.colPad > 0 {
			colMin = test.colPad
			colMax = dimW - 1 - test.colPad
			ctrl.IsColumnLimited = true
			ctrl.ColumnLimit.Min = colMin
			ctrl.ColumnLimit.Max = colMax
		} else {
			colMin = 0
			colMax = dimW - 1
		}

		x := colMin
		y := 4
		ctrl.SetXYP(x, y, test.p)

		switch test.wall {
		case 'R':
			for f._canPlacePiece(x+1, y, colMin, colMax, test.p, false, 0) {
				x++
				ctrl.X = x
			}
		case 'L':
			for f._canPlacePiece(x-1, y, colMin, colMax, test.p, false, 0) {
				x--
				ctrl.X = x
			}
		}

		success, _, dx, _, rotated := f.CanRotatePiece(0, false)

		if success {
			ctrl.X += dx
			ctrl.Piece = rotated
		}

		if success != test.success || dx != test.dx {
			t.Errorf("test '%s' failed, expected success=%t got %t, expected dx=%d got %d", test.name, test.success, success, test.dx, dx)
		}
	}
}

func TestField_GetDropHeight(t *testing.T) {
	f := Make(6, 6, 3)

	b := block.Block{Type: block.TypeRock}
	f.setXY(4, 3, b)

	p0 := piece.NewStandardTetromino(piece.TetrominoO)
	p1 := piece.NewStandardTetromino(piece.TetrominoI)

	p0x := int(p0.LeftEmptyColumns())
	p0y := int(p0.TopEmptyRows())
	p1x := int(p1.LeftEmptyColumns())
	p1y := int(p1.TopEmptyRows())

	f.Ctrl(0).SetXYP(0-p0x, 3+p0y, p0)
	f.Ctrl(1).SetXYP(0-p1x, 4+p1y, p1)
	f.Ctrl(2).SetXYP(1-p1x, 5+p1y, p1)

	// 5 . 2 2 2 2 .
	// 4 1 1 1 1 . .
	// 3 0 0 . . # .
	// 2 0 0 . . . .
	// 1 . . . . . .
	// 0 . . . . . .
	//   0 1 2 3 4 5

	tests := []struct {
		name string
		pIdx int
		exp  int
	}{
		{"P0 drop", 0, 2},
		{"P0 drop, again", 0, 0},
		{"P1 drop", 1, 2},
		{"P2 drop", 2, 1},
	}

	for _, test := range tests {
		r := f.GetDropHeight(test.pIdx, false)
		if r != test.exp {
			t.Errorf("test '%s' failed; expected %d got %d", test.name, test.exp, r)
		}

		f.pieces[test.pIdx].Y -= r
	}
}

func TestField_GetDropHeightFlipV(t *testing.T) {
	f := Make(6, 6, 2)

	b := block.Block{Type: block.TypeRock}
	f.setXY(4, 2, b)

	f.Ctrl(0).SetXYP(0, 3, piece.NewFlipVTetromino(0, b))
	f.Ctrl(1).SetXYP(1, 4, piece.NewFlipVTetromino(1, b))

	// 5 . 2 2 2 2 .
	// 4 . 1 1 1 1 .
	// 3 0 0 . . . .
	// 2 0 0 . . # .
	// 1 . . . . . .
	// 0 . . . . . .
	//   0 1 2 3 4 5

	if want, got := 2, f.GetDropHeight(0, false); want != got {
		t.Errorf("GetDropHeight(0,false) = %d, want %d", got, want)
	}
	if want, got := 2, f.GetDropHeight(0, true); want != got {
		t.Errorf("GetDropHeight(0,true) = %d, want %d", got, want)
	}
	if want, got := 0, f.GetDropHeight(1, false); want != got {
		t.Errorf("GetDropHeight(1,false) = %d, want %d", got, want)
	}
	if want, got := 1, f.GetDropHeight(1, true); want != got {
		t.Errorf("GetDropHeight(1,true) = %d, want %d", got, want)
	}
}

func TestField_GetPieceBlockLocations(t *testing.T) {
	f := Make(6, 6, 0)

	b := block.Block{Type: block.TypeRock}
	f.setXY(2, 1, b)

	const p0Color = 153459
	const p1Color = 7392
	const p2Color = 45674233

	b0 := block.Block{Type: block.TypeRock, Color: p0Color}
	b1 := block.Block{Type: block.TypeRock, Color: p1Color}
	b2 := block.Block{Type: block.TypeRock, Color: p2Color}

	p0 := piece.NewTetromino(0, b0)
	p1 := piece.NewTetromino(0, b1)
	p2 := piece.NewTetromino(0, b2)

	// 5 . . . . . .
	// 4 . . . . . .
	// 3 . . 2 2 . .
	// 2 . 1 1 2 . .
	// 1 . 1 # . 0 0
	// 0 . . . . 0 0
	//   0 1 2 3 4 5

	tests := []struct {
		name string
		x, y int
		p    piece.Piece
		exp  []block.XYB
	}{
		{"P0 meld", 4, 1, p0, []block.XYB{{block.XY{4, 0}, b0}, {block.XY{5, 0}, b0}, {block.XY{4, 1}, b0}, {block.XY{5, 1}, b0}}},
		{"P1 meld", 1, 2, p1, []block.XYB{{block.XY{1, 1}, b1}, {block.XY{1, 2}, b1}, {block.XY{2, 2}, b1}}},
		{"P2 meld", 2, 3, p2, []block.XYB{{block.XY{3, 2}, b2}, {block.XY{2, 3}, b2}, {block.XY{3, 3}, b2}}},
	}

	for _, test := range tests {
		blockLocations := f.GetPieceBlockLocations(test.x, test.y, test.p)
		if reflect.DeepEqual(blockLocations, test.exp) {
			t.Errorf("test %q failed; expected=%+v, got=%+v", test.name, test.exp, blockLocations)
		}
	}
}

func TestField_GetPieceStartPosition(t *testing.T) {
	const w = 10
	const h = 8
	f := Make(w, h, 2)

	widthPerPiece := f.w / 2

	p0 := piece.NewStandardTetromino(piece.TetrominoO)
	p1 := piece.NewStandardTetromino(piece.TetrominoT)

	dimP0X := int(p0.DimX())
	dimP1X := int(p1.DimX())

	tests := []struct {
		name       string
		pIdx       int
		piece      piece.Piece
		expX, expY int
	}{
		{"P0@0", 0, p0, (widthPerPiece - dimP0X) / 2, h - 1},
		{"P0@1", 1, p0, widthPerPiece + (widthPerPiece-dimP0X)/2, h - 1},
		{"P1@0", 0, p1, (widthPerPiece - dimP1X) / 2, h - 1},
		{"P1@1", 1, p1, widthPerPiece + (widthPerPiece-dimP1X)/2, h - 1},
	}

	for _, test := range tests {
		f.pieces[test.pIdx].Piece = nil
		f.pieces[test.pIdx].X = 0
		f.pieces[test.pIdx].Y = 0

		success, x, y := f.GetPieceStartPosition(test.pIdx, f.Ctrl(byte(test.pIdx)), test.piece, false)

		if success == false {
			t.Errorf("test '%s' failed; expected success, but it failed", test.name)
			continue
		}

		f.pieces[test.pIdx].X = x
		f.pieces[test.pIdx].Y = y
		f.pieces[test.pIdx].Piece = test.piece

		if f.pieces[test.pIdx].X != test.expX || f.pieces[test.pIdx].Y != test.expY {
			t.Errorf("test '%s' failed; expected (x, y)=(%d, %d), but got (%d, %d)", test.name, test.expX, test.expY, f.pieces[test.pIdx].X, f.pieces[test.pIdx].Y)
		}
	}
}

func TestField_GetPieceStartPosition2(t *testing.T) {
	const w = 12
	const h = 4
	f := Make(w, h, 3)
	f.setXY(0, 3, block.Block{Type: block.TypeRock})

	p0 := piece.NewStandardTetromino(piece.TetrominoO)
	p1 := piece.NewStandardTetromino(piece.TetrominoI)
	p2 := piece.NewStandardTetromino(piece.TetrominoI)
	p2.UndoActivate() // rotate CW

	p0y := int(p0.TopEmptyRows())
	p1y := int(p1.TopEmptyRows())
	p2y := int(p2.TopEmptyRows())

	// 3 # 0 0 . 1 1 1 1 . . 2 .
	// 2 . 0 0 . . . . . . . 2 .
	// 1 . . . . . . . . . . 2 .
	// 0 . . . . . . . . . . 2 .
	//   0 1 2 3 4 5 6 7 8 9 0 1

	tests := []struct {
		name       string
		piece      piece.Piece
		pIdx       int
		expSuccess bool
		expX, expY int
		liftAll    bool
		moveP1dX   int
	}{
		{"P1@0", p1, 0, false, 0, 0, false, 0}, // can't put the p1 at 0 because of the block at (0,3)
		{"P0@0", p0, 0, true, 1, 3 + p0y, false, 0},
		{"P1@1", p1, 1, true, 4, 3 + p1y, false, 0},
		{"P2@2", p2, 2, true, 8, 3 + p2y, false, 0},
		{"P2@2 over P1", p2, 2, false, 0, 0, false, 4}, // moves p1 to position 2: x+=4
		{"P2@2 over P1 lift all", p2, 2, true, 8, 3 + p2y, true, 0},
	}

	for _, test := range tests {
		f.pieces[test.pIdx].Piece = nil
		f.pieces[test.pIdx].X = 0
		f.pieces[test.pIdx].Y = 0
		if test.moveP1dX != 0 && f.pieces[1].Piece != nil {
			f.pieces[1].X += test.moveP1dX
		}

		success, x, y := f.GetPieceStartPosition(test.pIdx, f.Ctrl(byte(test.pIdx)), test.piece, test.liftAll)

		if success == false {
			if test.expSuccess {
				t.Errorf("test '%s' failed; expected success, but it failed", test.name)
			}
			continue
		}

		f.pieces[test.pIdx].X = x
		f.pieces[test.pIdx].Y = y
		f.pieces[test.pIdx].Piece = test.piece

		if f.pieces[test.pIdx].X != test.expX || f.pieces[test.pIdx].Y != test.expY {
			t.Errorf("test '%s' failed; expected (x, y)=(%d, %d), but got (%d, %d)", test.name, test.expX, test.expY, f.pieces[test.pIdx].X, f.pieces[test.pIdx].Y)
		}
	}
}

func TestField_GetTopmostEmpty(t *testing.T) {
	f := Make(6, 6, 0)
	f.setXY(0, 2, block.Block{Type: block.TypeRock})
	f.setXY(1, 5, block.Block{Type: block.TypeRock})

	// 5 . # . . . .
	// 4 . . . . . .
	// 3 . . . . . .
	// 2 # . . . . .
	// 1 . . . . . .
	// 0 . . . . . .
	//   0 1 2 3 4 5

	tests := []struct {
		name   string
		col    int
		expRow int
	}{
		{
			name:   "mid",
			col:    0,
			expRow: 3,
		},
		{
			name:   "top",
			col:    1,
			expRow: 6,
		},
		{
			name:   "empty",
			col:    2,
			expRow: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			row := f.GetTopmostEmpty(test.col)
			if row != test.expRow {
				t.Errorf("failed; expected row %d, got %d", test.expRow, row)
			}
		})
	}
}

func TestField_GetHeightToTopmostEmpty(t *testing.T) {
	const column = 0
	const height = 8

	tests := []struct {
		name     string
		y        int
		column   []block.Block
		expected int
	}{
		{
			name:     "fall to the bottom",
			y:        4,
			expected: 4,
		},
		{
			name:     "fall to block at row 0",
			y:        4,
			column:   []block.Block{{Type: block.TypeRock}},
			expected: 3,
		},
		{
			name:     "fall to block at row 2",
			y:        4,
			column:   []block.Block{{Type: block.TypeRock}, {Type: block.TypeRock}, {Type: block.TypeRock}},
			expected: 1,
		},
		{
			name:     "fall blocked",
			y:        4,
			column:   []block.Block{{Type: block.TypeRock}, {Type: block.TypeRock}, {Type: block.TypeRock}, {Type: block.TypeRock}},
			expected: 0,
		},
	}

	for _, test := range tests {
		f := Make(4, height, 0)

		for row, b := range test.column {
			f.setXY(column, row, b)
		}

		result := f.GetHeightToTopmostEmpty(column, test.y)

		if result != test.expected {
			t.Errorf("test %q failed: expected height=%d but got height=%d", test.name, test.expected, result)
		}
	}
}

func TestField_GetHeightToTopmostFull(t *testing.T) {
	const column = 0
	const height = 8

	tests := []struct {
		name     string
		y        int
		column   []block.Block
		expected int
	}{
		{
			name:     "no block below",
			y:        4,
			expected: 0,
		},
		{
			name:     "block at row 0",
			y:        4,
			column:   []block.Block{{Type: block.TypeRock}},
			expected: 4,
		},
		{
			name:     "all full below",
			y:        4,
			column:   []block.Block{{Type: block.TypeRock}, {Type: block.TypeRock}, {Type: block.TypeRock}, {Type: block.TypeRock}},
			expected: 1,
		},
	}

	for _, test := range tests {
		f := Make(4, height, 0)

		for row, b := range test.column {
			f.setXY(column, row, b)
		}

		result := f.GetHeightToTopmostFull(column, test.y)

		if result != test.expected {
			t.Errorf("test %q failed: expected height=%d but got height=%d", test.name, test.expected, result)
		}
	}
}

func TestField_GetHeightToHighestHole(t *testing.T) {
	const column = 0
	const height = 8

	bf := block.Block{Type: block.TypeRock}
	be := block.Block{Type: block.TypeEmpty}

	tests := []struct {
		name     string
		y        int
		column   []block.Block
		expected int
	}{
		{
			name:     "no block below",
			y:        4,
			expected: 0,
		},
		{
			name:     "block at row 0",
			y:        4,
			column:   []block.Block{bf},
			expected: 0,
		},
		{
			name:     "block at row 3",
			y:        4,
			column:   []block.Block{be, be, be, bf},
			expected: 2,
		},
		{
			name:     "blocks at row 1 and 3",
			y:        4,
			column:   []block.Block{be, bf, be, bf},
			expected: 2,
		},
		{
			name:     "block at row 2",
			y:        4,
			column:   []block.Block{be, be, bf},
			expected: 3,
		},
		{
			name:     "all full below",
			y:        4,
			column:   []block.Block{bf, bf, bf, bf},
			expected: 0,
		},
	}

	for _, test := range tests {
		f := Make(4, height, 0)

		for row, b := range test.column {
			f.setXY(column, row, b)
		}

		result := f.GetHeightToHighestHole(column, test.y)

		if result != test.expected {
			t.Errorf("test %q failed: expected height=%d but got height=%d", test.name, test.expected, result)
		}
	}
}

func TestField_GetHeightToLowestHole(t *testing.T) {
	const column = 2
	const height = 10

	bf := block.Block{Type: block.TypeRock}
	be := block.Block{Type: block.TypeEmpty}

	tests := []struct {
		name     string
		y        int
		column   []block.Block
		expected int
	}{
		{
			name:     "no block below",
			y:        4,
			expected: 0,
		},
		{
			name:     "block at row 0",
			y:        4,
			column:   []block.Block{bf},
			expected: 0,
		},
		{
			name:     "block at row 3",
			y:        4,
			column:   []block.Block{be, be, be, bf},
			expected: 4,
		},
		{
			name:     "blocks at row 1 and 3",
			y:        4,
			column:   []block.Block{be, bf, be, bf},
			expected: 4,
		},
		{
			name:     "block at row 0 and 2",
			y:        4,
			column:   []block.Block{bf, be, bf},
			expected: 3,
		},
		{
			name:     "block at row 0, 2, 3, 5",
			y:        9,
			column:   []block.Block{bf, be, bf, bf, be, bf},
			expected: 8,
		},
		{
			name:     "all full below",
			y:        4,
			column:   []block.Block{bf, bf, bf, bf},
			expected: 0,
		},
	}

	for _, test := range tests {
		f := Make(4, height, 0)

		for row, b := range test.column {
			f.setXY(column, row, b)
		}

		result := f.GetHeightToLowestHole(column, test.y)

		if result != test.expected {
			t.Errorf("test %q failed: expected height=%d but got height=%d", test.name, test.expected, result)
		}
	}
}

func TestField_GetDestroyInfo(t *testing.T) {
	w := 4
	h := 8

	type exp struct {
		Row    int
		Height int
		N      int
		Type   block.Type
	}

	tests := []struct {
		name        string
		rowsToClear []int
		columnDef   []block.Block
		expCount    int
		expected    []exp
		expHardDec  []int
		expHasImm   bool
	}{
		{
			name:        "simple test: do nothing",
			rowsToClear: []int{},
			columnDef:   []block.Block{{Type: block.TypeRock}},
			expected:    []exp{},
		},
		{
			name:        "simple test: destroy row 0",
			rowsToClear: []int{0},
			columnDef:   []block.Block{{Type: 42}},
			expCount:    1,
			expected:    []exp{{Row: 0, Height: 1, N: 7, Type: 42}},
		},
		{
			name:        "destroy 3 adjacent rows: 2, 3 and 4",
			rowsToClear: []int{2, 3, 4},
			columnDef:   []block.Block{{Type: block.TypeEmpty}, {Type: block.TypeEmpty}, {Type: 101}, {Type: 102}, {Type: 103}},
			expCount:    3,
			expected:    []exp{{Row: 2, Height: 1, N: 0, Type: 101}, {Row: 3, Height: 2, N: 0, Type: 102}, {Row: 4, Height: 3, N: 3, Type: 103}},
		},
		{
			name:        "destroy 2 nonadjacent rows: 1 and 4",
			rowsToClear: []int{1, 4},
			columnDef:   []block.Block{{Type: block.TypeEmpty}, {Type: 42}, {Type: block.TypeEmpty}, {Type: block.TypeEmpty}, {Type: 66}},
			expCount:    2,
			expected:    []exp{{Row: 1, Height: 1, N: 2, Type: 42}, {Row: 4, Height: 2, N: 3, Type: 66}},
		},
		{
			name:        "destroy 3 nonadjacent rows: 0, 3 and 5",
			rowsToClear: []int{0, 3, 5},
			columnDef:   []block.Block{{Type: 42}, {Type: block.TypeEmpty}, {Type: block.TypeEmpty}, {Type: 66}, {Type: block.TypeEmpty}, {Type: 102}},
			expCount:    3,
			expected:    []exp{{Row: 0, Height: 1, N: 2, Type: 42}, {Row: 3, Height: 2, N: 1, Type: 66}, {Row: 5, Height: 3, N: 2, Type: 102}},
		},
		{
			name:        "destroy row 2, immovable block at 4",
			rowsToClear: []int{2},
			columnDef:   []block.Block{{Type: block.TypeEmpty}, {Type: block.TypeRock}, {Type: block.TypeRock}, {Type: block.TypeRock}, {Type: block.TypeWall}, {Type: block.TypeRock}},
			expCount:    1,
			expected:    []exp{{Row: 2, Height: 1, N: 1, Type: block.TypeRock}},
			expHasImm:   true,
		},
		{
			name:        "destroy rows 0 and 6, immovable block in between at row 4",
			rowsToClear: []int{0, 6},
			columnDef:   []block.Block{{Type: 42}, {Type: block.TypeEmpty}, {Type: block.TypeEmpty}, {Type: block.TypeEmpty}, {Type: block.TypeWall}, {Type: block.TypeEmpty}, {Type: 66}},
			expCount:    2,
			expected:    []exp{{Row: 0, Height: 1, N: 3, Type: 42}, {Row: 6, Height: 1, N: 1, Type: 66}},
			expHasImm:   true,
		},
		{
			name:        "destroy rows 0, 3 and 5, immovable blocks at row 2 and 4",
			rowsToClear: []int{0, 3, 5},
			columnDef:   []block.Block{{Type: 42}, {Type: block.TypeEmpty}, {Type: block.TypeWall}, {Type: 66}, {Type: block.TypeWall}, {Type: 102}},
			expCount:    3,
			expected:    []exp{{Row: 0, Height: 1, N: 1, Type: 42}, {Row: 3, Height: 1, N: 0, Type: 66}, {Row: 5, Height: 1, N: 2, Type: 102}},
			expHasImm:   true,
		},
		{
			name:        "destroy rows 1, 3 and 6, immovable blocks at row 1 and 5, destroying immovable at row 1",
			rowsToClear: []int{1, 3, 6},
			columnDef:   []block.Block{{Type: block.TypeEmpty}, {Type: block.TypeWall}, {Type: block.TypeEmpty}, {Type: 102}, {Type: block.TypeEmpty}, {Type: block.TypeWall}, {Type: 103}},
			expCount:    3,
			expected:    []exp{{Row: 1, Height: 1, N: 1, Type: block.TypeWall}, {Row: 3, Height: 2, N: 1, Type: 102}, {Row: 6, Height: 1, N: 1, Type: 103}},
			expHasImm:   true,
		},
		{
			name:        "destroy rows 1, 4, hardness=1 as at row 4",
			rowsToClear: []int{1, 4},
			columnDef:   []block.Block{{Type: block.TypeEmpty}, {Type: 102}, {Type: block.TypeEmpty}, {Type: block.TypeEmpty}, {Type: block.TypeRock, Hardness: 1}},
			expCount:    2,
			expected:    []exp{{Row: 1, Height: 1, N: 2, Type: 102}},
			expHardDec:  []int{4},
		},
	}

	for column := range 1 {
		for _, test := range tests {
			f := Make(w, h, 0)

			for _, row := range test.rowsToClear {
				for col := range w {
					f.setXY(col, row, block.Block{Type: block.TypeRock})
				}
			}

			for row, b := range test.columnDef {
				f.setXY(column, row, b)
			}

			//fmt.Println(f.String())

			result := f.GetDestroyInfo()

			if result.RowCount != test.expCount {
				t.Errorf("test %q failed. expected row count=%d, but got row count=%d", test.name, test.expCount, result.RowCount)
				continue
			}

			if len(result.HardDec) != len(test.expHardDec) {
				t.Errorf("test %q failed. expected hardness decrease length=%d, but got length=%d", test.name, len(test.expHardDec), len(result.HardDec))
				continue
			}

			for i := 0; i < len(result.HardDec); i++ {
				hardDec := result.HardDec[i]
				if hardDec.X != column || hardDec.Y != test.expHardDec[i] {
					t.Errorf("test %q failed. expected hardness decrease index=%d at (col,row)=(%d, %d), but got (col, row)=(%d, %d)", test.name, i, column, test.expHardDec[i], hardDec.X, hardDec.Y)
					break
				}
			}

			if result.Columns != nil && result.Columns[column].HasImm != test.expHasImm {
				t.Errorf("test %q failed. expected has-immovable=%t, but got has-immovable=%t", test.name, test.expHasImm, result.Columns[column].HasImm)
				continue
			}

			var r []DestroyBlockInfo
			if result.Columns != nil {
				r = result.Columns[column].Rows
			}

			//fmt.Printf("%+v\n", r)

			if len(r) != len(test.expected) {
				t.Errorf("test %q failed. expected result length=%d, but got length=%d", test.name, len(test.expected), len(r))
				continue
			}

			for i := 0; i < len(r); i++ {
				if r[i].Row != test.expected[i].Row {
					t.Errorf("test %q failed. expected row=%d, but got row=%d at index=%d", test.name, test.expected[i].Row, r[i].Row, i)
				} else if r[i].Height != test.expected[i].Height {
					t.Errorf("test %q failed. expected height=%d, but got height=%d at index=%d", test.name, test.expected[i].Height, r[i].Height, i)
				} else if r[i].N != test.expected[i].N {
					t.Errorf("test %q failed. expected N=%d, but got N=%d at index=%d", test.name, test.expected[i].N, r[i].N, i)
				} else if r[i].Type != test.expected[i].Type {
					t.Errorf("test %q failed. expected type=%d, but got type=%d at index=%d", test.name, test.expected[i].Type, r[i].Type, i)
				}
			}
		}
	}
}

func TestField_GetDestroyInfo2(t *testing.T) {
	w := 4
	h := 4

	bR := block.Block{Type: block.TypeRock}
	bI := block.Block{Type: block.TypeWall}
	bH := block.Block{Type: block.TypeRock, Hardness: 2}
	bW := block.Block{Type: block.TypeRock, Hardness: block.HardnessMax}

	type exp struct {
		Row    int
		Height int
		N      int
		Type   block.Type
	}

	tests := []struct {
		name       string
		blocksDef  []block.XYB
		expCount   int
		expected   [][]exp
		expHardDec []block.XY
		expHasHard []bool
		expHasImm  []bool
	}{
		{
			name: "hard block",
			blocksDef: []block.XYB{
				// |[0][0][0]   | => |            |
				// |[0][0][0][2]|    |[0][0][0][1]|
				{block.XY{0, 1}, bR}, {block.XY{1, 1}, bR}, {block.XY{2, 1}, bR},
				{block.XY{0, 0}, bR}, {block.XY{1, 0}, bR}, {block.XY{2, 0}, bR}, {block.XY{3, 0}, bH},
			},
			expected: [][]exp{
				{{0, 1, 3, bR.Type}},
				{{0, 1, 3, bR.Type}},
				{{0, 1, 3, bR.Type}},
				{},
			},
			expCount:   1,
			expHardDec: []block.XY{{X: 3, Y: 0}},
			expHasHard: []bool{false, false, false, true},
			expHasImm:  []bool{false, false, false, false},
		},
		{
			name: "wall block",
			blocksDef: []block.XYB{
				// |[0][0][0]   | => |            |
				// |[0][0][0][W]|    |[0][0][0][W]|
				{block.XY{0, 1}, bR}, {block.XY{1, 1}, bR}, {block.XY{2, 1}, bR},
				{block.XY{0, 0}, bR}, {block.XY{1, 0}, bR}, {block.XY{2, 0}, bR}, {block.XY{3, 0}, bW},
			},
			expected: [][]exp{
				{{0, 1, 3, bR.Type}},
				{{0, 1, 3, bR.Type}},
				{{0, 1, 3, bR.Type}},
				{},
			},
			expCount:   1,
			expHasHard: []bool{false, false, false, true},
			expHasImm:  []bool{false, false, false, false},
		},
		{
			name: "immovable block",
			blocksDef: []block.XYB{
				// |   [0][0][0]|    |            |
				// |[I]         | => |[I][0][0][0]|
				// |[0][0][0][0]|    |            |
				{block.XY{1, 2}, bR}, {block.XY{2, 2}, bR}, {block.XY{3, 2}, bR},
				{block.XY{0, 1}, bI},
				{block.XY{0, 0}, bR}, {block.XY{1, 0}, bR}, {block.XY{2, 0}, bR}, {block.XY{3, 0}, bR},
			},
			expected: [][]exp{
				{{0, 1, 0, bR.Type}},
				{{0, 1, 3, bR.Type}},
				{{0, 1, 3, bR.Type}},
				{{0, 1, 3, bR.Type}},
			},
			expCount:   1,
			expHasHard: []bool{false, false, false, false},
			expHasImm:  []bool{true, false, false, false},
		},
		{
			name: "only wall blocks",
			blocksDef: []block.XYB{
				// |[W][W][W][W]| => |[W][W][W][W]|
				{block.XY{0, 0}, bW}, {block.XY{1, 0}, bW}, {block.XY{2, 0}, bW}, {block.XY{3, 0}, bW},
			},
			expected:   [][]exp{},
			expCount:   0,
			expHasHard: []bool{},
			expHasImm:  []bool{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := Make(w, h, 0)

			for _, xyb := range test.blocksDef {
				f.setXY(xyb.X, xyb.Y, xyb.Block)
			}

			result := f.GetDestroyInfo()

			if result.RowCount != test.expCount {
				t.Errorf("test %q failed. expected row count=%d, but got row count=%d", test.name, test.expCount, result.RowCount)
				return
			}

			if len(result.HardDec) != len(test.expHardDec) {
				t.Errorf("test %q failed. expected hardness decrease length=%d, but got length=%d", test.name, len(test.expHardDec), len(result.HardDec))
				return
			}

			for i := 0; i < len(result.HardDec); i++ {
				hardDec := result.HardDec[i]
				if hardDec.X != test.expHardDec[i].X || hardDec.Y != test.expHardDec[i].Y {
					t.Errorf("test %q failed. expected hardness decrease index=%d at (col,row)=(%+v), but got (col, row)=(%+v)", test.name, i, test.expHardDec, hardDec)
					break
				}
			}

			if result.Columns != nil {
				for column := 0; column < f.w; column++ {
					if result.Columns[column].HasImm != test.expHasImm[column] {
						t.Errorf("test %q failed. expected has-immovable=%t, but got has-immovable=%t", test.name, test.expHasImm[column], result.Columns[column].HasImm)
						continue
					}
				}
			}

			if result.Columns != nil {
				for column := 0; column < f.w; column++ {
					if result.Columns[column].HasHard != test.expHasHard[column] {
						t.Errorf("test %q failed. expected has-hard=%t, but got has-hard=%t", test.name, test.expHasHard[column], result.Columns[column].HasHard)
						continue
					}
				}
			}

			if result.Columns == nil {
				return
			}

			for column := 0; column < len(result.Columns); column++ {
				r := result.Columns[column].Rows
				if len(r) != len(test.expected[column]) {
					t.Errorf("test %q failed. expected result length=%d, but got length=%d", test.name, len(test.expected[column]), len(r))
					continue
				}

				for row := range r {
					if r[row].Row != test.expected[column][row].Row {
						t.Errorf("test %q failed. expected row=%d, but got row=%d at index=%d", test.name, test.expected[column][row].Row, r[row].Row, row)
					} else if r[row].Height != test.expected[column][row].Height {
						t.Errorf("test %q failed. expected height=%d, but got height=%d at index=%d", test.name, test.expected[column][row].Height, r[row].Height, row)
					} else if r[row].N != test.expected[column][row].N {
						t.Errorf("test %q failed. expected N=%d, but got N=%d at index=%d", test.name, test.expected[column][row].N, r[row].N, row)
					} else if r[row].Type != test.expected[column][row].Type {
						t.Errorf("test %q failed. expected type=%d, but got type=%d at index=%d", test.name, test.expected[column][row].Type, r[row].Type, row)
					}
				}
			}
		})
	}
}
