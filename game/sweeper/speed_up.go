// Copyright (c) 2025, 2026 by Marko Gaćeša
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

func (s *SpeedUp) Analyze(events event.Reader) {
	var removed int
	events.Range(func(e event.Event) {
		switch v := e.(type) {
		case *op.FieldBlockSet:
			if v.Op == op.TypeClear {
				removed++
			}
		case *op.FieldDestroyRow:
			removed += s.field.GetWidth()
		case *op.FieldDestroyColumn:
			removed++
		}
	})

	if removed == 0 {
		return
	}

	s.base.start()
}

func (s *SpeedUp) Sweep(p event.Pusher) {
	w := s.field.GetWidth()
	b := s.field.GetBlocksRemoved()
	n := byte(s.field.Ctrls())
	for i := range n {
		needed := field.LevelUpBlocks(int(s.field.CtrlLevel(i)+1), w)
		if b >= needed {
			p.Push(op.NewPieceSpeedUp(i, 1))
		}
	}
	s.endIteration()
}
