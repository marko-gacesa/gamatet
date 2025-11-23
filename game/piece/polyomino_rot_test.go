// Copyright (c) 2020, 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package piece

import (
	"testing"

	"github.com/marko-gacesa/gamatet/game/block"
)

func TestPolyomino_Rotate(t *testing.T) {
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
				__, XX,
				XX, __,
			},
			expRot: 1,
			exp: []bool{
				XX, __,
				__, XX,
			},
		},
		{
			name: "3x3",
			rots: 4,
			rot:  3,
			data: []bool{
				__, XX, __,
				XX, XX, XX,
				XX, __, __,
			},
			expRot: 0,
			exp: []bool{
				XX, XX, __,
				__, XX, XX,
				__, XX, __,
			},
		},
		{
			name: "4x4",
			rots: 4,
			rot:  2,
			data: []bool{
				__, XX, __, XX,
				XX, XX, XX, XX,
				XX, __, __, __,
				XX, XX, XX, XX,
			},
			expRot: 3,
			exp: []bool{
				XX, XX, XX, __,
				XX, __, XX, XX,
				XX, __, XX, __,
				XX, __, XX, XX,
			},
		},
		{
			name: "5x5",
			rots: 4,
			rot:  3,
			data: []bool{
				__, __, __, __, __,
				__, __, XX, XX, XX,
				XX, XX, XX, __, __,
				__, XX, XX, __, __,
				__, XX, __, __, __,
			},
			expRot: 0,
			exp: []bool{
				__, __, XX, __, __,
				XX, XX, XX, __, __,
				__, XX, XX, XX, __,
				__, __, __, XX, __,
				__, __, __, XX, __,
			},
		},
	}

	var b block.Block

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shape := _initShapeSquare(test.rots, test.data)
			piece := &polyominoRot{
				shapeSquare: shape,
				rot:         test.rot,
				block:       b,
			}

			expShape := _initShapeSquare(test.rots, test.exp)

			piece.UndoActivate()

			if piece.rot != test.expRot {
				t.Errorf("test %s CW rotate failed. expected rot=%d, but got %d", test.name, test.expRot, piece.rot)
				return
			}

			if piece.data != expShape.data {
				t.Errorf("test %s CW rotate failed. expected blocks=%b, but got %b", test.name, expShape.data, piece.data)
				return
			}

			piece.Activate()

			if piece.rot != test.rot {
				t.Errorf("test %s CCW rotate failed. expected rot=%d, but got %d", test.name, test.rot, piece.rot)
				return
			}

			if piece.data != shape.data {
				t.Errorf("test %s CCW rotate failed. expected blocks=%b, but got %b", test.name, shape.data, piece.data)
				return
			}

			if test.rots <= 0 {
				return
			}

			for i := byte(0); i < test.rots; i++ {
				piece.UndoActivate()
			}

			if piece.rot != test.rot {
				t.Errorf("test %s full CW rotate failed. expected rot=%d, but got %d", test.name, test.rot, piece.rot)
				return
			}

			if piece.data != shape.data {
				t.Errorf("test %s full CW rotate failed. expected blocks=%b, but got %b", test.name, shape.data, piece.data)
				return
			}

			for i := byte(0); i < test.rots; i++ {
				piece.Activate()
			}

			if piece.rot != test.rot {
				t.Errorf("test %s full CCW rotate failed. expected rot=%d, but got %d", test.name, test.rot, piece.rot)
				return
			}

			if piece.data != shape.data {
				t.Errorf("test %s full CCW rotate failed. expected blocks=%b, but got %b", test.name, shape.data, piece.data)
				return
			}
		})
	}
}

func TestPolyomino_IsRowEmpty(t *testing.T) {
	type expected struct {
		emptyRows   []bool
		topEmpty    byte
		bottomEmpty byte
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
				__, __,
			},
			exp: expected{[]bool{false, true}, 0, 1},
		},
		{
			name: "3x3",
			data: []bool{
				__, __, __,
				XX, XX, XX,
				__, __, __,
			},
			exp: expected{[]bool{true, false, true}, 1, 1},
		},
		{
			name: "4x4",
			data: []bool{
				__, __, __, __,
				XX, XX, XX, XX,
				__, __, __, __,
				__, __, __, __,
			},
			exp: expected{[]bool{true, false, true, true}, 1, 2},
		},
		{
			name: "4x4 V",
			data: []bool{
				__, __, XX, __,
				__, __, XX, __,
				__, __, XX, __,
				__, __, XX, __,
			},
			exp: expected{[]bool{false, false, false, false}, 0, 0},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shape := _initShapeSquare(4, test.data)
			piece := &polyominoRot{shapeSquare: shape, rot: 0, block: block.Rock}

			for i := byte(0); i < piece.DimY(); i++ {
				result := piece.data.isSquareRowEmpty(piece.dim, i)
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
		})
	}
}

func TestPolyomino_IsColumnEmpty(t *testing.T) {
	type expected struct {
		emptyColumns []bool
		leftEmpty    byte
		rightEmpty   byte
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
				__, XX,
				__, XX,
			},
			exp: expected{[]bool{true, false}, 1, 0},
		},
		{
			name: "3x3",
			data: []bool{
				__, XX, __,
				__, XX, __,
				__, XX, __,
			},
			exp: expected{[]bool{true, false, true}, 1, 1},
		},
		{
			name: "4x4",
			data: []bool{
				__, XX, __, __,
				__, XX, __, __,
				__, XX, __, __,
				__, XX, __, __,
			},
			exp: expected{[]bool{true, false, true, true}, 1, 2},
		},
		{
			name: "4x4 H",
			data: []bool{
				__, __, __, __,
				__, __, __, __,
				XX, XX, XX, XX,
				__, __, __, __,
			},
			exp: expected{[]bool{false, false, false, false}, 0, 0},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shape := _initShapeSquare(4, test.data)
			piece := &polyominoRot{shapeSquare: shape, rot: 0, block: block.Rock}

			for i := byte(0); i < piece.DimX(); i++ {
				result := piece.data.isSquareColumnEmpty(piece.dim, i)
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
		})
	}
}
