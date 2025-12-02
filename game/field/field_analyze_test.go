// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package field

import (
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/marko-gacesa/gamatet/game/block"
)

func TestField_RangeBlocks(t *testing.T) {
	const (
		fieldW = 6
		fieldH = 6
	)

	tests := []struct {
		name   string
		blocks []block.XYB
	}{
		{
			name:   "empty",
			blocks: []block.XYB{},
		},
		{
			name: "corners",
			blocks: []block.XYB{
				{XY: block.XY{0, 0}, Block: block.Rock},
				{XY: block.XY{fieldW - 1, 0}, Block: block.Ruby},
				{XY: block.XY{0, fieldH - 1}, Block: block.Iron},
				{XY: block.XY{fieldW - 1, fieldH - 1}, Block: block.Goal},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := Make(6, 6, 0)
			for _, xyb := range test.blocks {
				f.setXY(xyb.X, xyb.Y, xyb.Block)
			}

			idx := 0
			f.RangeBlocks(func(xyb block.XYB) bool {
				if test.blocks[idx] != xyb {
					t.Errorf("block index mismatch: got %v, want %v", xyb, test.blocks[idx])
				}
				idx++
				return true
			})
			if idx != len(test.blocks) {
				t.Errorf("len mismatch: got %v, want %v", idx, len(test.blocks))
			}

			var called bool
			f.RangeBlocks(func(xyb block.XYB) bool {
				if called {
					t.Error("already called")
				}
				called = true
				return false
			})
		})
	}
}

func TestField_FindBlizzardTops(t *testing.T) {
	b := block.Block{Type: block.TypeRock}

	fillToH := func(f *Field, col, h int) {
		for i := range h + 1 {
			f.setXY(col, i, b)
		}
	}

	tests := []struct {
		name    string
		heights []int
		want    []block.XY
	}{
		{
			name:    "empty",
			heights: []int{},
			want:    []block.XY{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}, {X: 4, Y: 0}, {X: 5, Y: 0}},
		},
		{
			// 5 . . # . . .
			// 4 # . # . . .
			// 3 # . # . # .
			// 2 # # # . # #
			// 1 # # # # # #
			// 0 # # # # # #
			//   0 1 2 3 4 5
			name:    "right_side_blocked",
			heights: []int{4, 2, 5, 1, 3, 2},
			want:    []block.XY{{X: 0, Y: 5}, {X: 4, Y: 4}},
		},
		{
			// 5 . . . . . .
			// 4 . . . . . .
			// 3 . . . . . .
			// 2 . # . # # #
			// 1 # # # # # #
			// 0 # # # # # #
			//   0 1 2 3 4 5
			name:    "left_side_blocked",
			heights: []int{1, 2, 1, 2, 2, 2},
			want:    []block.XY{{X: 1, Y: 3}, {X: 3, Y: 3}, {X: 4, Y: 3}, {X: 5, Y: 3}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := Make(6, 6, 0)
			for col, h := range test.heights {
				fillToH(f, col, h)
			}

			want := test.want
			got := f.FindBlizzardTops()

			if !cmp.Equal(got, want) {
				t.Errorf("failed:\n%s\n", cmp.Diff(got, want))
			}
		})
	}
}

func TestField_FindMovableColumnSections(t *testing.T) {
	bE := block.Block{Type: block.TypeEmpty}
	bR := block.Block{Type: block.TypeRock}
	bI := block.Block{Type: block.TypeWall}

	const dimW = 6
	const dimH = 6

	setColumn := func(f *Field, col int, blocks []block.Block) {
		for row := 0; row < dimH; row++ {
			f.setXY(col, row, blocks[row])
		}
	}

	tests := []struct {
		name     string
		column   []block.Block
		filter   func(f *Field, section ColumnSection) bool
		expected []ColumnSection
	}{
		{
			name:   "empty",
			column: []block.Block{bE, bE, bE, bE, bE, bE},
			expected: []ColumnSection{
				{RowFrom: 0, RowTo: dimH},
			},
		},
		{
			name:   "with_movable",
			column: []block.Block{bE, bR, bR, bE, bE, bR},
			expected: []ColumnSection{
				{RowFrom: 0, RowTo: dimH},
			},
		},
		{
			name:   "one",
			column: []block.Block{bE, bR, bI, bE, bE, bR},
			expected: []ColumnSection{
				{RowFrom: 0, RowTo: 2},
				{RowFrom: 3, RowTo: dimH},
			},
		},
		{
			name:   "two",
			column: []block.Block{bE, bR, bI, bE, bI, bR},
			expected: []ColumnSection{
				{RowFrom: 0, RowTo: 2},
				{RowFrom: 3, RowTo: 4},
				{RowFrom: 5, RowTo: dimH},
			},
		},
		{
			name:   "two_consecutive",
			column: []block.Block{bE, bI, bI, bE, bE, bR},
			expected: []ColumnSection{
				{RowFrom: 0, RowTo: 1},
				{RowFrom: 3, RowTo: dimH},
			},
		},
		{
			name:   "border",
			column: []block.Block{bI, bI, bR, bE, bE, bI},
			expected: []ColumnSection{
				{RowFrom: 2, RowTo: 5},
			},
		},
		{
			name:   "filter_top_empty",
			column: []block.Block{bI, bR, bE, bI, bE, bR},
			filter: func(f *Field, section ColumnSection) bool {
				b, _ := f.getXY(section.Column, section.RowTo-1)
				return b.Type == block.TypeEmpty
			},
			expected: []ColumnSection{
				{RowFrom: 1, RowTo: 3},
			},
		},
		{
			name:   "filter_bottom_empty",
			column: []block.Block{bE, bR, bI, bE, bE, bR},
			filter: func(f *Field, section ColumnSection) bool {
				b, _ := f.getXY(section.Column, section.RowTo-1)
				return b.Type == block.TypeEmpty
			},
			expected: []ColumnSection{
				{RowFrom: 0, RowTo: 2},
				{RowFrom: 3, RowTo: dimH},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := Make(dimW, dimH, 0)

			const col = 0

			setColumn(f, col, test.column)
			sections := f.FindMovableColumnSections(col, test.filter)

			if want, got := test.expected, sections; !slices.Equal(want, got) {
				t.Errorf("want: %v, got: %v", want, got)
			}
		})
	}
}
