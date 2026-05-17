// Copyright (c) 2025, 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package sweeper

import (
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/op"
	"github.com/marko-gacesa/gamatet/game/piece"
)

var _ Sweeper = (*SpeedUpOnDefeat)(nil)

func NewSpeedUpOnDefeat(f *field.Field, others []FieldPusher) *SpeedUpOnDefeat {
	b := newBase(f)
	return &SpeedUpOnDefeat{
		base:   *b,
		others: others,
	}
}

type SpeedUpOnDefeat struct {
	base
	others []FieldPusher
}

func (s *SpeedUpOnDefeat) Analyze(events event.Reader) {
	var endMode *field.Mode
	events.Range(func(e event.Event) {
		if v, ok := e.(*op.FieldMode); ok && (v.ModeNew == field.ModeGameOver || v.ModeNew == field.ModeVictory || v.ModeNew == field.ModeDefeat) {
			endMode = &v.ModeNew
		}
	})
	if endMode == nil {
		return
	}

	if endMode == nil || *endMode != field.ModeDefeat {
		return
	}

	s.base.start()
}

func (s *SpeedUpOnDefeat) Sweep(event.Pusher) {
	for _, other := range s.others {
		for ctrlIdx := range byte(other.Field.Ctrls()) {
			if other.Field.CtrlStateIsTerminal(ctrlIdx) {
				continue
			}

			level := other.Field.CtrlLevel(ctrlIdx)
			if level <= 5 {
				other.Pusher.Push(op.NewPieceSpeedUp(ctrlIdx, 2))
			} else if level < piece.MaxLevel {
				other.Pusher.Push(op.NewPieceSpeedUp(ctrlIdx, 1))
			}
		}
	}
}
