// Copyright (c) 2025 by Marko Gaćeša

package field

import (
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

func TestField_FindTops(t *testing.T) {
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
			got := f.FindTops()

			if !cmp.Equal(got, want) {
				t.Errorf("failed:\n%s\n", cmp.Diff(got, want))
			}
		})
	}
}
