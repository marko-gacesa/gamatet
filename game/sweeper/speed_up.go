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
		needed := levelUpBlocks(s.field.Ctrl(i).Level+1, w)
		if b >= needed {
			p.Push(op.NewPieceSpeedUp(i, 1))
		}
	}
	s.endIteration()
}

// levelUpBlocks returns number of destroyed blocks needed to reach the desired level of speed.
//
// Based on the code used in GMT1 30 years earlier:
// To reach the next level from current level it takes this many lines
// progression := []int{11, 13, 15, 17, 19, 21, 23, 25, 27, 29}
// neededLines := progression[current_level+1]
func levelUpBlocks(l, w int) int {
	needed := (10*l + l*l) * w
	return needed
}
