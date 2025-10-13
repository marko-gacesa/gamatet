// Copyright (c) 2025 by Marko Gaćeša

package piece

import "gamatet/game/block"

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
	return NewColorTetromino(id, colors[id])
}
