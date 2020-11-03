// Copyright (c) 2020 by Marko Gaćeša

package piece

import (
	"reflect"
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

func TestShapes(t *testing.T) {
	const XX = true
	const _I = false

	_monomino := _initPolyomino(0, []bool{
		XX,
	})

	_domino := _initPolyomino(2, []bool{
		XX, XX,
		_I, _I,
	})

	_trominoes := make([]polyominoShape, 2)

	_trominoes[0] = _initPolyomino(2, []bool{
		_I, _I, _I,
		XX, XX, XX,
		_I, _I, _I,
	})

	_trominoes[1] = _initPolyomino(4, []bool{
		_I, XX,
		XX, XX,
	})

	_tetrominoes := make([]polyominoShape, 7)

	for i := 0; i < len(_tetrominoes); i++ {
		switch i {
		case TetrominoO:
			_tetrominoes[i] = _initPolyomino(0, []bool{
				XX, XX,
				XX, XX,
			})
		case TetrominoI:
			_tetrominoes[i] = _initPolyomino(2, []bool{
				_I, _I, _I, _I,
				XX, XX, XX, XX,
				_I, _I, _I, _I,
				_I, _I, _I, _I,
			})
		case TetrominoT:
			_tetrominoes[i] = _initPolyomino(4, []bool{
				_I, XX, _I,
				XX, XX, XX,
				_I, _I, _I,
			})
		case TetrominoJ:
			_tetrominoes[i] = _initPolyomino(4, []bool{
				XX, _I, _I,
				XX, XX, XX,
				_I, _I, _I,
			})
		case TetrominoL:
			_tetrominoes[i] = _initPolyomino(4, []bool{
				_I, _I, XX,
				XX, XX, XX,
				_I, _I, _I,
			})
		case TetrominoS:
			_tetrominoes[i] = _initPolyomino(2, []bool{
				_I, XX, XX,
				XX, XX, _I,
				_I, _I, _I,
			})
		case TetrominoZ:
			_tetrominoes[i] = _initPolyomino(2, []bool{
				XX, XX, _I,
				_I, XX, XX,
				_I, _I, _I,
			})
		}
	}

	_pentominoes := make([]polyominoShape, 18)

	for i := 0; i < len(_pentominoes); i++ {
		switch i {
		case 0:
			_pentominoes[i] = _initPolyomino(2, []bool{
				_I, _I, _I, _I, _I,
				_I, _I, _I, _I, _I,
				XX, XX, XX, XX, XX,
				_I, _I, _I, _I, _I,
				_I, _I, _I, _I, _I,
			})
		case 1:
			_pentominoes[i] = _initPolyomino(4, []bool{
				XX, _I, _I,
				XX, XX, XX,
				_I, XX, _I,
			})
		case 2:
			_pentominoes[i] = _initPolyomino(4, []bool{
				_I, _I, XX,
				XX, XX, XX,
				_I, XX, _I,
			})
		case 3:
			_pentominoes[i] = _initPolyomino(4, []bool{
				_I, _I, _I, _I,
				XX, _I, _I, _I,
				XX, XX, XX, XX,
				_I, _I, _I, _I,
			})
		case 4:
			_pentominoes[i] = _initPolyomino(4, []bool{
				_I, _I, _I, _I,
				_I, _I, _I, XX,
				XX, XX, XX, XX,
				_I, _I, _I, _I,
			})
		case 5:
			_pentominoes[i] = _initPolyomino(4, []bool{
				_I, XX, XX,
				XX, XX, XX,
				_I, _I, _I,
			})
		case 6:
			_pentominoes[i] = _initPolyomino(4, []bool{
				XX, XX, _I,
				XX, XX, XX,
				_I, _I, _I,
			})
		case 7:
			_pentominoes[i] = _initPolyomino(4, []bool{
				_I, _I, _I, _I,
				XX, XX, _I, _I,
				_I, XX, XX, XX,
				_I, _I, _I, _I,
			})
		case 8:
			_pentominoes[i] = _initPolyomino(4, []bool{
				_I, _I, _I, _I,
				_I, _I, XX, XX,
				XX, XX, XX, _I,
				_I, _I, _I, _I,
			})
		case 9:
			_pentominoes[i] = _initPolyomino(4, []bool{
				_I, XX, _I,
				_I, XX, _I,
				XX, XX, XX,
			})
		case 10:
			_pentominoes[i] = _initPolyomino(4, []bool{
				XX, _I, XX,
				XX, XX, XX,
				_I, _I, _I,
			})
		case 11:
			_pentominoes[i] = _initPolyomino(4, []bool{
				_I, _I, XX,
				_I, _I, XX,
				XX, XX, XX,
			})
		case 12:
			_pentominoes[i] = _initPolyomino(4, []bool{
				_I, _I, XX,
				_I, XX, XX,
				XX, XX, _I,
			})
		case 13:
			_pentominoes[i] = _initPolyomino(0, []bool{
				_I, XX, _I,
				XX, XX, XX,
				_I, XX, _I,
			})
		case 14:
			_pentominoes[i] = _initPolyomino(4, []bool{
				_I, _I, _I, _I,
				_I, _I, XX, _I,
				XX, XX, XX, XX,
				_I, _I, _I, _I,
			})
		case 15:
			_pentominoes[i] = _initPolyomino(4, []bool{
				_I, _I, _I, _I,
				_I, XX, _I, _I,
				XX, XX, XX, XX,
				_I, _I, _I, _I,
			})
		case 16:
			_pentominoes[i] = _initPolyomino(4, []bool{
				_I, XX, XX,
				_I, XX, _I,
				XX, XX, _I,
			})
		case 17:
			_pentominoes[i] = _initPolyomino(4, []bool{
				XX, XX, _I,
				_I, XX, _I,
				_I, XX, XX,
			})
		}
	}

	if want, got := _monomino, monomino; want != got {
		t.Errorf("monomino mismatch: want=%+v got=%+v", want, got)
	}

	if want, got := _domino, domino; want != got {
		t.Errorf("domino mismatch: want=%+v got=%+v", want, got)
	}

	if want, got := _trominoes, trominoes; !reflect.DeepEqual(want, got) {
		t.Errorf("trominoes mismatch: want=%+v got=%+v", want, got)
	}

	if want, got := _tetrominoes, tetrominoes; !reflect.DeepEqual(want, got) {
		t.Errorf("tetrominoes mismatch: want=%+v got=%+v", want, got)
	}

	if want, got := _pentominoes, pentominoes; !reflect.DeepEqual(want, got) {
		t.Errorf("pentominoes mismatch: want=%+v got=%+v", want, got)
	}
}

func _initPolyomino(rots byte, boolData []bool) polyominoShape {
	n := len(boolData)
	size := byte(0)
	dim := 0

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
	for i := 0; i < n; i++ {
		if boolData[i] {
			data = data.set(i)
			size++
		}
	}

	if size == 0 {
		panic("empty polyomino")
	}

	return polyominoShape{
		dim:  dim,
		size: size,
		rots: rots,
		data: data,
	}
}
