// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package field

import (
	"slices"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/marko-gacesa/gamatet/game/block"
)

func TestField_Blizzard(t *testing.T) {
	f := Make(6, 6, 0)
	b := block.Rock

	heights := []int{4, 2, 5, 1, 3, 2}
	for col, h := range heights {
		for i := 0; i <= h; i++ {
			f.setXY(col, i, b)
		}
	}

	// 5 . . # . . .
	// 4 # . # . . .
	// 3 # . # . # .
	// 2 # # # . # #
	// 1 # # # # # #
	// 0 # # # # # #
	//   0 1 2 3 4 5

	// Asking for 4 blocks, but on the field there is only room for 3.
	got := f.Blizzard(4)

	sort.Slice(got, func(i, j int) bool {
		if got[i].Y == got[j].Y {
			return got[i].X < got[j].X
		}
		return got[i].Y < got[j].Y
	})

	want := []block.XY{{X: 4, Y: 4}, {X: 0, Y: 5}, {X: 4, Y: 5}}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestField_GetRandomBlock(t *testing.T) {
	const dimW = 10
	const dimH = 10

	tests := []struct {
		name   string
		blocks []block.XYB
		expect bool
	}{
		{
			name:   "empty",
			blocks: []block.XYB{},
			expect: false,
		},
		{
			name:   "one",
			blocks: []block.XYB{{block.XY{4, 5}, block.Rock}},
			expect: true,
		},
		{
			name: "two",
			blocks: []block.XYB{
				{block.XY{dimW - 1, dimH - 1}, block.Rock},
				{block.XY{0, 0}, block.Rock},
			},
			expect: true,
		},
		{
			name: "several",
			blocks: []block.XYB{
				{block.XY{0, 0}, block.Rock},
				{block.XY{1, 1}, block.Rock},
				{block.XY{2, 2}, block.Rock},
				{block.XY{3, 3}, block.Rock},
				{block.XY{4, 4}, block.Rock},
				{block.XY{5, 5}, block.Rock},
				{block.XY{6, 6}, block.Rock},
				{block.XY{7, 7}, block.Rock},
				{block.XY{8, 8}, block.Rock},
				{block.XY{9, 9}, block.Rock},
			},
			expect: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := Make(dimW, dimH, 0)
			for _, xyb := range test.blocks {
				f.setXY(xyb.X, xyb.Y, xyb.Block)
			}

			result, ok := f.GetRandomBlock()

			if ok != test.expect {
				t.Errorf("got %v, want %v", ok, test.expect)
				return
			}

			if ok && !slices.Contains(test.blocks, result) {
				t.Errorf("could not find %v in %v", result, test.blocks)
			}
		})
	}
}
