// Copyright (c) 2020, 2025 by Marko Gaćeša

package piece

import (
	"slices"
	"testing"
)

// when defining shape matrix for a piece for shapes pull the shape up:
// like this or this... ... and *NOT* like this:
// . . . .    . # #         . . . .   . . .  <- WRONG: empty row, the shape should be pulled up by 1
// # # # #    # # .         . . . .   . # #
// . . . .    . . .         # # # #   # # .
// . . . .                  . . . .
//
// also, define it in the horizontal position:
// like this or this... ... and *NOT* like this:
// . . . .    # . .         . # . .   . # .
// # # # #    # # #         . # . .   . # .
// . . . .    . . .         . # . .   # # .
// . . . .                  . # . .

func TestShapesRot(t *testing.T) {
	_tinyminoes := make([]shapeSquare, 0, 4) // monomino + dominoes + trominoes

	_tinyminoes = append(_tinyminoes, _initShapeSquare(0, []bool{
		XX,
	}))

	_tinyminoes = append(_tinyminoes, _initShapeSquare(2, []bool{
		XX, XX,
		__, __,
	}))

	_tinyminoes = append(_tinyminoes, _initShapeSquare(2, []bool{
		__, __, __,
		XX, XX, XX,
		__, __, __,
	}))

	_tinyminoes = append(_tinyminoes, _initShapeSquare(4, []bool{
		__, XX,
		XX, XX,
	}))

	_tetrominoes := make([]shapeSquare, 0, 7)

	_tetrominoes = append(_tetrominoes, _initShapeSquare(0, []bool{
		XX, XX,
		XX, XX,
	}))

	_tetrominoes = append(_tetrominoes, _initShapeSquare(2, []bool{
		__, __, __, __,
		XX, XX, XX, XX,
		__, __, __, __,
		__, __, __, __,
	}))

	_tetrominoes = append(_tetrominoes, _initShapeSquare(4, []bool{
		__, XX, __,
		XX, XX, XX,
		__, __, __,
	}))

	_tetrominoes = append(_tetrominoes, _initShapeSquare(4, []bool{
		XX, __, __,
		XX, XX, XX,
		__, __, __,
	}))

	_tetrominoes = append(_tetrominoes, _initShapeSquare(4, []bool{
		__, __, XX,
		XX, XX, XX,
		__, __, __,
	}))

	_tetrominoes = append(_tetrominoes, _initShapeSquare(2, []bool{
		__, XX, XX,
		XX, XX, __,
		__, __, __,
	}))

	_tetrominoes = append(_tetrominoes, _initShapeSquare(2, []bool{
		XX, XX, __,
		__, XX, XX,
		__, __, __,
	}))

	_pentominoes := make([]shapeSquare, 0, 18)

	_pentominoes = append(_pentominoes, _initShapeSquare(2, []bool{
		__, __, __, __, __,
		__, __, __, __, __,
		XX, XX, XX, XX, XX,
		__, __, __, __, __,
		__, __, __, __, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeSquare(4, []bool{
		XX, __, __,
		XX, XX, XX,
		__, XX, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeSquare(4, []bool{
		__, __, XX,
		XX, XX, XX,
		__, XX, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeSquare(4, []bool{
		__, __, __, __,
		XX, __, __, __,
		XX, XX, XX, XX,
		__, __, __, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeSquare(4, []bool{
		__, __, __, __,
		__, __, __, XX,
		XX, XX, XX, XX,
		__, __, __, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeSquare(4, []bool{
		__, XX, XX,
		XX, XX, XX,
		__, __, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeSquare(4, []bool{
		XX, XX, __,
		XX, XX, XX,
		__, __, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeSquare(4, []bool{
		__, __, __, __,
		XX, XX, __, __,
		__, XX, XX, XX,
		__, __, __, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeSquare(4, []bool{
		__, __, __, __,
		__, __, XX, XX,
		XX, XX, XX, __,
		__, __, __, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeSquare(4, []bool{
		__, XX, __,
		__, XX, __,
		XX, XX, XX,
	}))

	_pentominoes = append(_pentominoes, _initShapeSquare(4, []bool{
		XX, __, XX,
		XX, XX, XX,
		__, __, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeSquare(4, []bool{
		__, __, XX,
		__, __, XX,
		XX, XX, XX,
	}))

	_pentominoes = append(_pentominoes, _initShapeSquare(4, []bool{
		__, __, XX,
		__, XX, XX,
		XX, XX, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeSquare(0, []bool{
		__, XX, __,
		XX, XX, XX,
		__, XX, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeSquare(4, []bool{
		__, __, __, __,
		__, __, XX, __,
		XX, XX, XX, XX,
		__, __, __, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeSquare(4, []bool{
		__, __, __, __,
		__, XX, __, __,
		XX, XX, XX, XX,
		__, __, __, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeSquare(4, []bool{
		__, XX, XX,
		__, XX, __,
		XX, XX, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeSquare(4, []bool{
		XX, XX, __,
		__, XX, __,
		__, XX, XX,
	}))

	if want, got := _tinyminoes, shapesRotTinyminoes; !slices.Equal(want, got) {
		t.Errorf("tinyminoes mismatch: want=%+v got=%+v", want, got)
		outputShapes(_tinyminoes)
	}

	if want, got := _tetrominoes, shapesRotTetrominoes; !slices.Equal(want, got) {
		t.Errorf("tetrominoes mismatch: want=%+v got=%+v", want, got)
		outputShapes(_tetrominoes)
	}

	if want, got := _pentominoes, shapesRotPentominoes; !slices.Equal(want, got) {
		t.Errorf("pentominoes mismatch: want=%+v got=%+v", want, got)
		outputShapes(_pentominoes)
	}
}

func _initShapeSquare(rots byte, boolData []bool) shapeSquare {
	n := byte(len(boolData))

	var size byte
	var dim byte

	if rots == 1 {
		rots = 0
	}

	switch n {
	case 1:
		dim = 1
	case 4:
		dim = 2
	case 9:
		dim = 3
	case 16:
		dim = 4
	case 25:
		dim = 5
	default:
		panic("data slice has unsupported length")
	}

	var data bitarray
	for i := byte(0); i < n; i++ {
		if boolData[i] {
			data = data.set(i)
			size++
		}
	}

	if size == 0 {
		panic("empty polyomino")
	}

	return shapeSquare{
		dim:  dim,
		size: size,
		rots: rots,
		data: data,
	}
}
