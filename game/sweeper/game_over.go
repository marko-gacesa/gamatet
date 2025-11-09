// Copyright (c) 2025 by Marko Gaćeša

package sweeper

import (
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/op"
	"github.com/marko-gacesa/gamatet/game/piece"
	"time"
)

var _ Sweeper = (*GameOver)(nil)

func NewGameOver(f *field.Field) *GameOver {
	b := newBase(f)
	return &GameOver{
		base: *b,
	}
}

type GameOver struct {
	animFn func(*base, event.Pusher)
	state  *piece.State
	base
}

func (s *GameOver) Start(analyzer *Analyzer) bool {
	if analyzer.endMode == nil {
		return false
	}

	switch *analyzer.endMode {
	case field.ModeGameOver:
		s.animFn = blockVanish
	case field.ModeVictory:
		s.animFn = blockBurn
	case field.ModeDefeat:
		s.animFn = blockFall
	default:
		return false
	}

	return s.base.Start(analyzer)
}

func (s *GameOver) Sweep(p event.Pusher) {
	s.animFn(&s.base, p)
}

func blockBurn(s *base, p event.Pusher) {
	const n = 4

	xybs := s.field.AllXY(n)
	if len(xybs) == 0 {
		s.endIteration()
		return
	}

	for _, xyb := range xybs {
		p.Push(&op.FieldBlockSet{
			Col:       byte(xyb.X),
			Row:       byte(xyb.Y),
			Op:        op.TypeClear,
			AnimType:  field.AnimDestroy,
			AnimParam: 0,
			Block:     xyb.Block,
		})
	}

	s.reschedule(50 * time.Millisecond)
}

func blockFall(s *base, p event.Pusher) {
	const n = 4

	xybs := s.field.AllXY(n)
	if len(xybs) == 0 {
		s.endIteration()
		return
	}

	for _, xyb := range xybs {
		p.Push(&op.FieldBlockSet{
			Col:       byte(xyb.X),
			Row:       byte(xyb.Y),
			Op:        op.TypeClear,
			AnimType:  field.AnimFall,
			AnimParam: byte(xyb.Y),
			Block:     xyb.Block,
		})
	}

	s.reschedule(50 * time.Millisecond)
}

func blockVanish(s *base, p event.Pusher) {
	const n = 4

	xybs := s.field.AllXY(n)
	if len(xybs) == 0 {
		s.endIteration()
		return
	}

	for _, xyb := range xybs {
		p.Push(&op.FieldBlockSet{
			Col:       byte(xyb.X),
			Row:       byte(xyb.Y),
			Op:        op.TypeClear,
			AnimType:  field.AnimPop,
			AnimParam: 0,
			Block:     xyb.Block,
		})
	}

	s.reschedule(50 * time.Millisecond)
}
