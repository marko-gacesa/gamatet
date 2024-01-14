// Copyright (c) 2020-2024 by Marko Gaćeša

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
	const total = 12

	bagIdx := idx / total
	idx = idx % total

	r := random{uint32(f.seed + 857*bagIdx + 13), uint32(f.seed + 328*bagIdx + 17)}

	var m [total]int
	r.perm(m[:])

	switch m[idx] {
	default:
		fallthrough
	case 0:
		return Shooter(5, block.TypeLava)
	case 1:
		return Shooter(5, block.TypeAcid)
	case 2:
		return Shooter(5, block.TypeWave)
	case 3:
		return NewPentomino(0, block.Rock)
	case 4:
		return NewPentomino(0, block.Wave)
	case 5:
		b := block.Rock
		b.Hardness = 1
		return NewTetromino(TetrominoI, b)
	case 6:
		b := block.Rock
		b.Hardness = 2
		return NewTetromino(TetrominoO, b)
	case 7:
		return NewTetromino(TetrominoO, block.Acid)
	case 8:
		return NewTetromino(TetrominoO, block.Lava)
	case 9:
		b := block.Rock
		b.Hardness = 3
		return NewPentomino(13, b)
	case 10:
		return NewPentomino(13, block.Acid)
	case 11:
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
