// Copyright (c) 2025 by Marko Gaćeša

package field

import "github.com/marko-gacesa/gamatet/game/block"

func (f *Field) FindTops() []block.XY {
	tops := make([]block.XY, 0, f.w)

	for x := range f.w {
		y := f.GetTopmostEmpty(x)
		if y >= f.h {
			continue
		}

		idx := y*f.w + x

		var ok bool

		switch x {
		case 0:
			ok = f.blocks[idx+1].Type == block.TypeEmpty
		case f.w - 1:
			ok = f.blocks[idx-1].Type == block.TypeEmpty
		default:
			ok = f.blocks[idx-1].Type == block.TypeEmpty || f.blocks[idx+1].Type == block.TypeEmpty
		}

		if ok {
			tops = append(tops, block.XY{X: x, Y: y})
		}
	}

	return tops
}
