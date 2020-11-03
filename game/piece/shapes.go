// Copyright (c) 2020 by Marko Gaćeša

package piece

import (
	"gamatet/game/block"
)

const (
	TetrominoO = iota
	TetrominoI
	TetrominoT
	TetrominoJ
	TetrominoL
	TetrominoS
	TetrominoZ
)

var TinyminoColors = []uint32{
	0x7FFF0000, // lime
	0x007FFF00, // azure
	0x7F00FF00, // violet
	0x80008000, // purple
}

var TetrominoColors = []uint32{
	0xFFFF0000, // yellow
	0x00FFFF00, // cyan
	0xFF00FF00, // magenta
	0x0000FF00, // blue
	0xFF7F0000, // orange
	0x00FF0000, // green
	0xFF000000, // red
}

var PentominoColors = []uint32{
	0xFFD70000, // gold
	0x964B0000, // brown
	0xFFC0CB00, // pink
	0x3F00FF00, // indigo
	0xFF007F00, // pink-red
	0x00FF7F00, // spring-green
	0xFF007F00, // electric-magenta
	0xC3B09100, // khaki
	0xFFE5B400, // peach
	0x00808000, // teal
	0x80800000, // olive
	0xF5F5DC00, // beige
	0x40E0D000, // turquoise
}

func NewStandardTetromino(shape int) Piece {
	return &polyomino{polyominoShape: tetrominoes[shape], block: block.Block{Type: block.TypeRock, Color: TetrominoColors[shape]}}
}

func NewColorTetromino(shape int, color uint32) Piece {
	return &polyomino{polyominoShape: tetrominoes[shape], block: block.Block{Type: block.TypeRock, Color: color}}
}

func NewTetromino(shape int, b block.Block) Piece {
	return &polyomino{polyominoShape: tetrominoes[shape], block: b}
}

func NewPentomino(shape int, b block.Block) Piece {
	return &polyomino{polyominoShape: pentominoes[shape], block: b}
}

const (
	ShapeCountMonominoes  = 1
	ShapeCountDominoes    = 1
	ShapeCountTrominoes   = 2
	ShapeCountTetrominoes = 7
	ShapeCountPentominoes = 18

	ShapeCountAll = ShapeCountMonominoes + ShapeCountDominoes + ShapeCountTrominoes + ShapeCountTetrominoes + ShapeCountPentominoes
)

func NewAnyPolyomino(shape int) Piece {
	if shape < 7 {
		return &polyomino{polyominoShape: tetrominoes[shape], block: block.Block{Type: block.TypeRock, Color: TetrominoColors[shape]}}
	} else if shape == 7 {
		return &polyomino{polyominoShape: monomino, block: block.Block{Type: block.TypeRock, Color: TinyminoColors[0]}}
	} else if shape == 8 {
		return &polyomino{polyominoShape: domino, block: block.Block{Type: block.TypeRock, Color: TinyminoColors[1]}}
	} else if shape == 9 || shape == 10 {
		return &polyomino{polyominoShape: trominoes[shape-9], block: block.Block{Type: block.TypeRock, Color: TinyminoColors[shape-7]}}
	} else {
		return &polyomino{polyominoShape: pentominoes[shape-11], block: block.Block{Type: block.TypeRock, Color: PentominoColors[shape%len(PentominoColors)]}}
	}
}

// See: https://en.wikipedia.org/wiki/Polyomino

var (
	monomino = polyominoShape{dim: 1, size: 1, rots: 0, data: 1}
	domino   = polyominoShape{dim: 2, size: 2, rots: 2, data: 3}

	trominoes = []polyominoShape{
		0: {dim: 3, size: 3, rots: 2, data: 56},
		1: {dim: 2, size: 3, rots: 4, data: 14},
	}

	tetrominoes = []polyominoShape{
		0: {dim: 2, size: 4, rots: 0, data: 15},
		1: {dim: 4, size: 4, rots: 2, data: 240},
		2: {dim: 3, size: 4, rots: 4, data: 58},
		3: {dim: 3, size: 4, rots: 4, data: 57},
		4: {dim: 3, size: 4, rots: 4, data: 60},
		5: {dim: 3, size: 4, rots: 2, data: 30},
		6: {dim: 3, size: 4, rots: 2, data: 51},
	}

	pentominoes = []polyominoShape{
		0:  {dim: 5, size: 5, rots: 2, data: 31744},
		1:  {dim: 3, size: 5, rots: 4, data: 185},
		2:  {dim: 3, size: 5, rots: 4, data: 188},
		3:  {dim: 4, size: 5, rots: 4, data: 3856},
		4:  {dim: 4, size: 5, rots: 4, data: 3968},
		5:  {dim: 3, size: 5, rots: 4, data: 62},
		6:  {dim: 3, size: 5, rots: 4, data: 59},
		7:  {dim: 4, size: 5, rots: 4, data: 3632},
		8:  {dim: 4, size: 5, rots: 4, data: 1984},
		9:  {dim: 3, size: 5, rots: 4, data: 466},
		10: {dim: 3, size: 5, rots: 4, data: 61},
		11: {dim: 3, size: 5, rots: 4, data: 484},
		12: {dim: 3, size: 5, rots: 4, data: 244},
		13: {dim: 3, size: 5, rots: 0, data: 186},
		14: {dim: 4, size: 5, rots: 4, data: 3904},
		15: {dim: 4, size: 5, rots: 4, data: 3872},
		16: {dim: 3, size: 5, rots: 4, data: 214},
		17: {dim: 3, size: 5, rots: 4, data: 403},
	}
)
