// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package field

import "github.com/marko-gacesa/gamatet/game/block"

// RangeBlocks goes over every block on the field and runs inspect function for each.
// The operation is terminated when inspect function returns false.
func (f *Field) RangeBlocks(inspect func(xyb block.XYB) bool) {
	idx := 0
	w := f.w
	h := f.h
	for y := range h {
		for x := range w {
			b := f.blocks[idx].Block
			idx++

			if b.Type == block.TypeEmpty {
				continue
			}
			if !inspect(block.XYB{
				XY:    block.XY{X: x, Y: y},
				Block: b,
			}) {
				return
			}
		}
	}
}

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
