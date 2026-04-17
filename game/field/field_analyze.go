// Copyright (c) 2025, 2026 by Marko Gaćeša
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

func (f *Field) FindBlizzardTops() []block.XY {
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

func (f *Field) FindAcidRainTops() []block.XY {
	tops := make([]block.XY, 0, f.w)

	for x := range f.w {
		y := f.GetTopmostFull(x)
		if y < 0 || !f.blocks[y*f.w+x].Type.Shootable() {
			continue
		}

		tops = append(tops, block.XY{X: x, Y: y})
	}

	return tops
}

type ColumnSection struct {
	Column  int
	RowFrom int
	RowTo   int
}

func (f *Field) FindMovableColumnSections(col int, filter func(*Field, ColumnSection) bool) []ColumnSection {
	var start int
	var sections []ColumnSection
	start = -1
	for row, idx := 0, col; row < f.h; row, idx = row+1, idx+f.w {
		b := f.blocks[idx].Block
		if b.Type.IsImmovable() {
			if start >= 0 {
				section := ColumnSection{Column: col, RowFrom: start, RowTo: row}
				if filter == nil || filter(f, section) {
					sections = append(sections, section)
				}
				start = -1
			}
		} else {
			if start < 0 {
				start = row
			}
		}
	}
	if start >= 0 {
		section := ColumnSection{Column: col, RowFrom: start, RowTo: f.h}
		if filter == nil || filter(f, section) {
			sections = append(sections, section)
		}
	}

	return sections
}

func (f *Field) FindMovableSections(filter func(*Field, ColumnSection) bool) []ColumnSection {
	sections := make([]ColumnSection, 0, f.w)
	for col := range f.w {
		sections = append(sections, f.FindMovableColumnSections(col, filter)...)
	}
	return sections
}
