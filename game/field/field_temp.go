// Copyright (c) 2024 by Marko Gaćeša

package field

import (
	"gamatet/game/block"
	"gamatet/logic/anim"
)

func NewTemp(f *Field) *temp {
	return &temp{
		f: f,
	}
}

type temp struct {
	f       *Field
	entries []tempEntry
}

type tempEntry struct {
	x, y int
	elem elem
}

func (t *temp) Set(x, y int, b block.Block) {
	f := t.f
	idx := y*f.w + x

	e := f.blocks[idx]
	t.entries = append(t.entries, tempEntry{
		x:    x,
		y:    y,
		elem: e,
	})

	f.blocks[idx] = elem{
		Block: b,
		List:  anim.List{},
	}
}

func (t *temp) Revert() {
	f := t.f
	for i := len(t.entries) - 1; i >= 0; i-- {
		x := t.entries[i].x
		y := t.entries[i].y
		idx := y*f.w + x
		f.blocks[idx] = t.entries[i].elem
	}
}
