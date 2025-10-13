// Copyright (c) 2020-2025 by Marko Gaćeša

package piece

import (
	"gamatet/game/block"
	"strconv"
	"sync"
)

type Feed interface {
	Get(idx int) Piece
}

const MaxBagSize = 4

type GenericFeed struct {
	seed          int
	pieceBagCount int
	shapeCount    int
	fn            func(idx int) Piece
	pool          *sync.Pool
}

func NewGenericFeed(bagSize int, seed int, shapeCount int, fn func(idx int) Piece) GenericFeed {
	if bagSize < 1 || bagSize > MaxBagSize {
		panic("bagSize must be from 1 to " + strconv.Itoa(MaxBagSize))
	}

	pieceBagCount := shapeCount * bagSize

	return GenericFeed{
		seed:          seed,
		pieceBagCount: pieceBagCount,
		shapeCount:    shapeCount,
		fn:            fn,
		pool: &sync.Pool{
			New: func() interface{} {
				return make([]int, pieceBagCount)
			},
		},
	}
}

func (f GenericFeed) shapeIdx(idx int) int {
	bagIdx := idx / f.pieceBagCount
	idx = idx % f.pieceBagCount

	r := random{uint32(f.seed + 857*bagIdx + 13), uint32(f.seed + 328*bagIdx + 17)}

	m := f.pool.Get().([]int)
	r.perm(m)
	shapeIdx := m[idx] % f.shapeCount
	f.pool.Put(m)

	return shapeIdx
}

func (f GenericFeed) Get(idx int) Piece {
	return f.fn(f.shapeIdx(idx))
}

func NewRotTetrominoFeed(bagSize int, seed int) GenericFeed {
	return NewGenericFeed(bagSize, seed, len(shapesRotTetrominoes), func(idx int) Piece {
		return &polyominoRot{
			shapeSquare: shapesRotTetrominoes[idx],
			rot:         0,
			block: block.Block{
				Type:     block.TypeRock,
				Hardness: 0,
				Color:    colors[idx%len(colors)],
			},
		}
	})
}

func NewFlipVTetrominoFeed(bagSize int, seed int) GenericFeed {
	return NewGenericFeed(bagSize, seed, len(shapesFlipVTetrominoes), func(idx int) Piece {
		return &polyominoFlip{
			shapeRect: shapesFlipVTetrominoes[idx],
			block: block.Block{
				Type:     block.TypeRock,
				Hardness: 0,
				Color:    colors[idx%len(colors)],
			},
		}
	})
}
