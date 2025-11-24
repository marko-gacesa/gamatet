// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package piece

import "github.com/marko-gacesa/gamatet/game/block"

// See: https://en.wikipedia.org/wiki/Polyomino
const (
	TetrominoO = iota
	TetrominoI
	TetrominoT
	TetrominoJ
	TetrominoL
	TetrominoS
	TetrominoZ
)

func NewTetromino(id int, b block.Block) Piece {
	return &polyominoRot{
		shapeSquare: shapesRotTetrominoes[id],
		rot:         0,
		block:       b,
	}
}

func NewColorTetromino(id int, color uint32) Piece {
	return NewTetromino(id, block.Block{
		Type:     block.TypeRock,
		Hardness: 0,
		Color:    color,
	})
}

func NewStandardTetromino(id int) Piece {
	return NewColorTetromino(id, DefaultColor{}.Color(id, 0))
}

func NewFlipVTetromino(id int, b block.Block) Piece {
	return &polyominoFlipV{
		shapeRect: shapesFlipVTetrominoes[id],
		block:     b,
	}
}

func NewPentominos(id int, b block.Block) Piece {
	return &polyominoRot{
		shapeSquare: shapesRotPentominoes[id],
		rot:         0,
		block:       b,
	}
}

func NewO(b block.Block) Piece {
	return &polyominoDumb{
		shapeRect: shapeO,
		block:     b,
	}
}
