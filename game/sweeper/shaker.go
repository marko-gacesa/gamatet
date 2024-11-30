// Copyright (c) 2024 by Marko Gaćeša

package sweeper

import (
	"gamatet/game/event"
	"gamatet/game/field"
	"gamatet/game/op"
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

func (s *Shaker) Start(analyzer *Analyzer) {
	w := s.field.GetWidth()

	intensity := (analyzer.removed-w>>1)/w + 1
	if intensity < 2 {
		return
	}
	if intensity > 5 {
		intensity = 5
	}

	s.intensity = byte(intensity)
	s.base.Start(analyzer)
}

func (s *Shaker) Sweep(p event.Pusher) {
	p.Push(op.NewFieldQuake(s.intensity))
	s.intensity = 0
	s.endIteration()
}
