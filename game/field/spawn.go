// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package field

import (
	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/logic/random"
)

type SpawnLocation byte

const (
	SpawnLocationLowerMid SpawnLocation = iota
	SpawnLocationBottomMid
)

func (f *Field) SpawnLocation(loc SpawnLocation) (block.XY, bool) {
	w := f.GetWidth()
	h := f.GetHeight()

	var (
		start   block.XY
		reach   int
		entropy int
	)

	switch loc {
	case SpawnLocationLowerMid:
		start = block.XY{X: w / 2, Y: h / 3}
		reach = 2 * w / 3
	case SpawnLocationBottomMid:
		start = block.XY{X: w / 2, Y: 0}
		reach = w/2 + 1
	}

	potentialXY := make([]block.XY, 0, 10)

	f.FindNearest8(start, reach, func(xyb block.XYB, i int) bool {
		idx := xyb.Y*w + xyb.X
		entropy += idx

		if xyb.Block.Type != block.TypeEmpty {
			entropy += idx * int(xyb.Block.Type)
			return false
		}

		potentialXY = append(potentialXY, xyb.XY)
		return len(potentialXY) == cap(potentialXY)
	})

	if len(potentialXY) == 0 {
		return block.XY{}, false
	}

	r := random.New(uint64(entropy), uint64(f.seed))
	n := r.Int(len(potentialXY))
	xy := potentialXY[n]

	return xy, true
}
