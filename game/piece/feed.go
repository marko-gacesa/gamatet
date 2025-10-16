// Copyright (c) 2020-2025 by Marko Gaćeša

package piece

import (
	"gamatet/game/block"
	"slices"
	"strconv"
	"sync"
)

type Feed interface {
	Get(idx int) Piece
}

const MaxBagSize = 4

// SamePieceFeed is the feed that always return the same piece. Useful for testing.
type SamePieceFeed struct{ Piece }

func (p SamePieceFeed) Get(int) Piece { return p.Piece }

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

func (f GenericFeed) Get(idx int) Piece {
	bagIdx := idx / f.pieceBagCount
	idx = idx % f.pieceBagCount

	r := random{uint32(f.seed + 857*bagIdx + 13), uint32(f.seed + 328*bagIdx + 17)}

	m := f.pool.Get().([]int)
	r.perm(m)
	shapeIdx := m[idx] % f.shapeCount
	f.pool.Put(m)

	return f.fn(shapeIdx)
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

func NewFlipVFeed(bagSize int, seed int) GenericFeed {
	shapes := slices.Concat(shapesFlipVTinyminoes, shapesFlipVTetrominoes)
	return NewGenericFeed(bagSize, seed, len(shapes), func(idx int) Piece {
		return &polyominoFlipV{
			shapeRect: shapes[idx],
			block: block.Block{
				Type:     block.TypeRock,
				Hardness: 0,
				Color:    colors[idx%len(colors)],
			},
		}
	})
}

func NewFlipHFeed(bagSize int, seed int) GenericFeed {
	shapes := slices.Concat(shapesFlipHTinyminoes, shapesFlipHTetrominoes)
	return NewGenericFeed(bagSize, seed, len(shapes), func(idx int) Piece {
		return &polyominoFlipH{
			shapeRect: shapes[idx],
			block: block.Block{
				Type:     block.TypeRock,
				Hardness: 0,
				Color:    colors[idx%len(colors)],
			},
		}
	})
}

type TestFeed struct{}

func NewTestFeed(bagSize int, seed int) GenericFeed {

	return NewGenericFeed(bagSize, seed, 4, func(idx int) Piece {
		switch idx {
		default:
			fallthrough
		case 0:
			return &polyominoFlipH{
				shapeRect: shapesFlipHTetrominoes[5],
				block: block.Block{
					Type:     block.TypeRock,
					Hardness: 1,
					Color:    colors[idx%len(colors)],
				},
			}
		case 1:
			return &polyominoFlipV{
				shapeRect: shapesFlipVTetrominoes[3],
				block: block.Block{
					Type:     block.TypeRock,
					Hardness: 2,
					Color:    colors[idx%len(colors)],
				},
			}
		case 2:
			return &polyominoRot{
				shapeSquare: shapesRotTetrominoes[TetrominoJ],
				block: block.Block{
					Type:     block.TypeLava,
					Hardness: 0,
					Color:    block.Lava.Color,
				},
			}
		case 3:
			return &shooter{
				bulletType: block.TypeAcid,
				ammo:       5,
			}
		}
	})
}
