// Copyright (c) 2024, 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package sweeper

import (
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/op"
)

var _ Sweeper = (*Score)(nil)

func NewScore(f *field.Field) *Score {
	b := newBase(f)
	return &Score{base: *b}
}

type Score struct {
	base
	delta int
}

func (s *Score) Start(analyzer *Analyzer) bool {
	n := analyzer.stats.removed + analyzer.stats.softened
	if n == 0 {
		return false
	}

	w := s.field.GetWidth()
	multiplier := (n + w>>1) / w

	if multiplier == 0 {
		s.delta = 5 * n
	} else {
		s.delta = 10 * multiplier * n
	}

	return s.base.Start(analyzer)
}

func (s *Score) Sweep(p event.Pusher) {
	for i := range s.field.Ctrls() {
		ctrl := s.field.Ctrl(byte(i))
		level := ctrl.Level
		scoreDelta := s.delta * int(level)
		p.Push(op.NewPieceScore(i, scoreDelta))
	}
	s.delta = 0
	s.endIteration()
}
