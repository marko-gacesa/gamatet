// Copyright (c) 2024-2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package sweeper

import (
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/op"
)

var _ Sweeper = (*Shaker)(nil)

func NewShaker(f *field.Field) *Shaker {
	b := newBase(f)
	return &Shaker{base: *b}
}

type Shaker struct {
	base
	intensity byte
}

func (s *Shaker) Analyze(events event.Reader) {
	var blocks int
	events.Range(func(e event.Event) {
		switch v := e.(type) {
		case *op.FieldBlockSet:
			if v.Op == op.TypeClear {
				blocks++
			}
		case *op.FieldDestroyRow:
			blocks += s.field.GetWidth()
		case *op.FieldDestroyColumn:
			blocks++
		}
	})

	w := s.field.GetWidth()

	intensity := (blocks + w>>1) / w
	if intensity < 2 {
		return
	}
	if intensity > 5 {
		intensity = 5
	}

	s.intensity = byte(intensity)
	s.base.start()
}

func (s *Shaker) Sweep(p event.Pusher) {
	p.Push(op.NewFieldQuake(s.intensity))
	s.intensity = 0
	s.endIteration()
}
