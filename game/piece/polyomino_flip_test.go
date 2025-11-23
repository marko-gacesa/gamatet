// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package piece

import (
	"testing"

	"github.com/marko-gacesa/gamatet/game/block"
)

func TestPolyomino_FlipV(t *testing.T) {
	tests := []struct {
		name string
		w, h byte
		data []bool
		exp  []bool
	}{
		{
			name: "1x1",
			w:    1,
			h:    1,
			data: []bool{XX},
			exp:  []bool{XX},
		},
		{
			name: "2x2",
			w:    2,
			h:    2,
			data: []bool{
				__, XX,
				XX, __,
			},
			exp: []bool{
				XX, __,
				__, XX,
			},
		},
		{
			name: "3x3",
			w:    3,
			h:    3,
			data: []bool{
				__, XX, __,
				XX, XX, XX,
				XX, __, __,
			},
			exp: []bool{
				XX, __, __,
				XX, XX, XX,
				__, XX, __,
			},
		},
		{
			name: "4x4",
			w:    4,
			h:    4,
			data: []bool{
				__, XX, __, XX,
				XX, XX, XX, XX,
				XX, __, __, __,
				XX, XX, XX, XX,
			},
			exp: []bool{
				XX, XX, XX, XX,
				XX, __, __, __,
				XX, XX, XX, XX,
				__, XX, __, XX,
			},
		},
		{
			name: "5x4",
			w:    5,
			h:    4,
			data: []bool{
				__, __, XX, XX, XX,
				XX, XX, XX, __, __,
				__, XX, XX, __, __,
				__, XX, __, __, __,
			},
			exp: []bool{
				__, XX, __, __, __,
				__, XX, XX, __, __,
				XX, XX, XX, __, __,
				__, __, XX, XX, XX,
			},
		},
	}

	var b block.Block

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shape := _initShapeRect(test.w, test.h, test.data)
			piece := &polyominoFlipV{
				shapeRect: shape,
				block:     b,
			}

			piece.Activate()

			expShape := _initShapeRect(test.w, test.h, test.exp)

			if piece.data != expShape.data {
				t.Errorf("test %s FlipV flopped. expected blocks=%b, but got %b", test.name, expShape.data, piece.data)
				return
			}
		})
	}
}

func TestPolyomino_FlipH(t *testing.T) {
	tests := []struct {
		name string
		w, h byte
		data []bool
		exp  []bool
	}{
		{
			name: "1x1",
			w:    1,
			h:    1,
			data: []bool{XX},
			exp:  []bool{XX},
		},
		{
			name: "2x2",
			w:    2,
			h:    2,
			data: []bool{
				__, XX,
				XX, __,
			},
			exp: []bool{
				XX, __,
				__, XX,
			},
		},
		{
			name: "3x3",
			w:    3,
			h:    3,
			data: []bool{
				__, XX, __,
				XX, XX, XX,
				XX, __, __,
			},
			exp: []bool{
				__, XX, __,
				XX, XX, XX,
				__, __, XX,
			},
		},
		{
			name: "4x4",
			w:    4,
			h:    4,
			data: []bool{
				__, XX, __, XX,
				XX, XX, XX, XX,
				XX, __, __, __,
				XX, XX, XX, XX,
			},
			exp: []bool{
				XX, __, XX, __,
				XX, XX, XX, XX,
				__, __, __, XX,
				XX, XX, XX, XX,
			},
		},
		{
			name: "5x4",
			w:    5,
			h:    4,
			data: []bool{
				__, __, XX, XX, XX,
				XX, XX, XX, __, __,
				__, XX, XX, __, __,
				__, XX, __, __, __,
			},
			exp: []bool{
				XX, XX, XX, __, __,
				__, __, XX, XX, XX,
				__, __, XX, XX, __,
				__, __, __, XX, __,
			},
		},
	}

	var b block.Block

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shape := _initShapeRect(test.w, test.h, test.data)
			piece := &polyominoFlipH{
				shapeRect: shape,
				block:     b,
			}

			piece.Activate()

			expShape := _initShapeRect(test.w, test.h, test.exp)

			if piece.data != expShape.data {
				t.Errorf("test %s FlipH flopped. expected blocks=%b, but got %b", test.name, expShape.data, piece.data)
				return
			}
		})
	}
}
