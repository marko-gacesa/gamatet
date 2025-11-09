// Copyright (c) 2024, 2025 by Marko Gaćeša

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

func (s *Shaker) Start(analyzer *Analyzer) bool {
	w := s.field.GetWidth()

	intensity := (analyzer.blocks.removed + w>>1) / w
	if intensity < 2 {
		return false
	}
	if intensity > 5 {
		intensity = 5
	}

	s.intensity = byte(intensity)
	return s.base.Start(analyzer)
}

func (s *Shaker) Sweep(p event.Pusher) {
	p.Push(op.NewFieldQuake(s.intensity))
	s.intensity = 0
	s.endIteration()
}
