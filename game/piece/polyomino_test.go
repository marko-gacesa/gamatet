// Copyright (c) 2020 by Marko Gaćeša

package piece

import (
	"gamatet/game/block"
	"testing"
)

func TestPolyomino_Rotate(t *testing.T) {
	const _I = false
	const XX = true
	tests := []struct {
		name   string
		rots   byte
		rot    byte
		data   []bool
		expRot byte
		exp    []bool
	}{
		{
			name:   "1x1",
			rots:   0,
			rot:    0,
			data:   []bool{XX},
			expRot: 0,
			exp:    []bool{XX},
		},
		{
			name: "2x2",
			rots: 2,
			rot:  0,
			data: []bool{
				_I, XX,
				XX, _I,
			},
			expRot: 1,
			exp: []bool{
				XX, _I,
				_I, XX,
			},
		},
		{
			name: "3x3",
			rots: 4,
			rot:  3,
			data: []bool{
				_I, XX, _I,
				XX, XX, XX,
				XX, _I, _I,
			},
			expRot: 0,
			exp: []bool{
				XX, XX, _I,
				_I, XX, XX,
				_I, XX, _I,
			},
		},
		{
			name: "4x4",
			rots: 4,
			rot:  2,
			data: []bool{
				_I, XX, _I, XX,
				XX, XX, XX, XX,
				XX, _I, _I, _I,
				XX, XX, XX, XX,
			},
			expRot: 3,
			exp: []bool{
				XX, XX, XX, _I,
				XX, _I, XX, XX,
				XX, _I, XX, _I,
				XX, _I, XX, XX,
			},
		},
		{
			name: "5x5",
			rots: 4,
			rot:  3,
			data: []bool{
				_I, _I, _I, _I, _I,
				_I, _I, XX, XX, XX,
				XX, XX, XX, _I, _I,
				_I, XX, XX, _I, _I,
				_I, XX, _I, _I, _I,
			},
			expRot: 0,
			exp: []bool{
				_I, _I, XX, _I, _I,
				XX, XX, XX, _I, _I,
				_I, XX, XX, XX, _I,
				_I, _I, _I, XX, _I,
				_I, _I, _I, XX, _I,
			},
		},
	}

	var b block.Block

	for _, test := range tests {
		shape := _initPolyomino(test.rots, test.data)
		piece := &polyomino{
			polyominoShape: shape,
			rot:            test.rot,
			block:          b,
		}

		expShape := _initPolyomino(test.rots, test.exp)

		piece.RotateCW()

		if piece.rot != test.expRot {
			t.Errorf("test %s CW rotate failed. expexted rot=%d, but got %d", test.name, test.expRot, piece.rot)
			continue
		}

		if piece.data != expShape.data {
			t.Errorf("test %s CW rotate failed. expexted blocks=%b, but got %b", test.name, expShape.data, piece.data)
			continue
		}

		piece.RotateCCW()

		if piece.rot != test.rot {
			t.Errorf("test %s CCW rotate failed. expexted rot=%d, but got %d", test.name, test.rot, piece.rot)
			continue
		}

		if piece.data != shape.data {
			t.Errorf("test %s CCW rotate failed. expexted blocks=%b, but got %b", test.name, shape.data, piece.data)
			continue
		}

		if test.rots <= 0 {
			continue
		}

		for i := byte(0); i < test.rots; i++ {
			piece.RotateCW()
		}

		if piece.rot != test.rot {
			t.Errorf("test %s full CW rotate failed. expexted rot=%d, but got %d", test.name, test.rot, piece.rot)
			continue
		}

		if piece.data != shape.data {
			t.Errorf("test %s full CW rotate failed. expexted blocks=%b, but got %b", test.name, shape.data, piece.data)
			continue
		}

		for i := byte(0); i < test.rots; i++ {
			piece.RotateCCW()
		}

		if piece.rot != test.rot {
			t.Errorf("test %s full CCW rotate failed. expexted rot=%d, but got %d", test.name, test.rot, piece.rot)
			continue
		}

		if piece.data != shape.data {
			t.Errorf("test %s full CCW rotate failed. expexted blocks=%b, but got %b", test.name, shape.data, piece.data)
			continue
		}

	}
}

func TestPolyomino_IsRowEmpty(t *testing.T) {
	const _I = false
	const XX = true

	type expected struct {
		emptyRows   []bool
		topEmpty    int
		bottomEmpty int
	}

	tests := []struct {
		name string
		data []bool
		exp  expected
	}{
		{
			name: "1x1",
			data: []bool{XX},
			exp:  expected{[]bool{false}, 0, 0},
		},
		{
			name: "2x2",
			data: []bool{
				XX, XX,
				_I, _I,
			},
			exp: expected{[]bool{false, true}, 0, 1},
		},
		{
			name: "3x3",
			data: []bool{
				_I, _I, _I,
				XX, XX, XX,
				_I, _I, _I,
			},
			exp: expected{[]bool{true, false, true}, 1, 1},
		},
		{
			name: "4x4",
			data: []bool{
				_I, _I, _I, _I,
				XX, XX, XX, XX,
				_I, _I, _I, _I,
				_I, _I, _I, _I,
			},
			exp: expected{[]bool{true, false, true, true}, 1, 2},
		},
		{
			name: "4x4 V",
			data: []bool{
				_I, _I, XX, _I,
				_I, _I, XX, _I,
				_I, _I, XX, _I,
				_I, _I, XX, _I,
			},
			exp: expected{[]bool{false, false, false, false}, 0, 0},
		},
	}

	for _, test := range tests {
		shape := _initPolyomino(4, test.data)
		piece := &polyomino{polyominoShape: shape, rot: 0, block: block.Rock}

		for i := 0; i < piece.DimY(); i++ {
			result := piece.IsRowEmpty(i)
			if result != test.exp.emptyRows[i] {
				t.Errorf("test '%s' failed for row=%d, expected=%t. got=%t", test.name, i, test.exp.emptyRows[i], result)
			}
		}

		result := piece.TopEmptyRows()
		if result != test.exp.topEmpty {
			t.Errorf("test '%s' failed for top empty rows, expected=%d. got=%d", test.name, test.exp.topEmpty, result)
		}

		result = piece.BottomEmptyRows()
		if result != test.exp.bottomEmpty {
			t.Errorf("test '%s' failed for bottom empty rows, expected=%d. got=%d", test.name, test.exp.bottomEmpty, result)
		}
	}
}

func TestPolyomino_IsColumnEmpty(t *testing.T) {
	const _I = false
	const XX = true

	type expected struct {
		emptyColumns []bool
		leftEmpty    int
		rightEmpty   int
	}

	tests := []struct {
		name string
		data []bool
		exp  expected
	}{
		{
			name: "1x1",
			data: []bool{XX},
			exp:  expected{[]bool{false}, 0, 0},
		},
		{
			name: "2x2",
			data: []bool{
				_I, XX,
				_I, XX,
			},
			exp: expected{[]bool{true, false}, 1, 0},
		},
		{
			name: "3x3",
			data: []bool{
				_I, XX, _I,
				_I, XX, _I,
				_I, XX, _I,
			},
			exp: expected{[]bool{true, false, true}, 1, 1},
		},
		{
			name: "4x4",
			data: []bool{
				_I, XX, _I, _I,
				_I, XX, _I, _I,
				_I, XX, _I, _I,
				_I, XX, _I, _I,
			},
			exp: expected{[]bool{true, false, true, true}, 1, 2},
		},
		{
			name: "4x4 H",
			data: []bool{
				_I, _I, _I, _I,
				_I, _I, _I, _I,
				XX, XX, XX, XX,
				_I, _I, _I, _I,
			},
			exp: expected{[]bool{false, false, false, false}, 0, 0},
		},
	}

	for _, test := range tests {
		shape := _initPolyomino(4, test.data)
		piece := &polyomino{polyominoShape: shape, rot: 0, block: block.Rock}

		for i := 0; i < piece.DimX(); i++ {
			result := piece.IsColumnEmpty(i)
			if result != test.exp.emptyColumns[i] {
				t.Errorf("test '%s' failed for column=%d, expected=%t. got=%t", test.name, i, test.exp.emptyColumns[i], result)
			}
		}

		result := piece.LeftEmptyColumns()
		if result != test.exp.leftEmpty {
			t.Errorf("test '%s' failed for left empty columns, expected=%d. got=%d", test.name, test.exp.leftEmpty, result)
		}

		result = piece.RightEmptyColumns()
		if result != test.exp.rightEmpty {
			t.Errorf("test '%s' failed for right empty columns, expected=%d. got=%d", test.name, test.exp.rightEmpty, result)
		}
	}
}
