// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package sweeper

import (
	"math"
	"time"

	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/op"
	"github.com/marko-gacesa/gamatet/logic/random"
)

var _ Sweeper = (*Gnaw)(nil)

func NewGnaw(f *field.Field, seed int) *Gnaw {
	b := newBase(f)

	g := &Gnaw{
		base:     *b,
		seed:     uint(seed),
		aliveMap: make(map[uint32]*gnawData),
		killed:   0,
	}

	return g
}

type Gnaw struct {
	base
	seed     uint
	aliveMap map[uint32]*gnawData
	killed   int
}

type gnawData struct {
	color        uint32
	lastPos      block.XY
	hungry       int
	restDuration time.Duration
	lastAction   time.Time
}

// Spawn is a temporary function.
// TODO: TEMPORARY
func (g *Gnaw) Spawn(p event.Pusher) {
	g.spawn(0x007F000FF, math.MaxInt, time.Second, p)
	g.scheduleNext()
}

func (g *Gnaw) Start(a *Analyzer) bool {
	return false
}

func (g *Gnaw) Sweep(p event.Pusher) {
	for _, gnaw := range g.aliveMap {
		if g.process(gnaw, p) {
			break // at most one gnaw can be processed in one iteration
		}
	}

	g.scheduleNext()
}

func (g *Gnaw) spawn(color uint32, hungry int, restDuration time.Duration, p event.Pusher) bool {
	if _, ok := g.aliveMap[color]; ok {
		return false
	}

	r := random.New(uint(color), g.seed)
	f := g.base.field
	w := f.GetWidth()
	h := f.GetHeight()

	potentialXY := make([]block.XY, 0, 10)
	f.FindNearest8(block.XY{X: w / 2, Y: h / 3}, 2*w/3, func(xyb block.XYB, i int) bool {
		if xyb.Block.Type != block.TypeEmpty {
			return false
		}

		potentialXY = append(potentialXY, xyb.XY)
		return len(potentialXY) == 10
	})

	if len(potentialXY) == 0 {
		return false
	}

	n := r.Int(len(potentialXY))
	xy := potentialXY[n]

	b := block.Block{
		Type:     block.TypeGnaw,
		Hardness: 0,
		Color:    color,
	}

	gnaw := &gnawData{
		color:        color,
		lastPos:      xy,
		hungry:       hungry,
		restDuration: restDuration,
		lastAction:   time.Now(),
	}

	g.aliveMap[color] = gnaw

	p.Push(op.NewFieldBlockSet(xy.X, xy.Y, op.TypeSet, field.AnimPop, 0, b))

	return true
}

// process processes one gnaw. It returns true if it was processed (generated at least one event).
// Return value false means that it's not yet time for this gnaw to be processed.
func (g *Gnaw) process(gnaw *gnawData, p event.Pusher) (processed bool) {
	now := time.Now()

	ok := g.find(gnaw)
	if !ok {
		g.killed++
		g.remove(gnaw, false, p) // not in field
		return true
	}

	if gnaw.hungry == 0 {
		g.remove(gnaw, true, p) // not hungry anymore
		return true
	}

	if now.Before(gnaw.lastAction.Add(gnaw.restDuration)) {
		return false
	}

	gnaw.lastAction = now
	xy := gnaw.lastPos

	f := g.base.field
	w := f.GetWidth()
	h := f.GetHeight()

	r := max(w-xy.X-1, xy.X, h-xy.Y-1, xy.Y)
	target, ok := f.FindNearest8(xy, r, func(xyb block.XYB, i int) bool {
		return xyb.Block.Type.Gnawable() && f.HasLOS(xy, xyb.XY)
	})
	if !ok {
		g.move1(gnaw, p)
		return true
	}

	path := f.Path8(xy, target.XY)
	if len(path) < 2 {
		g.move1(gnaw, p)
		return true
	}

	moveTo := path[1]

	b := f.GetXY(moveTo.X, moveTo.Y)
	switch {
	case b.Type == block.TypeEmpty:
		// move to empty space
		p.Push(op.NewFieldBlockSwap(xy.X, xy.Y, moveTo.X, moveTo.Y, field.AnimSlide, 0))
		gnaw.lastPos = moveTo
	case b.Type.Gnawable():
		if b.Hardness == 0 {
			// eat the block and move there
			p.Push(op.NewFieldBlockSet(moveTo.X, moveTo.Y, op.TypeClear, field.AnimDestroy, 0, b))
			p.Push(op.NewFieldBlockSwap(xy.X, xy.Y, moveTo.X, moveTo.Y, field.AnimSlide, 0))
			gnaw.lastPos = moveTo
		} else {
			// reduce the block's hardness
			p.Push(op.NewFieldBlockHardness(moveTo.X, moveTo.Y, -1, field.AnimSpin, 0))
		}

		gnaw.hungry--
	default:
		g.move1(gnaw, p)
	}

	return true
}

func (g *Gnaw) move1(gnaw *gnawData, p event.Pusher) {
	f := g.base.field
	neighbors := f.Neighbors8(gnaw.lastPos)

	potentialXY := make([]block.XY, 0, 4)
	neighbors.ForEach(f, gnaw.lastPos, func(xyb block.XYB) {
		if xyb.Block.Type == block.TypeEmpty {
			potentialXY = append(potentialXY, xyb.XY)
		}
	})

	if len(potentialXY) == 0 {
		g.remove(gnaw, false, p) // unable to move
		return
	}

	r := random.New(uint(gnaw.color)+uint(gnaw.lastPos.Y*f.GetWidth()+gnaw.lastPos.X), g.seed)
	n := r.Int(len(potentialXY))

	moveTo := potentialXY[n]

	p.Push(op.NewFieldBlockSwap(gnaw.lastPos.X, gnaw.lastPos.Y, moveTo.X, moveTo.Y, field.AnimSlide, 0))
	gnaw.lastPos = moveTo
}

func (g *Gnaw) find(gnaw *gnawData) bool {
	x := gnaw.lastPos.X
	y := gnaw.lastPos.Y

	// first try the last known location

	b := g.base.field.GetXY(x, y)
	if b.Type == block.TypeGnaw && b.Color == gnaw.color {
		return true
	}

	// next try few locations below the last known location

	for i := 1; i < 4; i++ {
		y1 := y - i
		if y1 < 0 {
			break
		}

		b = g.base.field.GetXY(x, y1)
		if b.Type == block.TypeGnaw && b.Color == gnaw.color {
			gnaw.lastPos.Y = y1
			return true
		}
	}

	// finally search the entire field

	var found bool
	g.field.RangeBlocks(func(xyb block.XYB) bool {
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

func (g *Gnaw) scheduleNext() {
	if len(g.aliveMap) == 0 {
		g.base.endIteration()
		return
	}

	var activateAtMin time.Time
	for _, q := range g.aliveMap {
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

	g.reschedule(activateDelay)
}

func (g *Gnaw) remove(gnaw *gnawData, noHunger bool, p event.Pusher) {
	x, y := gnaw.lastPos.X, gnaw.lastPos.Y
	b := g.base.field.GetXY(x, y)

	if b.Type == block.TypeGnaw && b.Color == gnaw.color {
		animGnaw := field.AnimDestroy
		if noHunger {
			animGnaw = field.AnimPop
		}

		p.Push(op.NewFieldBlockSet(x, y, op.TypeClear, animGnaw, 0, b))
	}

	delete(g.aliveMap, gnaw.color)
}
