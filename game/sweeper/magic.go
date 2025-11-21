// Copyright (c) 2025 by Marko Gaćeša

package sweeper

import (
	"fmt"
	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/op"
	"github.com/marko-gacesa/gamatet/logic/random"
	"time"
)

var _ Sweeper = (*Magic)(nil)

const (
	magicDurWait = 2 * time.Second
	magicDur     = 5 * time.Second

	magicEffectCount = 5
)

type magicState byte

const (
	magicStateWaiting magicState = iota
	magicStateRunning
	magicStateActivated
)

func NewMagic(f *field.Field, others []FieldPusher, seed int) *Magic {
	b := newBase(f)

	m := &Magic{
		base:   *b,
		others: others,
		seed:   seed,
		state:  magicStateWaiting,
	}

	m.base.Start(nil) // it's always active
	m.waiting()

	return m
}

// Magic assumes that there is max one block type=Goal on the field.
// It creates it, monitors it and restores the original block when it expires.
type Magic struct {
	base
	others []FieldPusher
	seed   int
	state  magicState

	count    uint32
	oldBlock block.Block
	effect   byte
}

func (s *Magic) Start(analyzer *Analyzer) bool {
	if analyzer.blocks.goal > 0 {
		s.activate()
	}

	return false
}

func (s *Magic) Sweep(p event.Pusher) {
	switch s.state {
	case magicStateWaiting:
		s.create(p)
	case magicStateRunning:
		s.expired(p)
	case magicStateActivated:
		s.activated(p)
	}
}

func (s *Magic) activate() {
	s.state = magicStateActivated
	s.base.reschedule(time.Microsecond)
}

func (s *Magic) activated(p event.Pusher) {
	fmt.Println("ACTIVATED", s.effect)
	switch s.effect {
	// TODO: Finish implementation of magic activation effects
	default:
		for _, o := range s.others {
			o.Pusher.Push(op.NewFieldExBlock(0, 0, field.AnimDestroy, 0, block.Rock))
		}
	}
	s.waiting()
}

func (s *Magic) expired(p event.Pusher) {
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
	s.state = magicStateWaiting
	s.base.reschedule(magicDurWait)
}

func (s *Magic) create(p event.Pusher) {
	var buffer [64]block.XYB
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
		s.waiting()
		return
	}

	randomIndex := random.New(s.count, uint32(s.seed))

	xyb := blocks[randomIndex.Int(len(blocks))]

	s.oldBlock = xyb.Block
	p.Push(op.NewFieldBlockTransform(xyb.X, xyb.Y, xyb.Block, block.Goal, field.AnimNo, 0))
	p.Push(op.NewFieldExBlock(xyb.X, xyb.Y, field.AnimDestroy, 0, xyb.Block))

	randomEffect := random.New(s.count, uint32(s.seed))
	s.effect = byte(randomEffect.Int(magicEffectCount))

	s.count++
	s.state = magicStateRunning
	s.base.reschedule(magicDur)
}

func (s *Magic) waiting() {
	s.state = magicStateWaiting
	s.base.reschedule(magicDurWait)
}
