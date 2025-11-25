// Copyright (c) 2025 by Marko Gaćeša
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

func (s *SpeedUpOnDefeat) Start(analyzer *Analyzer) bool {
	if analyzer.endMode == nil || *analyzer.endMode != field.ModeDefeat {
		return false
	}

	return s.base.Start(analyzer)
}

func (s *SpeedUpOnDefeat) Sweep(event.Pusher) {
	for _, other := range s.others {
		for ctrlIdx := range byte(other.Field.Ctrls()) {
			ctrl := other.Field.Ctrl(ctrlIdx)
			if ctrl.State.IsTerminal() {
				continue
			}

			level := ctrl.Level
			if level <= 5 {
				other.Pusher.Push(op.NewPieceSpeedUp(ctrlIdx, 2))
			} else if level < piece.MaxLevel {
				other.Pusher.Push(op.NewPieceSpeedUp(ctrlIdx, 1))
			}
		}
	}
}
