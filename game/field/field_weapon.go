// Copyright (c) 2025, 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package field

import (
	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/logic/random"
)

func (f *Field) Blizzard(intensity int) []block.XY {
	if intensity <= 0 {
		return nil
	}

	tops := f.FindBlizzardTops()
	result := make([]block.XY, 0, intensity)

	var k int
	for _, top := range tops {
		k += top.Y*f.w + top.X
	}

	r := random.New(uint64(intensity+k), uint64(f.seed))

	for range intensity {
		if len(tops) == 0 {
			break
		}

		idx := r.Int(len(tops))
		result = append(result, tops[idx])

		tops[idx].Y++
		if tops[idx].Y >= f.h {
			tops = append(tops[:idx], tops[idx+1:]...)
		}
	}

	return result
}

func (f *Field) GetRandomBlock() (block.XYB, bool) {
	var count int
	f.RangeBlocks(func(xyb block.XYB) bool {
		if xyb.Block.Type != block.TypeRock || xyb.Block.Hardness != 0 {
			return true
		}
		count++
		return true
	})
	if count == 0 {
		return block.XYB{}, false
	}

	r := random.New(uint64(count), uint64(f.seed))
	idx := r.Int(count)

	count = 0
	var chosen block.XYB
	var ok bool
	f.RangeBlocks(func(xyb block.XYB) bool {
		if xyb.Block.Type != block.TypeRock || xyb.Block.Hardness != 0 {
			return true
		}
		if count == idx {
			chosen = xyb
			ok = true
			return false
		}
		count++
		return true
	})

	return chosen, ok
}
