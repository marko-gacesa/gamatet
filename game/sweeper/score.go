// Copyright (c) 2024-2026 by Marko Gaćeša
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
	g := len(analyzer.blocks.goalsRemoved) + len(analyzer.blocks.gnawsKilled)
	if n == 0 && g == 0 {
		return false
	}

	if n > 0 {
		// Score bonus per number of removed blocks (example for width 10):
		// n=1..4   : bonus=n*5  ; so clearing one block (n=1) gives 5
		// n=5..14  : bonus=n*10 ; so clearing one row (n=10) gives 100
		// n=15..24 : bonus=n*20 ; so clearing two rows (n=20) gives 400
		// n=25..34 : bonus=n*30 ; so clearing three rows (n=30) gives 900
		// n=35..44 : bonus=n*40 ; so clearing four rows (n=40) gives 1600
		// n=45..54 : bonus=n*50 ; so clearing five rows (n=50) gives 2500
		// Everything is times the current game speed level.
		// ... so clearing 3 rows at level=7 gives 6300 points.

		w := s.field.GetWidth()
		multiplier := (n + w>>1) / w

		if multiplier == 0 {
			s.delta = 5 * n
		} else {
			s.delta = 10 * multiplier * n
		}
	}

	if g > 0 {
		// Score bonus for each goal or gnaw is 1000 times the current game speed level.
		// ... so removing a goal block at level=7 gives 7000 points.

		s.delta += g * 1000
	}

	return s.base.Start(analyzer)
}

func (s *Score) Sweep(p event.Pusher) {
	for ctrlIdx := range s.field.Ctrls() {
		ctrl := s.field.Ctrl(byte(ctrlIdx))
		level := ctrl.Level
		scoreDelta := s.delta * int(level)
		p.Push(op.NewPieceScore(ctrlIdx, scoreDelta))
	}
	s.delta = 0
	s.endIteration()
}
