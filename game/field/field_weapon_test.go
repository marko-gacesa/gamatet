// Copyright (c) 2025 by Marko Gaćeša

package field

import (
	"github.com/google/go-cmp/cmp"
	"github.com/marko-gacesa/gamatet/game/block"
	"sort"
	"testing"
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
