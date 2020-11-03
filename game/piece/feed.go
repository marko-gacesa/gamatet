// Copyright (c) 2020 by Marko Gaćeša

package piece

import (
	"gamatet/game/block"
	"strconv"
)

type Feed interface {
	Get(idx int) Piece
}

const MaxBagSize = 4

type tetrominoFeed struct {
	bagSize int
	seed    int
}

func NewTetrominoFeed(bagSize int, seed int) Feed {
	if bagSize < 1 || bagSize > MaxBagSize {
		panic("bagSize must be from 1 to " + strconv.Itoa(MaxBagSize))
	}
	return &tetrominoFeed{
		bagSize: ShapeCountTetrominoes * bagSize,
		seed:    seed,
	}
}

func (f *tetrominoFeed) Get(idx int) Piece {
	bagIdx := idx / f.bagSize
	idx = idx % f.bagSize

	r := random{uint32(f.seed + 857*bagIdx + 13), uint32(f.seed + 328*bagIdx + 17)}

	var m [MaxBagSize * ShapeCountTetrominoes]int
	r.perm(m[:f.bagSize])

	return NewStandardTetromino(m[idx] % ShapeCountTetrominoes)
}

type pentaFeed struct {
	bagSize int
	seed    int
}

func NewPentaFeed(bagSize, seed int) Feed {
	if bagSize < 1 || bagSize > MaxBagSize {
		panic("bagSize must be from 1 to " + strconv.Itoa(MaxBagSize))
	}
	return &pentaFeed{
		bagSize: ShapeCountAll * bagSize,
		seed:    seed,
	}
}

func (f *pentaFeed) Get(idx int) Piece {
	bagIdx := idx / f.bagSize
	idx = idx % f.bagSize

	r := random{uint32(f.seed + 857*bagIdx + 13), uint32(f.seed + 328*bagIdx + 17)}

	var m [MaxBagSize * ShapeCountAll]int
	r.perm(m[:f.bagSize])

	return NewAnyPolyomino(m[idx] % ShapeCountAll)
}

type debug struct {
	seed int
}

func NewDebugFeed(seed int) Feed {
	return &debug{seed: seed}
}

func (f *debug) Get(idx int) Piece {
	bagIdx := idx / 10
	idx = idx % 10

	r := random{uint32(f.seed + 857*bagIdx + 13), uint32(f.seed + 328*bagIdx + 17)}

	var m [10]int
	r.perm(m[:])

	switch m[idx] {
	default:
		fallthrough
	case 0:
		return Shooter(5, block.TypeLava)
	case 1:
		return Shooter(5, block.TypeAcid)
	case 2:
		return NewPentomino(0, block.Rock)
	case 3:
		return NewTetromino(TetrominoI, block.Rock)
	case 4:
		return NewTetromino(TetrominoO, block.Rock)
	case 5:
		return NewTetromino(TetrominoO, block.Acid)
	case 6:
		return NewTetromino(TetrominoO, block.Lava)
	case 7:
		return NewPentomino(13, block.Rock)
	case 8:
		return NewPentomino(13, block.Acid)
	case 9:
		return NewPentomino(13, block.Lava)
	}
}

type battle struct {
	bagSize int
	seed    int
}

func NewBattleFeed(bagSize, seed int) Feed {
	return &battle{
		bagSize: ShapeCountTetrominoes*bagSize + 2,
		seed:    seed,
	}
}

func (f *battle) Get(idx int) Piece {
	bagIdx := idx / f.bagSize
	idx = idx % f.bagSize

	r := random{uint32(f.seed + 857*bagIdx + 13), uint32(f.seed + 328*bagIdx + 17)}

	// TODO: Unfinished. Improve the implementation.

	var m [MaxBagSize*ShapeCountTetrominoes + 2]int
	r.perm(m[:f.bagSize])

	switch m[idx] {
	case 0:
		return Shooter(5, block.TypeLava)
	case 1:
		return Shooter(5, block.TypeAcid)
	default:
		shape := (m[idx] - 2) % len(tetrominoes)
		r = random{uint32(idx), uint32(idx + 17)}
		t := r.int(3)
		switch t {
		case 0:
			return NewTetromino(shape, block.Lava)
		case 1:
			return NewTetromino(shape, block.Acid)
		default:
			return NewStandardTetromino(shape)
		}
	}
}
