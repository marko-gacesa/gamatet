// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package setup

import (
	"strconv"

	"github.com/marko-gacesa/gamatet/game/piece"
)

type PieceType byte

const (
	PieceTypeRotatingPolyominoes PieceType = iota
	PieceTypeHMirroringPolyominoes
	PieceTypeVMirroringPolyominoes
)

var PieceTypeAll = []PieceType{
	PieceTypeRotatingPolyominoes,
	PieceTypeHMirroringPolyominoes,
	PieceTypeVMirroringPolyominoes,
}

func (p PieceType) String() string {
	switch p {
	case PieceTypeRotatingPolyominoes:
		return "R"
	case PieceTypeHMirroringPolyominoes:
		return "H"
	case PieceTypeVMirroringPolyominoes:
		return "V"
	default:
		return "?(" + strconv.Itoa(int(p)) + ")"
	}
}

const (
	PieceSizeMin     = piece.SizeMin
	PieceSizeMax     = piece.SizeMax
	PieceSizeDefault = piece.SizeDefault
	PieceSize3       = piece.Size3
	PieceSize4       = piece.Size4
	PieceSize5       = piece.Size5
)

var PieceSizeAll = []byte{PieceSize3, PieceSize4, PieceSize5}

const (
	BagSizeDefault = 2
	BagSizeMax     = piece.MaxBagSize
)
