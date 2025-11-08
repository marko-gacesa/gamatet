// Copyright (c) 2020-2025 by Marko Gaćeša

package piece

import (
	"gamatet/game/block"
	"slices"
	"strconv"
	"sync"
)

type Feed interface {
	Get(idx, playerIdx int) Piece
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

func (p SamePieceFeed) Get(int, int) Piece { return p.Piece }

type QFeed struct{}

func (p QFeed) Get(int, int) Piece { return &polyominoFlipH{shapeRect: shapeQ, block: block.Rock} }

type GenericFeed struct {
	seed          int
	pieceBagCount int
	shapeCount    int
	fn            func(idx, playerIdx int) Piece
	pool          *sync.Pool
}

func NewGenericFeed(bagSize int, seed int, shapeCount int, fn func(idx, playerIdx int) Piece) GenericFeed {
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

func (f GenericFeed) Get(idx, playerIdx int) Piece {
	bagIdx := idx / f.pieceBagCount
	idx = idx % f.pieceBagCount

	r := random{uint32(f.seed + 857*bagIdx + 13), uint32(f.seed + 328*bagIdx + 17)}

	m := f.pool.Get().([]int)
	r.perm(m)
	shapeIdx := m[idx] % f.shapeCount
	f.pool.Put(m)

	return f.fn(shapeIdx, playerIdx)
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
	return NewGenericFeed(bagSize, seed, len(shapes), func(idx, playerIdx int) Piece {
		return &polyominoRot{
			shapeSquare: shapes[idx],
			rot:         0,
			block: block.Block{
				Type:     block.TypeRock,
				Hardness: 0,
				Color:    c.Color(idx, playerIdx),
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
	return NewGenericFeed(bagSize, seed, len(shapes), func(idx, playerIndex int) Piece {
		return &polyominoFlipV{
			shapeRect: shapes[idx],
			block: block.Block{
				Type:     block.TypeRock,
				Hardness: 0,
				Color:    c.Color(idx, playerIndex),
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
	return NewGenericFeed(bagSize, seed, len(shapes), func(idx, playerIndex int) Piece {
		return &polyominoFlipH{
			shapeRect: shapes[idx],
			block: block.Block{
				Type:     block.TypeRock,
				Hardness: 0,
				Color:    c.Color(idx, playerIndex),
			},
		}
	})
}

type CtrlFeed struct {
	internal Feed
	fIdx     int
	ctrlIdx  int
	same     bool
}

func NewCtrlFeed(internal Feed, fIdx, ctrlIdx int, same bool) *CtrlFeed {
	return &CtrlFeed{
		internal: internal,
		fIdx:     fIdx,
		ctrlIdx:  ctrlIdx,
		same:     same,
	}
}

func (f *CtrlFeed) Get(idx, playerIndex int) Piece {
	if f.same {
		return f.internal.Get(idx, playerIndex)
	}
	return f.internal.Get(idx+f.fIdx*137+f.ctrlIdx*5, playerIndex)
}
