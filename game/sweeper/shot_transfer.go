// Copyright (c) 2025, 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package sweeper

import (
	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/op"
)

var _ Sweeper = (*ShotTransfer)(nil)

func NewShotTransfer(f *field.Field, others []FieldPunisher) *ShotTransfer {
	b := newBase(f)
	return &ShotTransfer{
		base:   *b,
		others: others,
	}
}

type ShotTransfer struct {
	base
	others []FieldPunisher
	shots  []block.Type
}

func (s *ShotTransfer) Start(analyzer *Analyzer) bool {
	if len(analyzer.shots) == 0 {
		return false
	}

	s.shots = analyzer.shots

	return s.base.Start(analyzer)
}

func (s *ShotTransfer) Sweep(event.Pusher) {
	for i := range s.others {
		f := s.others[i].Field
		if f.IsFinished() {
			continue
		}

		for _, sh := range s.shots {
			switch sh {
			case block.TypeAcid:
				h := f.GetHeight()
				xys := f.Blizzard(1)
				if len(xys) == 0 {
					continue
				}

				s.others[i].Pusher.Push(&op.FieldBlockSet{
					Col:       byte(xys[0].X),
					Row:       byte(xys[0].Y),
					Op:        op.TypeSet,
					AnimType:  field.AnimFall,
					AnimParam: byte(h - xys[0].Y),
					Block:     block.Rock,
				})
			default:
				xyb, ok := f.GetRandomBlock()
				if !ok {
					continue
				}

				s.others[i].Pusher.Push(&op.FieldBlockSet{
					Col:      byte(xyb.X),
					Row:      byte(xyb.Y),
					Op:       op.TypeClear,
					AnimType: field.AnimDestroy,
					Block:    xyb.Block,
				})
			}
		}
	}

	s.endIteration()
}
