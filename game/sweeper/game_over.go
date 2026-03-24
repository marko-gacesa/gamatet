// Copyright (c) 2025, 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package sweeper

import (
	"time"

	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/op"
	"github.com/marko-gacesa/gamatet/game/piece"
)

var _ Sweeper = (*GameOver)(nil)

func NewGameOver(f *field.Field) *GameOver {
	b := newBase(f)
	return &GameOver{
		base: *b,
	}
}

type gameOverMethod byte

const (
	gameOverMethodVanish gameOverMethod = iota
	gameOverMethodCurtain
	gameOverMethodBurn
	gameOverMethodFall
)

type GameOver struct {
	base
	state  *piece.State
	method gameOverMethod
}

func (s *GameOver) Start(analyzer *Analyzer) bool {
	if analyzer.endMode == nil {
		return false
	}

	switch *analyzer.endMode {
	case field.ModeGameOver:
		s.method = gameOverMethodCurtain
	case field.ModeVictory:
		s.method = gameOverMethodBurn
	case field.ModeDefeat:
		s.method = gameOverMethodFall
	default:
		return false
	}

	return s.base.Start(analyzer)
}

func (s *GameOver) Sweep(p event.Pusher) {
	switch s.method {
	case gameOverMethodVanish:
		s.blockVanish(p)
	case gameOverMethodCurtain:
		s.blockCurtain(p)
	case gameOverMethodBurn:
		s.blockBurn(p)
	case gameOverMethodFall:
		s.blockFall(p)
	}
}

func (s *GameOver) blockCurtain(p event.Pusher) {
	const n = 10

	xybs := s.findAllDestructible(n)
	if len(xybs) == 0 {
		s.endIteration()
		return
	}

	for _, xyb := range xybs {
		p.Push(&op.FieldBlockSet{
			Col:       byte(xyb.X),
			Row:       byte(xyb.Y),
			Op:        op.TypeClear,
			AnimType:  field.AnimCurtain,
			AnimParam: 1,
			Block:     xyb.Block,
		})
	}

	s.reschedule(2 * time.Millisecond)
}

func (s *GameOver) blockVanish(p event.Pusher) {
	const n = 4

	xybs := s.findAllDestructible(n)
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

func (s *GameOver) blockBurn(p event.Pusher) {
	const n = 4

	xybs := s.findAllDestructible(n)
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

	s.reschedule(75 * time.Millisecond)
}

func (s *GameOver) blockFall(p event.Pusher) {
	const n = 4

	xybs := s.findAllDestructible(n)
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

func (s *GameOver) findAllDestructible(max int) []block.XYB {
	result := make([]block.XYB, 0, max)

	s.field.RangeBlocks(func(xyb block.XYB) bool {
		if t := xyb.Block.Type; t != block.TypeRock && t != block.TypeRuby {
			return true
		}

		result = append(result, xyb)
		return len(result) < max
	})

	return result
}
