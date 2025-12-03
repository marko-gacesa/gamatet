// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package sweeper

import (
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/op"
)

var _ Sweeper = (*SpeedUp)(nil)

func NewSpeedUp(f *field.Field) *SpeedUp {
	b := newBase(f)
	return &SpeedUp{
		base: *b,
	}
}

type SpeedUp struct {
	base
}

func (s *SpeedUp) Start(analyzer *Analyzer) bool {
	if analyzer.stats.removed <= 0 {
		return false
	}

	return s.base.Start(analyzer)
}

func (s *SpeedUp) Sweep(p event.Pusher) {
	w := s.field.GetWidth()
	b := s.field.GetBlocksRemoved()
	n := byte(s.field.Ctrls())
	for i := range n {
		needed := field.LevelUpBlocks(int(s.field.Ctrl(i).Level+1), w)
		if b >= needed {
			p.Push(op.NewPieceSpeedUp(i, 1))
		}
	}
	s.endIteration()
}
