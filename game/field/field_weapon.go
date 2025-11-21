// Copyright (c) 2025 by Marko Gaćeša

package field

import (
	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/logic/random"
)

func (f *Field) Blizzard(intensity int) []block.XY {
	if intensity <= 0 {
		return nil
	}

	tops := f.FindTops()
	result := make([]block.XY, 0, intensity)

	var k int
	for _, top := range tops {
		k += top.Y*f.w + top.X
	}

	r := random.New(uint32(f.seed), uint32(intensity+k))

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
