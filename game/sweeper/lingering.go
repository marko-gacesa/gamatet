// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package sweeper

import (
	"time"

	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/op"
)

var _ Sweeper = (*Lingering)(nil)

func NewLingering(f *field.Field) *Lingering {
	b := newBase(f)
	return &Lingering{
		base: *b,
	}
}

type Lingering struct {
	base
	counter uint64
}

func (s *Lingering) Analyze(events event.Reader) {
	var has bool
	events.Range(func(e event.Event) {
		if v, ok := e.(*op.FieldLingering); ok && v.Delta > 0 {
			has = true
		}
	})

	if has {
		s.base.start()
	}
}

func (s *Lingering) Sweep(p event.Pusher) {
	effect, amount := s.field.LingeringEffect()
	if effect == field.EffectNone || amount == 0 {
		s.endIteration()
		return
	}

	p.Push(op.NewFieldLingering(effect, -1))

	switch effect {
	case field.EffectAcidRain:
		tops := s.field.FindAcidRainTops()
		if len(tops) == 0 {
			break
		}

		r := s.field.Random(s.counter).Int(len(tops))
		s.counter++

		top := tops[r]
		x := top.X
		y := top.Y

		b := s.field.GetXY(x, y)
		fh := s.field.GetHeight()
		height := fh - top.Y

		p.Push(op.NewFieldExBlock(x, y, field.AnimFall, height, block.Acid))
		if b.Hardness > 0 {
			p.Push(op.NewFieldBlockHardness(x, y, -1, field.AnimSpin, height))
		} else {
			p.Push(op.NewFieldBlockSet(x, y, op.TypeClear, field.AnimPop, 0, b))
		}
	}

	s.reschedule(time.Second)
}
