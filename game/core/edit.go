// Copyright (c) 2020-2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package core

import (
	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/op"
	"github.com/marko-gacesa/gamatet/game/piece"
	"github.com/marko-gacesa/gamatet/game/setup"
)

func conjureBlock(p event.Pusher, x, y int, b block.Block) {
	p.Push(op.NewFieldBlockSet(x, y, op.TypeSet, field.AnimPop, 0, b))
}

func putBlock(p event.Pusher, x, y int, b block.Block) {
	p.Push(op.NewFieldBlockSet(x, y, op.TypeSet, field.AnimNo, 0, b))
}

func InitRandomBlocks(f *field.Field, p event.Pusher) {
	rnd := f.Random(0)
	w := f.GetWidth()
	n := w + w/2
	m := make(map[block.XY]struct{})
	c := piece.NewRandomColor(setup.ColorRGB[:], 0)
	for len(m) < n {
		xy := block.XY{
			X: rnd.Int(w),
			Y: rnd.Int(4),
		}
		if _, exists := m[xy]; exists {
			continue
		}

		conjureBlock(p, xy.X, xy.Y, block.Block{
			Type:     block.TypeRock,
			Hardness: 0,
			Color:    c.Color(uint(xy.Y*w+xy.X), 0),
		})
		m[xy] = struct{}{}
	}
}

func InitTriangle(f *field.Field, p event.Pusher) {
	w := f.GetWidth()
	h := f.GetHeight()
	c := piece.NewRandomColor(setup.ColorRGB[:], 0)
	d := min(w, h)
	for y := 0; y < d; y++ {
		for x := w - d + 1 + y; x < w; x++ {
			idx := y*w + x
			putBlock(p, x, y, block.Block{
				Type:     block.TypeRock,
				Hardness: 0,
				Color:    c.Color(uint(idx), byte(idx%8)),
			})
		}
	}
}

func InitFunnel(f *field.Field, p event.Pusher) {
	w := f.GetWidth()
	h := f.GetHeight()
	d := min((w-1)/2, h)
	for y := 0; y < d; y++ {
		for x := 0; x < d-y; x++ {
			putBlock(p, x, y, block.Wall)
			putBlock(p, w-x-1, y, block.Wall)
		}
	}
}
