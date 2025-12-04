// Copyright (c) 2020-2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package piece

import (
	"slices"
	"strconv"
	"sync"

	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/logic/random"
)

type Feed interface {
	Get(idx uint, playerIdx byte) Piece
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

func (p SamePieceFeed) Get(uint, byte) Piece { return p.Piece }

type QFeed struct{}

func (p QFeed) Get(uint, byte) Piece { return NewQ(block.Rock) }

type GenericFeed struct {
	seed          uint
	pieceBagCount uint
	shapeCount    uint
	fn            func(idx uint, playerIdx byte) Piece
	pool          *sync.Pool
}

func NewGenericFeed(bagSize int, seed int, shapeCount int, fn func(idx uint, playerIdx byte) Piece) GenericFeed {
	if bagSize < 1 || bagSize > MaxBagSize {
		panic("bagSize must be from 1 to " + strconv.Itoa(MaxBagSize))
	}
	if shapeCount < 1 {
		panic("shapeCount must a positive integer")
	}

	pieceBagCount := shapeCount * bagSize

	return GenericFeed{
		seed:          uint(seed),
		pieceBagCount: uint(pieceBagCount),
		shapeCount:    uint(shapeCount),
		fn:            fn,
		pool: &sync.Pool{
			New: func() any {
				return make([]uint, pieceBagCount)
			},
		},
	}
}

func (f GenericFeed) Get(idx uint, playerIdx byte) Piece {
	bagIdx := idx / f.pieceBagCount
	offs := idx % f.pieceBagCount

	r := random.New(f.seed+857*bagIdx+13, f.seed+328*bagIdx+17)

	m := f.pool.Get().([]uint)
	r.Perm(m)
	shapeIdx := m[offs] % f.shapeCount
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
	return NewGenericFeed(bagSize, seed, len(shapes), func(idx uint, playerIdx byte) Piece {
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

var shapesFlipV = map[byte][]shapeRectV{
	Size3: shapesFlipVTinyminoes,
	Size4: slices.Concat(shapesFlipVTinyminoes, shapesFlipVTetrominoes),
	Size5: slices.Concat(shapesFlipVTinyminoes, shapesFlipVTinyminoes, shapesFlipVTetrominoes, shapesFlipVTetrominoes, shapesFlipVPentominoes),
}

func NewFlipVFeed(pieceSize byte, bagSize int, seed int, c Color) Feed {
	if pieceSize < SizeMin || pieceSize > SizeMax {
		return QFeed{}
	}
	shapes := shapesFlipV[pieceSize]
	return NewGenericFeed(bagSize, seed, len(shapes), func(idx uint, playerIndex byte) Piece {
		return &polyominoFlipV{
			shapeRect: shapeRect(shapes[idx]),
			block: block.Block{
				Type:     block.TypeRock,
				Hardness: 0,
				Color:    c.Color(idx, playerIndex),
			},
		}
	})
}

var shapesFlipH = map[byte][]shapeRectH{
	Size3: shapesFlipHTinyminoes,
	Size4: slices.Concat(shapesFlipHTinyminoes, shapesFlipHTetrominoes),
	Size5: slices.Concat(shapesFlipHTinyminoes, shapesFlipHTinyminoes, shapesFlipHTetrominoes, shapesFlipHTetrominoes, shapesFlipHPentominoes),
}

func NewFlipHFeed(pieceSize byte, bagSize int, seed int, c Color) Feed {
	if pieceSize < SizeMin || pieceSize > SizeMax {
		return QFeed{}
	}
	shapes := shapesFlipH[pieceSize]
	return NewGenericFeed(bagSize, seed, len(shapes), func(idx uint, playerIndex byte) Piece {
		return &polyominoFlipH{
			shapeRect: shapeRect(shapes[idx]),
			block: block.Block{
				Type:     block.TypeRock,
				Hardness: 0,
				Color:    c.Color(idx, playerIndex),
			},
		}
	})
}

func NewMixedFeed(bagSize int, seed int, c Color, shapes ...any) Feed {
	return NewGenericFeed(bagSize, seed, len(shapes), func(idx uint, playerIndex byte) Piece {
		switch s := shapes[idx].(type) {
		case shapeSquare:
			return &polyominoRot{
				shapeSquare: s,
				block: block.Block{
					Type:     block.TypeRock,
					Hardness: 0,
					Color:    c.Color(idx, playerIndex),
				},
			}
		case shapeRectV:
			return &polyominoFlipV{
				shapeRect: shapeRect(s),
				block: block.Block{
					Type:     block.TypeRock,
					Hardness: 0,
					Color:    c.Color(idx, playerIndex),
				},
			}
		case shapeRectH:
			return &polyominoFlipH{
				shapeRect: shapeRect(s),
				block: block.Block{
					Type:     block.TypeRock,
					Hardness: 0,
					Color:    c.Color(idx, playerIndex),
				},
			}
		default:
			return NewQ(block.Rock)
		}
	})
}

type CtrlFeed struct {
	internal Feed
	fIdx     int
	ctrlIdx  int
	same     bool
	override map[uint]Piece
}

func NewCtrlFeed(internal Feed, fIdx, ctrlIdx int, same bool) *CtrlFeed {
	return &CtrlFeed{
		internal: internal,
		fIdx:     fIdx,
		ctrlIdx:  ctrlIdx,
		same:     same,
		override: map[uint]Piece{},
	}
}

func (f *CtrlFeed) Get(idx uint, playerIndex byte) Piece {
	if override, ok := f.override[idx]; ok {
		return override
	}
	if f.same {
		return f.internal.Get(idx, playerIndex)
	}
	return f.internal.Get(idx+uint(f.fIdx)*137+uint(f.ctrlIdx)*5, playerIndex)
}

func (f *CtrlFeed) Overridden(idx uint) bool {
	_, ok := f.override[idx]
	return ok
}

func (f *CtrlFeed) Override(idx uint, piece Piece) {
	if _, ok := f.override[idx]; ok {
		return
	}
	f.override[idx] = piece
}

func (f *CtrlFeed) OverrideClear(idx uint) {
	delete(f.override, idx)
}
