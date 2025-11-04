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

const (
	MaxBagSize = 6

	SizeMin     = 3
	SizeMax     = 5
	SizeDefault = 4
	Size3       = 3
	Size4       = 4
	Size5       = 5
)

// SamePieceFeed is the feed that always return the same piece. Useful for testing.
type SamePieceFeed struct{ Piece }

func (p SamePieceFeed) Get(int) Piece { return p.Piece }

type QFeed struct{}

func (p QFeed) Get(int) Piece { return &polyominoFlipH{shapeRect: shapeQ, block: block.Rock} }

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
			New: func() any {
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

var shapesRot = map[byte][]shapeSquare{
	Size3: shapesRotTinyminoes,
	Size4: shapesRotTetrominoes,
	Size5: slices.Concat(shapesRotTinyminoes, shapesRotTinyminoes, shapesRotTetrominoes, shapesRotTetrominoes, shapesRotPentominoes),
}

func NewRotTetrominoFeed(pieceSize byte, bagSize int, seed int, c Color) Feed {
	if pieceSize < SizeMin || pieceSize > SizeMax {
		return QFeed{}
	}
	shapes := shapesRot[pieceSize]
	return NewGenericFeed(bagSize, seed, len(shapes), func(idx int) Piece {
		return &polyominoRot{
			shapeSquare: shapes[idx],
			rot:         0,
			block: block.Block{
				Type:     block.TypeRock,
				Hardness: 0,
				Color:    c.Color(idx),
			},
		}
	})
}

var shapesFlipV = map[byte][]shapeRect{
	Size3: shapesFlipVTinyminoes,
	Size4: slices.Concat(shapesFlipVTinyminoes, shapesFlipVTetrominoes),
	Size5: slices.Concat(shapesFlipVTinyminoes, shapesFlipVTinyminoes, shapesFlipVTetrominoes, shapesFlipVTetrominoes, shapesFlipHPentominoes),
}

func NewFlipVFeed(pieceSize byte, bagSize int, seed int, c Color) Feed {
	if pieceSize < SizeMin || pieceSize > SizeMax {
		return QFeed{}
	}
	shapes := shapesFlipV[pieceSize]
	return NewGenericFeed(bagSize, seed, len(shapes), func(idx int) Piece {
		return &polyominoFlipV{
			shapeRect: shapes[idx],
			block: block.Block{
				Type:     block.TypeRock,
				Hardness: 0,
				Color:    c.Color(idx),
			},
		}
	})
}

var shapesFlipH = map[byte][]shapeRect{
	Size3: shapesFlipHTinyminoes,
	Size4: slices.Concat(shapesFlipHTinyminoes, shapesFlipHTetrominoes),
	Size5: slices.Concat(shapesFlipHTinyminoes, shapesFlipHTinyminoes, shapesFlipHTetrominoes, shapesFlipHTetrominoes, shapesFlipHPentominoes),
}

func NewFlipHFeed(pieceSize byte, bagSize int, seed int, c Color) Feed {
	if pieceSize < SizeMin || pieceSize > SizeMax {
		return QFeed{}
	}
	shapes := shapesFlipH[pieceSize]
	return NewGenericFeed(bagSize, seed, len(shapes), func(idx int) Piece {
		return &polyominoFlipH{
			shapeRect: shapes[idx],
			block: block.Block{
				Type:     block.TypeRock,
				Hardness: 0,
				Color:    c.Color(idx),
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
					Color:    DefaultColor{}.Color(idx),
				},
			}
		case 1:
			return &polyominoFlipV{
				shapeRect: shapesFlipVTetrominoes[3],
				block: block.Block{
					Type:     block.TypeRock,
					Hardness: 2,
					Color:    DefaultColor{}.Color(idx),
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
