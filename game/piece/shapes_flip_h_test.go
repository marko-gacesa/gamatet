// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package piece

import (
	"testing"
)

func TestShapesFlipH(t *testing.T) {
	_tinyminoes := make([]shapeRect, 0)

	_tinyminoes = append(_tinyminoes, _initShapeRect(1, 1, []bool{
		XX,
	}))

	_tinyminoes = append(_tinyminoes, _initShapeRect(2, 1, []bool{
		XX, XX,
	}))

	_tinyminoes = append(_tinyminoes, _initShapeRect(1, 2, []bool{
		XX,
		XX,
	}))

	_tinyminoes = append(_tinyminoes, _initShapeRect(3, 1, []bool{
		XX, XX, XX,
	}))

	_tinyminoes = append(_tinyminoes, _initShapeRect(1, 3, []bool{
		XX,
		XX,
		XX,
	}))

	_tinyminoes = append(_tinyminoes, _initShapeRect(2, 2, []bool{
		XX, __,
		XX, XX,
	}))

	_tinyminoes = append(_tinyminoes, _initShapeRect(2, 2, []bool{
		XX, XX,
		XX, __,
	}))

	_tetrominoes := make([]shapeRect, 0)

	_tetrominoes = append(_tetrominoes, _initShapeRect(2, 2, []bool{
		XX, XX,
		XX, XX,
	}))

	_tetrominoes = append(_tetrominoes, _initShapeRect(4, 1, []bool{
		XX, XX, XX, XX,
	}))

	_tetrominoes = append(_tetrominoes, _initShapeRect(1, 4, []bool{
		XX,
		XX,
		XX,
		XX,
	}))

	_tetrominoes = append(_tetrominoes, _initShapeRect(3, 2, []bool{
		__, XX, __,
		XX, XX, XX,
	}))

	_tetrominoes = append(_tetrominoes, _initShapeRect(3, 2, []bool{
		XX, XX, XX,
		__, XX, __,
	}))

	_tetrominoes = append(_tetrominoes, _initShapeRect(2, 3, []bool{
		XX, __,
		XX, XX,
		XX, __,
	}))

	_tetrominoes = append(_tetrominoes, _initShapeRect(2, 3, []bool{
		XX, __,
		XX, __,
		XX, XX,
	}))

	_tetrominoes = append(_tetrominoes, _initShapeRect(2, 3, []bool{
		XX, XX,
		XX, __,
		XX, __,
	}))

	_tetrominoes = append(_tetrominoes, _initShapeRect(3, 2, []bool{
		XX, __, __,
		XX, XX, XX,
	}))

	_tetrominoes = append(_tetrominoes, _initShapeRect(3, 2, []bool{
		XX, XX, XX,
		XX, __, __,
	}))

	_tetrominoes = append(_tetrominoes, _initShapeRect(2, 3, []bool{
		XX, __,
		XX, XX,
		__, XX,
	}))

	_tetrominoes = append(_tetrominoes, _initShapeRect(3, 2, []bool{
		__, XX, XX,
		XX, XX, __,
	}))

	_pentominoes := make([]shapeRect, 0)

	_pentominoes = append(_pentominoes, _initShapeRect(5, 1, []bool{
		XX, XX, XX, XX, XX,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(1, 5, []bool{
		XX,
		XX,
		XX,
		XX,
		XX,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(3, 3, []bool{
		XX, __, __,
		XX, XX, XX,
		__, XX, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(3, 3, []bool{
		__, XX, __,
		XX, XX, XX,
		XX, __, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(3, 3, []bool{
		__, XX, __,
		XX, XX, __,
		__, XX, XX,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(3, 3, []bool{
		__, XX, XX,
		XX, XX, __,
		__, XX, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(4, 2, []bool{
		XX, __, __, __,
		XX, XX, XX, XX,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(4, 2, []bool{
		XX, XX, XX, XX,
		XX, __, __, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(2, 4, []bool{
		XX, __,
		XX, __,
		XX, __,
		XX, XX,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(2, 4, []bool{
		XX, XX,
		XX, __,
		XX, __,
		XX, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(3, 2, []bool{
		XX, XX, __,
		XX, XX, XX,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(3, 2, []bool{
		XX, XX, XX,
		XX, XX, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(2, 3, []bool{
		XX, __,
		XX, XX,
		XX, XX,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(2, 3, []bool{
		XX, XX,
		XX, XX,
		XX, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(4, 2, []bool{
		XX, XX, __, __,
		__, XX, XX, XX,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(4, 2, []bool{
		__, XX, XX, XX,
		XX, XX, __, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(2, 4, []bool{
		XX, __,
		XX, __,
		XX, XX,
		__, XX,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(2, 4, []bool{
		__, XX,
		XX, XX,
		XX, __,
		XX, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(3, 3, []bool{
		__, XX, __,
		__, XX, __,
		XX, XX, XX,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(3, 3, []bool{
		XX, XX, XX,
		__, XX, __,
		__, XX, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(3, 3, []bool{
		XX, __, __,
		XX, XX, XX,
		XX, __, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(3, 2, []bool{
		XX, __, XX,
		XX, XX, XX,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(3, 2, []bool{
		XX, XX, XX,
		XX, __, XX,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(2, 3, []bool{
		XX, XX,
		XX, __,
		XX, XX,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(3, 3, []bool{
		XX, __, __,
		XX, __, __,
		XX, XX, XX,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(3, 3, []bool{
		XX, XX, XX,
		XX, __, __,
		XX, __, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(3, 3, []bool{
		XX, __, __,
		XX, XX, __,
		__, XX, XX,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(3, 3, []bool{
		__, XX, XX,
		XX, XX, __,
		XX, __, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(3, 3, []bool{
		__, XX, __,
		XX, XX, XX,
		__, XX, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(4, 2, []bool{
		__, XX, __, __,
		XX, XX, XX, XX,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(4, 2, []bool{
		XX, XX, XX, XX,
		__, XX, __, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(2, 4, []bool{
		XX, __,
		XX, __,
		XX, XX,
		XX, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(2, 4, []bool{
		XX, __,
		XX, XX,
		XX, __,
		XX, __,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(3, 3, []bool{
		XX, XX, __,
		__, XX, __,
		__, XX, XX,
	}))

	_pentominoes = append(_pentominoes, _initShapeRect(3, 3, []bool{
		__, XX, XX,
		__, XX, __,
		XX, XX, __,
	}))

	if want, got := _tinyminoes, shapesFlipHTinyminoes; !shapesRectEqual(want, got) {
		t.Error("tinymino mismatch")
		outputShapes(_tinyminoes)
	}

	if want, got := _tetrominoes, shapesFlipHTetrominoes; !shapesRectEqual(want, got) {
		t.Error("tetrominoes mismatch")
		outputShapes(_tetrominoes)
	}

	if want, got := _pentominoes, shapesFlipHPentominoes; !shapesRectEqual(want, got) {
		t.Error("pentominoes mismatch")
		outputShapes(_pentominoes)
	}
}
