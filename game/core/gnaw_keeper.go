// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package core

import (
	"fmt"
	"math"
	"math/rand/v2"
	"time"

	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/op"
)

func NewGnawKeeper(f *field.Field, seed uint) *GnawKeeper {
	k := &GnawKeeper{
		f:         f,
		gnaws:     nil,
		timer:     time.NewTimer(24 * time.Hour),
		seed:      seed,
		killCount: 0,
	}

	k.timer.Stop()

	return k
}

type GnawKeeper struct {
	f         *field.Field
	gnaws     map[uint32]*gnawData
	timer     *time.Timer
	seed      uint
	killCount int
}

type gnawData struct {
	color        uint32
	lastPos      block.XY
	hunger       int
	restDuration time.Duration
	lastAction   time.Time
}

func (k *GnawKeeper) StopTimer() {
	k.timer.Stop()
}

func (k *GnawKeeper) StartTimer() {
	d := k.getDelayUntilNextIter(time.Now())
	if d == 0 {
		return
	}

	k.timer.Reset(d)
}

func (k *GnawKeeper) Chan() <-chan time.Time {
	return k.timer.C
}

func (k *GnawKeeper) Exists(color uint32) bool {
	_, exists := k.gnaws[color]
	return exists
}

func (k *GnawKeeper) RandomColor() uint32 {
	for {
		color := rand.N[uint32](192)<<24 + rand.N[uint32](256)<<16 + rand.N[uint32](192)<<8 + 0xFF
		if _, ok := k.gnaws[color]; ok {
			continue
		}

		return color
	}
}

func (k *GnawKeeper) AddSmall(x, y int) {
	k.Add(x, y, k.RandomColor(), 10, time.Second)
}

func (k *GnawKeeper) AddInfinite(x, y int) {
	k.Add(x, y, k.RandomColor(), math.MaxInt, 1200*time.Millisecond)
}

func (k *GnawKeeper) Add(x, y int, color uint32, hunger int, restDuration time.Duration) {
	if _, ok := k.gnaws[color]; ok {
		panic(fmt.Sprintf("can't add gnaw; %X already exists", color))
	}

	if b := k.f.GetXY(x, y); b.Type != block.TypeEmpty {
		panic(fmt.Sprintf("can't add gnaw; (%d,%d) contains %d", x, y, b.Type))
	}

	gnaw := &gnawData{
		color:        color,
		lastPos:      block.XY{X: x, Y: y},
		hunger:       hunger,
		restDuration: restDuration,
		lastAction:   time.Now(),
	}

	if k.gnaws == nil {
		k.gnaws = make(map[uint32]*gnawData)
	}
	k.gnaws[color] = gnaw

	k.f.SetXY(x, y, field.AnimPop, 0, block.Block{
		Type:     block.TypeGnaw,
		Hardness: 0,
		Color:    color,
	})

	k.StartTimer()
}

func (k *GnawKeeper) UndoAdd(color uint32, animType, animParam int) {
	gnaw, ok := k.gnaws[color]
	if !ok {
		panic("failed to remove gnaw: gnaw not found")
	}

	if found := k.find(gnaw); found {
		b := k.f.ClearXY(gnaw.lastPos.X, gnaw.lastPos.Y, animType, animParam)
		if b.Type != block.TypeGnaw || b.Color != color {
			panic("failed to remove gnaw: ClearXY didn't the correct block")
		}
	}

	delete(k.gnaws, color)
}

func (k *GnawKeeper) ProcessAll(p event.Pusher) {
	now := time.Now()

	for _, gnaw := range k.gnaws {
		if k.process(gnaw, now, p) {
			break // at most one gnaw can be processed in one iteration
		}
	}

	d := k.getDelayUntilNextIter(now)
	if d == 0 {
		return
	}

	k.timer.Reset(d)
}

func (k *GnawKeeper) getDelayUntilNextIter(now time.Time) time.Duration {
	if len(k.gnaws) == 0 {
		return 0
	}

	var activateAtMin time.Time
	for _, q := range k.gnaws {
		activateAt := q.lastAction.Add(q.restDuration)
		if activateAtMin.IsZero() {
			activateAtMin = activateAt
		} else if activateAtMin.After(activateAt) {
			activateAtMin = activateAt
		}
	}

	activateDelay := activateAtMin.Sub(now)
	if activateDelay <= 0 {
		activateDelay = time.Nanosecond
	}

	return activateDelay
}

func (k *GnawKeeper) process(gnaw *gnawData, now time.Time, p event.Pusher) (success bool) {
	ok := k.find(gnaw)
	if !ok {
		k.killCount++
		k.remove(gnaw) // not in field
		return true
	}

	if gnaw.hunger == 0 {
		k.die(gnaw, true, p) // not hungry anymore
		return true
	}

	if now.Before(gnaw.lastAction.Add(gnaw.restDuration)) {
		return false
	}

	gnaw.lastAction = now
	xy := gnaw.lastPos

	f := k.f
	w := f.GetWidth()
	h := f.GetHeight()

	r := max(w-xy.X-1, xy.X, h-xy.Y-1, xy.Y)
	target, ok := f.FindNearest8(xy, r, func(xyb block.XYB, i int) bool {
		return validGnawFood(xyb.Block.Type) && f.HasLOS(xy, xyb.XY)
	})
	if !ok {
		k.move1(gnaw, p) // target to eat
		return true
	}

	path := f.Path8(xy, target.XY, validGnawMove)
	if len(path) < 2 {
		k.move1(gnaw, p) // no route to target
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
		k.move1(gnaw, p)
	}

	return true
}

func (k *GnawKeeper) move1(gnaw *gnawData, p event.Pusher) {
	f := k.f
	neighbors := f.Neighbors8(gnaw.lastPos, validGnawMove)

	potentialXY := make([]block.XY, 0, 4)
	neighbors.ForEach(f, gnaw.lastPos, func(xyb block.XYB) {
		if xyb.Block.Type == block.TypeEmpty {
			potentialXY = append(potentialXY, xyb.XY)
		}
	})

	if len(potentialXY) == 0 {
		k.die(gnaw, false, p) // unable to move
		return
	}

	r := f.Random(uint64(gnaw.color) + uint64(gnaw.lastPos.Y*f.GetWidth()+gnaw.lastPos.X))
	n := r.Int(len(potentialXY))

	moveTo := potentialXY[n]

	p.Push(op.NewFieldBlockSwap(gnaw.lastPos.X, gnaw.lastPos.Y, moveTo.X, moveTo.Y, field.AnimSlide, 0))
	gnaw.lastPos = moveTo
}

func (k *GnawKeeper) find(gnaw *gnawData) bool {
	f := k.f
	x := gnaw.lastPos.X
	y := gnaw.lastPos.Y

	// first try the last known location

	b := f.GetXY(x, y)
	if b.Type == block.TypeGnaw && b.Color == gnaw.color {
		return true
	}

	// next try few locations below the last known location

	for i := 1; i < 4; i++ {
		y1 := y - i
		if y1 < 0 {
			break
		}

		b = f.GetXY(x, y1)
		if b.Type == block.TypeGnaw && b.Color == gnaw.color {
			gnaw.lastPos.Y = y1
			return true
		}
	}

	// finally search the entire field

	var found bool
	f.RangeBlocks(func(xyb block.XYB) bool {
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

func (k *GnawKeeper) die(gnaw *gnawData, noHunger bool, p event.Pusher) {
	x, y := gnaw.lastPos.X, gnaw.lastPos.Y
	b := k.f.GetXY(x, y)

	if b.Type == block.TypeGnaw && b.Color == gnaw.color {
		animGnaw := field.AnimDestroy
		if noHunger {
			animGnaw = field.AnimPop
		}

		p.Push(op.NewFieldBlockSet(x, y, op.TypeClear, animGnaw, 0, b))
	}

	k.remove(gnaw)
}

func (k *GnawKeeper) remove(gnaw *gnawData) {
	delete(k.gnaws, gnaw.color)
}

// validGnawFood returns true if a gnaw can eat this block type.
func validGnawFood(t block.Type) bool { return t == block.TypeRock }

// validGnawMove returns true if a gnaw can move into this block type.
func validGnawMove(t block.Type) bool { return t == block.TypeEmpty || t == block.TypeRock }
