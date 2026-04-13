// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package field

import (
	"container/heap"
	"math"
	"slices"

	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/logic/geommath"
)

// HasLOS return true if between the two coordinates are only empty blocks.
func (f *Field) HasLOS(p0, p1 block.XY) bool {
	los := true
	geommath.Line(geommath.P[int](p0), geommath.P[int](p1), func(p geommath.P[int]) bool {
		if p == geommath.P[int](p0) {
			return true
		}
		if p == geommath.P[int](p1) {
			return false
		}

		los = los && f.GetXY(p.X, p.Y).Type == block.TypeEmpty

		return los
	})

	return los
}

const (
	Neighbor8BottomLeft = iota
	Neighbor8BottomMid
	Neighbor8BottomRight
	Neighbor8MidLeft
	Neighbor8MidRight
	Neighbor8TopLeft
	Neighbor8TopMid
	Neighbor8TopRight
)

type Neighbors8 [8]bool

func (f *Field) Neighbors8(pos block.XY) Neighbors8 {
	w := f.w
	h := f.h
	idx := pos.Y*w + pos.X

	var (
		hasTop    = pos.Y < h-1
		hasBottom = pos.Y > 0
		hasRight  = pos.X < w-1
		hasLeft   = pos.X > 0
	)

	isOk := func(f *Field, i int) bool {
		t := f.blocks[i].Block.Type
		return t == block.TypeEmpty || t.Gnawable()
	}

	var result Neighbors8

	if hasBottom {
		i := idx - w
		result[Neighbor8BottomLeft] = hasLeft && isOk(f, i-1)
		result[Neighbor8BottomMid] = isOk(f, i)
		result[Neighbor8BottomRight] = hasRight && isOk(f, i+1)
	}

	result[Neighbor8MidLeft] = hasLeft && isOk(f, idx-1)
	result[Neighbor8MidRight] = hasRight && isOk(f, idx+1)

	if hasTop {
		i := idx + w
		result[Neighbor8TopLeft] = hasLeft && isOk(f, i-1)
		result[Neighbor8TopMid] = isOk(f, i)
		result[Neighbor8TopRight] = hasRight && isOk(f, i+1)
	}

	return result
}

func (n Neighbors8) ForEach(f *Field, pos block.XY, fn func(xyb block.XYB)) {
	for idx, ok := range n {
		if !ok {
			continue
		}

		x := pos.X
		y := pos.Y
		switch idx {
		case Neighbor8BottomLeft, Neighbor8MidLeft, Neighbor8TopLeft:
			x--
		case Neighbor8BottomRight, Neighbor8MidRight, Neighbor8TopRight:
			x++
		}
		switch idx {
		case Neighbor8BottomLeft, Neighbor8BottomMid, Neighbor8BottomRight:
			y--
		case Neighbor8TopLeft, Neighbor8TopMid, Neighbor8TopRight:
			y++
		}

		xyb := block.XYB{
			XY:    block.XY{X: x, Y: y},
			Block: f.GetXY(x, y),
		}

		fn(xyb)
	}
}

const (
	Neighbor4Bottom = iota
	Neighbor4Left
	Neighbor4Right
	Neighbor4Top
)

type Neighbors4 [4]bool

func (f *Field) Neighbors4(pos block.XY) Neighbors4 {
	w := f.w
	h := f.h
	idx := pos.Y*w + pos.X

	var (
		hasTop    = pos.Y < h-1
		hasBottom = pos.Y > 0
		hasRight  = pos.X < w-1
		hasLeft   = pos.X > 0
	)

	isOk := func(f *Field, i int) bool {
		t := f.blocks[i].Block.Type
		return t == block.TypeEmpty || t.Gnawable()
	}

	var result Neighbors4

	result[Neighbor4Bottom] = hasBottom && isOk(f, idx-w)
	result[Neighbor4Left] = hasLeft && isOk(f, idx-1)
	result[Neighbor4Right] = hasRight && isOk(f, idx+1)
	result[Neighbor4Top] = hasTop && isOk(f, idx+w)

	return result
}

func (n Neighbors4) ForEach(f *Field, pos block.XY, fn func(xyb block.XYB)) {
	for idx, ok := range n {
		if !ok {
			continue
		}

		x := pos.X
		y := pos.Y
		switch idx {
		case Neighbor4Bottom:
			y--
		case Neighbor4Left:
			x--
		case Neighbor4Right:
			x++
		case Neighbor4Top:
			y++
		}

		xyb := block.XYB{
			XY:    block.XY{X: x, Y: y},
			Block: f.GetXY(x, y),
		}

		fn(xyb)
	}
}

type pathNode struct {
	xy        block.XY
	length    int
	heuristic float32
	origin    *pathNode
}

type pathQueue []*pathNode

func (pq pathQueue) Len() int      { return len(pq) }
func (pq pathQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq pathQueue) Less(i, j int) bool {
	return float32(pq[i].length)+pq[i].heuristic < float32(pq[j].length)+pq[j].heuristic
}

func (pq *pathQueue) Push(n any) {
	*pq = append(*pq, n.(*pathNode))
}

func (pq *pathQueue) Pop() any {
	count := len(*pq)
	n := (*pq)[count-1]
	*pq = (*pq)[:count-1]
	return n
}

// Path4 finds the shortest path between the start and the goal. It doesn't contain diagonal movement.
// If it's found, the resulting slice will contain all coordinates of the path, including the start and the goal.
// So at least two elements must be in the list. However, if no path could be found it returns nil.
func (f *Field) Path4(start, goal block.XY) []block.XY {
	return f.path(start, goal, false)
}

// Path8 finds the shortest path between the start and the goal. It contains diagonal movement.
// If it's found, the resulting slice will contain all coordinates of the path, including the start and the goal.
// So at least two elements must be in the list. However, if no path could be found it returns nil.
func (f *Field) Path8(start, goal block.XY) []block.XY {
	return f.path(start, goal, true)
}

func (f *Field) path(start, goal block.XY, diagonalMove bool) []block.XY {
	openSet := &pathQueue{}
	heap.Init(openSet)

	var heuristic func(block.XY, block.XY) float32

	if diagonalMove {
		heuristic = func(p0, p1 block.XY) float32 {
			dx := p0.X - p1.X
			dy := p0.Y - p1.Y
			dd := dx*dx + dy*dy
			return float32(math.Sqrt(float64(dd)))
		}
	} else {
		heuristic = func(p0, p1 block.XY) float32 {
			return float32(geommath.Manhattan[int](geommath.P[int](p0), geommath.P[int](p1)))
		}
	}

	heap.Push(openSet, &pathNode{
		xy:        start,
		length:    1,
		heuristic: heuristic(start, goal),
		origin:    nil,
	})

	visited := make(map[block.XY]struct{})

	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*pathNode)

		if current.xy == goal {
			path := make([]block.XY, 0, current.length)
			for n := current; n != nil; n = n.origin {
				path = append(path, n.xy)
			}
			slices.Reverse(path)

			return path
		}

		visited[current.xy] = struct{}{}

		var neighbors interface {
			ForEach(*Field, block.XY, func(xyb block.XYB))
		}

		if diagonalMove {
			neighbors = f.Neighbors8(current.xy)
		} else {
			neighbors = f.Neighbors4(current.xy)
		}

		neighbors.ForEach(f, current.xy, func(neighbor block.XYB) {
			if _, yes := visited[neighbor.XY]; yes {
				return
			}

			heap.Push(openSet, &pathNode{
				xy:        neighbor.XY,
				length:    current.length + 1,
				heuristic: heuristic(neighbor.XY, goal),
				origin:    current,
			})
		})
	}

	return nil
}

// FindNearest8 runs the provided function on every block on coordinates around the starting position,
// starting from the bounding rectangle around the position, and then it's bounding rectangle, and so on
// up to the provided range, until the provided function returns true.
func (f *Field) FindNearest8(pos block.XY, r int, fn func(block.XYB, int) bool) (block.XYB, bool) {
	for d := 1; d <= r; d++ {
		x0 := pos.X - d
		x1 := pos.X + d
		y0 := pos.Y - d
		y1 := pos.Y + d

		xi0 := max(x0, 0)
		xi1 := min(x1, f.w-1)
		yj0 := max(y0+1, 0)
		yj1 := min(y1-1, f.h-1)

		if y0 >= 0 {
			for i := xi0; i <= xi1; i++ {
				xy := block.XY{X: i, Y: y0}
				b := f.GetXY(xy.X, xy.Y)
				xyb := block.XYB{XY: xy, Block: b}
				if fn(xyb, d) {
					return xyb, true
				}
			}
		}

		for j := yj0; j <= yj1; j++ {
			if x0 >= 0 {
				xy := block.XY{X: x0, Y: j}
				b := f.GetXY(xy.X, xy.Y)
				xyb := block.XYB{XY: xy, Block: b}
				if fn(xyb, d) {
					return xyb, true
				}
			}
			if x1 < f.w {
				xy := block.XY{X: x1, Y: j}
				b := f.GetXY(xy.X, xy.Y)
				xyb := block.XYB{XY: xy, Block: b}
				if fn(xyb, d) {
					return xyb, true
				}
			}
		}

		if y1 < f.h {
			for i := xi0; i <= xi1; i++ {
				xy := block.XY{X: i, Y: y1}
				b := f.GetXY(xy.X, xy.Y)
				xyb := block.XYB{XY: xy, Block: b}
				if fn(xyb, d) {
					return xyb, true
				}
			}
		}
	}

	return block.XYB{XY: pos, Block: block.Block{Type: block.TypeEmpty}}, false
}

// FindNearest4 runs the provided function on every block on coordinates around the starting position,
// starting from the bounding rhombus around the position, and then it's bounding rhombus, and so on
// up to the provided range, until the provided function returns true.
func (f *Field) FindNearest4(pos block.XY, r int, fn func(block.XYB, int) bool) (block.XYB, bool) {
	w := f.w
	h := f.h
	fn1 := func(x, y, r int) (block.XYB, bool) {
		if x < 0 || x >= w || y < 0 || y >= h {
			return block.XYB{}, false
		}

		xy := block.XY{X: x, Y: y}
		b := f.GetXY(xy.X, xy.Y)
		xyb := block.XYB{XY: xy, Block: b}
		return xyb, fn(xyb, r)
	}

	for d := 1; d <= r; d++ {
		yi := pos.Y - d
		if xyb, stop := fn1(pos.X, yi, d); stop {
			return xyb, true
		}
		yi++
		for i := 1; i < d; i++ {
			if xyb, stop := fn1(pos.X-i, yi, d); stop {
				return xyb, true
			}
			if xyb, stop := fn1(pos.X+i, yi, d); stop {
				return xyb, true
			}
			yi++
		}
		fn1(pos.X-d, yi, d)
		fn1(pos.X+d, yi, d)
		yi++
		for i := d - 1; i > 0; i-- {
			if xyb, stop := fn1(pos.X-i, yi, d); stop {
				return xyb, true
			}
			if xyb, stop := fn1(pos.X+i, yi, d); stop {
				return xyb, true
			}
			yi++
		}
		if xyb, stop := fn1(pos.X, yi, d); stop {
			return xyb, true
		}
	}

	return block.XYB{XY: pos, Block: block.Block{Type: block.TypeEmpty}}, false
}
