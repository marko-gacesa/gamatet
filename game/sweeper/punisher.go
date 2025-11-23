// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package sweeper

import (
	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/op"
)

var _ Sweeper = (*Punisher)(nil)

func NewPunisher(f *field.Field, others []FieldPusher) *Punisher {
	b := newBase(f)
	return &Punisher{
		base:   *b,
		others: others,
	}
}

type Punisher struct {
	base
	others    []FieldPusher
	intensity int
}

func (s *Punisher) Start(analyzer *Analyzer) bool {
	w := s.field.GetWidth()

	intensity := blockCount((analyzer.stats.removed + w>>1) / w)
	if intensity < 1 {
		return false
	}

	s.intensity = intensity

	return s.base.Start(analyzer)
}

func (s *Punisher) Sweep(event.Pusher) {
	for i := range s.others {
		f := s.others[i].Field
		if f.IsFinished() {
			continue
		}

		h := f.GetHeight()
		xys := f.Blizzard(s.intensity)

		for _, xy := range xys {
			s.others[i].Pusher.Push(&op.FieldBlockSet{
				Col:       byte(xy.X),
				Row:       byte(xy.Y),
				Op:        op.TypeSet,
				AnimType:  field.AnimFall,
				AnimParam: byte(h - xy.Y),
				Block:     block.Rock,
			})
		}
	}

	s.intensity = 0
	s.endIteration()
}

func blockCount(intensity int) int {
	if intensity <= 1 {
		return intensity
	}
	return fib(intensity + 1)
}

func fib(n int) int {
	prev, curr := 0, 1
	for range n {
		prev, curr = curr, prev+curr
	}
	return curr
}
