// Copyright (c) 2025 by Marko Gaćeša

package field

import (
	"github.com/marko-gacesa/gamatet/game/block"
)

func (f *Field) Blizzard(intensity int) []block.XY {
	if intensity <= 0 {
		return nil
	}

	tops := f.FindTops()
	result := make([]block.XY, 0, intensity)

	for range intensity {
		if len(tops) == 0 {
			break
		}

		idx := f.rand.IntN(len(tops))
		result = append(result, tops[idx])

		tops[idx].Y++
		if tops[idx].Y >= f.h {
			tops = append(tops[:idx], tops[idx+1:]...)
		}
	}

	return result
}
