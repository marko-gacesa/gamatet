// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package sweeper

import (
	"time"

	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/op"
	"github.com/marko-gacesa/gamatet/game/piece"
	"github.com/marko-gacesa/gamatet/logic/random"
)

var _ Sweeper = (*Magic)(nil)

const (
	magicWaitSeconds   = 5
	magicActiveSeconds = 15
)

type magicState byte

const (
	magicStateRunning magicState = iota
	magicStateActivated
	magicStateFinished
)

type MagicType byte

const (
	MagicTypeSelf   MagicType = 1
	MagicTypeOthers MagicType = 2
	MagicTypeAll    MagicType = MagicTypeSelf | MagicTypeOthers
)

const (
	magicEffectNone field.Effect = iota

	magicEffectLid
	magicEffectBigO
	magicRaise

	magicEffectPatchHoles

	magicEffectTotal
)

var (
	magicEffectsSelf   = []field.Effect{magicEffectPatchHoles}
	magicEffectsOthers = []field.Effect{magicEffectLid, magicEffectBigO, magicRaise}
)

func NewMagic(f *field.Field, others []FieldPusher, seed int, types MagicType) *Magic {
	b := newBase(f)

	m := &Magic{
		base:   *b,
		others: others,
		seed:   uint32(seed),
		state:  magicStateRunning,
		types:  types,
	}

	m.base.Start(nil) // it's always active

	return m
}

// Magic assumes that there is max one block type=Goal on the field.
// It creates it, monitors it and restores the original block when it expires.
type Magic struct {
	base
	others []FieldPusher
	seed   uint32
	state  magicState
	types  MagicType

	count    uint32
	oldBlock block.Block
}

func (s *Magic) Start(analyzer *Analyzer) bool {
	if analyzer.endMode != nil && s.state != magicStateFinished {
		s.state = magicStateFinished
		s.base.reschedule(time.Nanosecond)
		return false
	}

	if analyzer.blocks.goal > 0 {
		s.state = magicStateActivated
		s.base.reschedule(time.Microsecond)
	}

	return false
}

func (s *Magic) Sweep(p event.Pusher) {
	effect, seconds := s.field.GetEffect()

	if s.state == magicStateFinished {
		if effect != magicEffectNone {
			s.restoreBlock(p)
		}
		p.Push(op.NewFieldEffect(effect, magicEffectNone, 0, 0))
		s.endIteration()
		return
	}

	if s.state == magicStateActivated {
		s.activated(effect, p)
		s.state = magicStateRunning
		p.Push(op.NewFieldEffect(effect, magicEffectNone, seconds, magicWaitSeconds))
		s.base.reschedule(time.Second)
		return
	}

	if seconds > 0 {
		p.Push(op.NewFieldEffect(effect, effect, seconds, seconds-1))
		s.base.reschedule(time.Second)
		return
	}

	if effect != magicEffectNone {
		s.restoreBlock(p)
		p.Push(op.NewFieldEffect(effect, magicEffectNone, 0, magicWaitSeconds))
		s.base.reschedule(time.Second)
		return
	}

	xyb, ok := s.possessBlock()
	if !ok {
		p.Push(op.NewFieldEffect(magicEffectNone, magicEffectNone, 0, magicWaitSeconds))
		s.base.reschedule(time.Second)
		return
	}

	s.oldBlock = xyb.Block
	effect = s.randomEffect()
	if effect == magicEffectNone {
		s.state = magicStateFinished
		s.endIteration()
		return
	}

	p.Push(op.NewFieldBlockTransform(xyb.X, xyb.Y, xyb.Block, block.Goal, field.AnimNo, 0))
	p.Push(op.NewFieldExBlock(xyb.X, xyb.Y, field.AnimDestroy, 0, xyb.Block))
	p.Push(op.NewFieldEffect(magicEffectNone, effect, 0, magicActiveSeconds))

	s.count++
	s.base.reschedule(time.Second)
}

func (s *Magic) activated(effect field.Effect, p event.Pusher) {
	switch effect {
	case magicEffectLid:
		s.effectLid()
	case magicEffectBigO:
		s.effectBigO()
	case magicRaise:
		s.effectRaise()

	case magicEffectPatchHoles:
		s.effectPatchHoles(p)
	}
}

func (s *Magic) restoreBlock(p event.Pusher) {
	s.field.RangeBlocks(func(xyb block.XYB) bool {
		if xyb.Block.Type == block.TypeGoal {
			p.Push(op.NewFieldBlockTransform(xyb.X, xyb.Y, xyb.Block, s.oldBlock, field.AnimNo, 0))
			p.Push(op.NewFieldExBlock(xyb.X, xyb.Y, field.AnimDestroy, 0, block.Block{
				Type: block.TypeRock, Color: xyb.Color,
			}))
			return false
		}
		return true
	})
}

func (s *Magic) possessBlock() (block.XYB, bool) {
	var buffer [128]block.XYB
	blocks := buffer[:0]
	s.field.RangeBlocks(func(xyb block.XYB) bool {
		if xyb.Block.Type == block.TypeRock && xyb.Block.Hardness == 0 {
			blocks = append(blocks, xyb)
			if len(blocks) == cap(blocks) {
				return false
			}
		}
		return true
	})
	if len(blocks) == 0 {
		return block.XYB{}, false
	}

	randomIndex := random.New(s.count, s.seed)
	xyb := blocks[randomIndex.Int(len(blocks))]

	return xyb, true
}

func (s *Magic) randomEffect() field.Effect {
	var effectsBuffer [magicEffectTotal]field.Effect
	effects := effectsBuffer[:0]
	if s.types&MagicTypeSelf > 0 {
		effects = append(effects, magicEffectsSelf...)
	}
	if s.types&MagicTypeOthers > 0 {
		effects = append(effects, magicEffectsOthers...)
	}
	if len(effects) == 0 {
		return magicEffectNone
	}

	r := random.New(s.count, s.seed)
	return effects[r.Int(len(effects))]
}

func (s *Magic) effectLid() {
	for idx, o := range s.others {
		f := o.Field
		if f.IsFinished() {
			continue
		}

		r := random.New(s.count*10+uint32(idx), s.seed)

		w := f.GetWidth()
		h := f.GetHeight()

		var topRow int
		for col := range w {
			topRow = max(f.GetTopmostEmpty(col), topRow)
		}

		for j := range 2 {
			skipCol := r.Int(w)
			if topRow+j >= h {
				break
			}
			for col := range w {
				if col != skipCol {
					o.Pusher.Push(op.NewFieldBlockSet(col, topRow+j, op.TypeSet, field.AnimPop, 0, block.Rock))
				}
			}
		}
	}
}

func (s *Magic) effectBigO() {
	for _, o := range s.others {
		f := o.Field
		if f.IsFinished() {
			continue
		}

		ctrls := byte(f.Ctrls())
		for ctrlIdx := range ctrls {
			ctrl := f.Ctrl(ctrlIdx)
			for pieceCount := ctrl.PieceCount; ; pieceCount++ {
				if !ctrl.Feed.Overridden(pieceCount) {
					o.Pusher.Push(op.NewPieceOverride(ctrlIdx, piece.NewO(block.Rock), pieceCount))
					break
				}
			}
		}
	}
}

func (s *Magic) effectRaise() {
	for _, o := range s.others {
		f := o.Field
		if f.IsFinished() {
			continue
		}

		sections := f.FindMovableSections(func(f *field.Field, section field.ColumnSection) bool {
			return f.GetXY(section.Column, section.RowTo-1).Type == block.TypeEmpty
		})
		for _, section := range sections {
			o.Pusher.Push(op.NewFieldColumnShift(section, 1))
		}
	}
}

func (s *Magic) effectPatchHoles(p event.Pusher) {
	var holesBuffer [64]block.XY
	holes := holesBuffer[:0]

	w := s.field.GetWidth()

	for col := range w {
		top := s.field.GetTopmostEmpty(col) - 2
		for row := 0; row <= top; row++ {
			if s.field.GetXY(col, row).Type == block.TypeEmpty {
				holes = append(holes, block.XY{X: col, Y: row})
			}
		}
	}

	r := random.New(s.count, s.seed)
	random.Shuffle(r, holes)

	const patches = 10

	for i := range min(patches, len(holes)) {
		p.Push(op.NewFieldBlockSet(holes[i].X, holes[i].Y, op.TypeSet, field.AnimPop, 0, block.Rock))
	}
}
