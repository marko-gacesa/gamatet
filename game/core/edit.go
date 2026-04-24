// Copyright (c) 2020-2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package core

import (
	"math/rand"

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

func FieldRandomBlocks(f *field.Field, p event.Pusher) {
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

func FieldInit1(f *field.Field, p event.Pusher) {
	w := f.GetWidth()
	h := f.GetHeight()
	c := piece.NewRandomColor(setup.ColorRGB[:], 0)
	for yi := range h / 2 {
		for xi := range w {
			if xi == yi%w {
				continue
			}
			putBlock(p, xi, yi, block.Block{
				Type:     block.TypeRock,
				Hardness: byte(rand.Int31n(4)),
				Color:    c.Color(uint(yi*w+xi), 0),
			})
		}
	}
}
