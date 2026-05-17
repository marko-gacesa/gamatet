// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package sweeper

import (
	"time"

	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/op"
)

var _ Sweeper = (*GnawKeeper)(nil)

func NewGnawKeeper(f *field.Field) *GnawKeeper {
	b := newBase(f)

	g := &GnawKeeper{
		base:     *b,
		aliveMap: make(map[uint32]*gnawData),
	}

	return g
}

type GnawKeeper struct {
	base
	aliveMap    map[uint32]*gnawData
	toRemoveMap map[uint32]*gnawData
	killed      int
}

type gnawData struct {
	color        uint32
	lastPos      block.XY
	hunger       int
	restDuration time.Duration
	lastAction   time.Time
}

func (s *GnawKeeper) Analyze(events event.Reader) {
	var removedBuffer [4]block.XYB
	removed := removedBuffer[:0]

	var addedBuffer [4]gnawData
	added := addedBuffer[:0]

	events.Range(func(e event.Event) {
		switch v := e.(type) {
		case *op.FieldBlockSet:
			if v.Block.Type == block.TypeGnaw {
				switch v.Op {
				case op.TypeSet:
					panic("gnaws must not be added with op.FieldBlockSet")
				case op.TypeClear:
					removed = append(removed, block.XYB{
						XY:    block.XY{X: int(v.Col), Y: int(v.Row)},
						Block: v.Block,
					})
				}
			}
		case *op.FieldDestroyRow:
			for col, b := range v.Blocks {
				if b.Type == block.TypeGnaw {
					removed = append(removed, block.XYB{
						XY:    block.XY{X: int(col), Y: int(v.Row)},
						Block: b,
					})
				}
			}
		case *op.FieldDestroyColumn:
			if v.Block.Type == block.TypeGnaw {
				removed = append(removed, block.XYB{
					XY:    block.XY{X: int(v.Col), Y: int(v.Row)},
					Block: v.Block,
				})
			}
		case *op.FieldGnaw:
			added = append(added, gnawData{
				color:        v.Color,
				lastPos:      block.XY{X: int(v.Col), Y: int(v.Row)},
				hunger:       int(v.Hunger),
				restDuration: time.Millisecond * time.Duration(v.Delay),
				lastAction:   time.Now(),
			})
		}
	})

	if len(added)+len(removed) == 0 {
		return
	}

	for _, g := range added {
		s.aliveMap[g.color] = &g
	}

	s.base.start()
}

func (s *GnawKeeper) Sweep(p event.Pusher) {
	for _, gnaw := range s.aliveMap {
		if s.process(gnaw, p) {
			break // at most one gnaw can be processed in one iteration
		}
	}

	s.scheduleNext()
}

// process processes one gnaw. It returns true if it was processed (generated at least one event).
// Return value false means that it's not yet time for this gnaw to be processed.
func (s *GnawKeeper) process(gnaw *gnawData, p event.Pusher) (processed bool) {
	now := time.Now()

	ok := s.find(gnaw)
	if !ok {
		s.killed++
		s.remove(gnaw) // not in field
		return true
	}

	if gnaw.hunger == 0 {
		s.kill(gnaw, true, p) // not hungry anymore
		return true
	}

	if now.Before(gnaw.lastAction.Add(gnaw.restDuration)) {
		return false
	}

	gnaw.lastAction = now
	xy := gnaw.lastPos

	f := s.base.field
	w := f.GetWidth()
	h := f.GetHeight()

	r := max(w-xy.X-1, xy.X, h-xy.Y-1, xy.Y)
	target, ok := f.FindNearest8(xy, r, func(xyb block.XYB, i int) bool {
		return validGnawFood(xyb.Block.Type) && f.HasLOS(xy, xyb.XY)
	})
	if !ok {
		s.move1(gnaw, p) // no target to eat
		return true
	}

	path := f.Path8(xy, target.XY, validGnawMove)
	if len(path) < 2 {
		s.move1(gnaw, p) // no route to target
		return true
	}

	moveTo := path[1]

	b := f.GetXY(moveTo.X, moveTo.Y)
	switch b.Type {
	case block.TypeEmpty:
		// move to empty space
		p.Push(op.NewFieldBlockSwap(xy.X, xy.Y, moveTo.X, moveTo.Y, field.AnimSlide, 0))
		gnaw.lastPos = moveTo
	case block.TypeRock:
		if b.Hardness == 0 {
			// eat the block and move there
			p.Push(op.NewFieldBlockSet(moveTo.X, moveTo.Y, op.TypeClear, field.AnimDestroy, 0, b))
			p.Push(op.NewFieldBlockSwap(xy.X, xy.Y, moveTo.X, moveTo.Y, field.AnimSlide, 0))
			gnaw.lastPos = moveTo
			gnaw.hunger--
		} else {
			// reduce the block's hardness
			p.Push(op.NewFieldBlockHardness(moveTo.X, moveTo.Y, -1, field.AnimSpin, 0))
		}
	default:
		s.move1(gnaw, p)
	}

	return true
}

func (s *GnawKeeper) move1(gnaw *gnawData, p event.Pusher) {
	f := s.base.field
	neighbors := f.Neighbors8(gnaw.lastPos, validGnawMove)

	potentialXY := make([]block.XY, 0, 4)
	neighbors.ForEach(f, gnaw.lastPos, func(xyb block.XYB) {
		if xyb.Block.Type == block.TypeEmpty {
			potentialXY = append(potentialXY, xyb.XY)
		}
	})

	if len(potentialXY) == 0 {
		s.kill(gnaw, false, p) // unable to move
		return
	}

	r := f.Random(uint64(gnaw.color) + uint64(gnaw.lastPos.Y*f.GetWidth()+gnaw.lastPos.X))
	n := r.Int(len(potentialXY))

	moveTo := potentialXY[n]

	p.Push(op.NewFieldBlockSwap(gnaw.lastPos.X, gnaw.lastPos.Y, moveTo.X, moveTo.Y, field.AnimSlide, 0))
	gnaw.lastPos = moveTo
}

func (s *GnawKeeper) find(gnaw *gnawData) bool {
	x := gnaw.lastPos.X
	y := gnaw.lastPos.Y

	// first try the last known location

	b := s.base.field.GetXY(x, y)
	if b.Type == block.TypeGnaw && b.Color == gnaw.color {
		return true
	}

	// next try few locations below the last known location

	for i := 1; i < 4; i++ {
		y1 := y - i
		if y1 < 0 {
			break
		}

		b = s.base.field.GetXY(x, y1)
		if b.Type == block.TypeGnaw && b.Color == gnaw.color {
			gnaw.lastPos.Y = y1
			return true
		}
	}

	// finally search the entire field

	var found bool
	s.field.RangeBlocks(func(xyb block.XYB) bool {
		if b.Type == block.TypeGnaw && b.Color == gnaw.color {
			found = true
			gnaw.lastPos = xyb.XY
			return false
		}

		return true
	})
	if found {
		return true
	}

	return false
}

func (s *GnawKeeper) scheduleNext() {
	if len(s.aliveMap) == 0 {
		s.base.endIteration()
		return
	}

	var activateAtMin time.Time
	for _, q := range s.aliveMap {
		activateAt := q.lastAction.Add(q.restDuration)
		if activateAtMin.IsZero() {
			activateAtMin = activateAt
		} else if activateAtMin.After(activateAt) {
			activateAtMin = activateAt
		}
	}

	now := time.Now()
	activateDelay := activateAtMin.Sub(now)
	if activateDelay <= 0 {
		activateDelay = time.Nanosecond
	}

	s.reschedule(activateDelay)
}

func (s *GnawKeeper) kill(gnaw *gnawData, noHunger bool, p event.Pusher) {
	x, y := gnaw.lastPos.X, gnaw.lastPos.Y
	b := s.base.field.GetXY(x, y)

	if b.Type == block.TypeGnaw && b.Color == gnaw.color {
		animGnaw := field.AnimDestroy
		if noHunger {
			animGnaw = field.AnimPop
		}

		p.Push(op.NewFieldBlockSet(x, y, op.TypeClear, animGnaw, 0, b))
	}

	s.remove(gnaw)
}

func (s *GnawKeeper) remove(gnaw *gnawData) {
	delete(s.aliveMap, gnaw.color)
}

// validGnawFood returns true if a gnaw can eat this block type.
func validGnawFood(t block.Type) bool { return t == block.TypeRock }

// validGnawMove returns true if a gnaw can move into this block type.
func validGnawMove(t block.Type) bool { return t == block.TypeEmpty || t == block.TypeRock }
