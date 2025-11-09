// Copyright (c) 2020 by Marko Gaćeša

package field

import (
	"fmt"
	"github.com/marko-gacesa/gamatet/game/block"
	"testing"
)

func TestField_ShiftColumnDownByN(t *testing.T) {
	const fieldHeight = 8
	const s = block.TypeRock
	const e = block.TypeEmpty
	tests := []struct {
		name         string
		y, n, height int
		blocks       []block.Block
		expect       []block.Block
		expectedAnim []int
	}{
		{
			name: "just one block, no block moves",
			y:    2, n: 0, height: 1,
			blocks:       []block.Block{{Type: s, Color: 1}, {Type: s, Color: 2}, {Type: s, Color: 3}, {Type: e}},
			expect:       []block.Block{{Type: s, Color: 1}, {Type: s, Color: 2}, {Type: e}, {Type: e}},
			expectedAnim: []int{0, 0, 0, 0},
		},
		{
			name: "move full row down",
			y:    2, n: 5, height: 3,
			blocks:       []block.Block{{Type: e, Color: 0}, {Type: e, Color: 0}, {Type: s, Color: 3}, {Type: s, Color: 4}, {Type: s, Color: 5}, {Type: s, Color: 6}, {Type: s, Color: 7}, {Type: s, Color: 8}},
			expect:       []block.Block{{Type: s, Color: 4}, {Type: s, Color: 5}, {Type: s, Color: 6}, {Type: s, Color: 7}, {Type: s, Color: 8}, {Type: e}, {Type: e}, {Type: e}},
			expectedAnim: []int{1, 1, 1, 1, 1, 0, 0, 0},
		},
		{
			name: "move sparsely filled row down by 2",
			y:    1, n: 6, height: 2,
			blocks:       []block.Block{{Type: e, Color: 0}, {Type: e, Color: 2}, {Type: s, Color: 3}, {Type: e, Color: 0}, {Type: s, Color: 5}, {Type: s, Color: 6}, {Type: e, Color: 0}, {Type: s, Color: 8}},
			expect:       []block.Block{{Type: s, Color: 3}, {Type: e, Color: 0}, {Type: s, Color: 5}, {Type: s, Color: 6}, {Type: e, Color: 0}, {Type: s, Color: 8}, {Type: e}, {Type: e}},
			expectedAnim: []int{1, 0, 1, 1, 0, 1, 0, 0},
		},
		{
			name: "high fall",
			y:    6, n: 1, height: 5,
			blocks:       []block.Block{{Type: s, Color: 1}, {Type: s, Color: 2}, {Type: e, Color: 0}, {Type: s, Color: 4}, {Type: s, Color: 5}, {Type: s, Color: 6}, {Type: s, Color: 7}, {Type: s, Color: 8}},
			expect:       []block.Block{{Type: s, Color: 1}, {Type: s, Color: 2}, {Type: s, Color: 8}, {Type: s, Color: 4}, {Type: s, Color: 5}, {Type: s, Color: 6}, {Type: e}, {Type: e}},
			expectedAnim: []int{0, 0, 1, 0, 0, 0, 0, 0},
		},
	}

	const column = 1

	for _, test := range tests {
		f := Make(4, fieldHeight, 0)
		f.Config.Anim = true

		for row := 0; row < len(test.blocks); row++ {
			f.setXY(column, row, test.blocks[row])
		}

		removedBlock := f.GetXY(column, test.y)
		f.ShiftColumnDownByN(column, test.y, test.n, test.height)

		def := fmt.Sprintf("row=%d n=%d h=%d", test.y, test.n, test.height)

		for row := 0; row < len(test.expect); row++ {
			b, a := f.getXY(column, row)
			animCount := a.Count()

			if b != test.expect[row] {
				t.Errorf("test %q (%s) failed: at row=%d expected block=%+v, but got block=%+v", test.name, def, row, test.expect[row], b)
			}

			if animCount != test.expectedAnim[row] {
				t.Errorf("test %q (%s) failed: at row=%d expected anim count=%d, but got count=%d", test.name, def, row, test.expectedAnim[row], animCount)
			}
		}

		f.UndoShiftColumnByN(column, test.y, test.n, test.height, removedBlock)

		for row := 0; row < len(test.blocks); row++ {
			b, _ := f.getXY(column, row)
			if b != test.blocks[row] {
				t.Errorf("test %q (%s) failed: at row=%d expected block=%+v, but got block=%+v", test.name, def, row, test.expect[row], b)
			}
		}
	}
}

func TestField_ShiftRowsDown(t *testing.T) {
	f := Make(6, 6, 0)
	b := block.Block{Type: block.TypeRock}

	colMin := 0
	colMax := 5

	// 5 . . . . . #
	// 4 # . . . . .
	// 3 . . . # . .
	// 2 . . # . . .
	// 1 . . . . # .
	// 0 . # . . . .
	//   0 1 2 3 4 5

	f.setXY(1, 0, b)
	f.setXY(4, 1, b)
	f.setXY(2, 2, b)
	f.setXY(3, 3, b)
	f.setXY(0, 4, b)
	f.setXY(5, 5, b)

	// 5 . . . . . .
	// 4 . . . . . #
	// 3 # . . . . .
	// 2 . . . # . .
	// 1 . . # . . .
	// 0 . . . . # .
	//   0 1 2 3 4 5

	f.ShiftRowsDown(0)

	if f._isXYEmpty(4, 0, colMin, colMax, false, -1) || !f._isXYEmpty(4, 1, colMin, colMax, false, -1) ||
		f._isXYEmpty(2, 1, colMin, colMax, false, -1) || !f._isXYEmpty(2, 2, colMin, colMax, false, -1) ||
		f._isXYEmpty(3, 2, colMin, colMax, false, -1) || !f._isXYEmpty(3, 3, colMin, colMax, false, -1) ||
		f._isXYEmpty(0, 3, colMin, colMax, false, -1) || !f._isXYEmpty(0, 4, colMin, colMax, false, -1) ||
		f._isXYEmpty(5, 4, colMin, colMax, false, -1) || !f._isXYEmpty(5, 5, colMin, colMax, false, -1) {
		t.Errorf("test failed for y=0")
	}

	// 5 . . . . . .
	// 4 . . . . . .
	// 3 . . . . . #
	// 2 # . . . . .
	// 1 . . # . . .
	// 0 . . . . # .
	//   0 1 2 3 4 5

	f.ShiftRowsDown(2)

	if f._isXYEmpty(4, 0, colMin, colMax, false, -1) ||
		f._isXYEmpty(2, 1, colMin, colMax, false, -1) ||
		!f._isXYEmpty(3, 2, colMin, colMax, false, -1) ||
		f._isXYEmpty(0, 2, colMin, colMax, false, -1) || !f._isXYEmpty(0, 3, colMin, colMax, false, -1) ||
		f._isXYEmpty(5, 3, colMin, colMax, false, -1) || !f._isXYEmpty(5, 4, colMin, colMax, false, -1) {
		t.Errorf("test failed for y=2")
	}
}
