// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package setup

import "github.com/marko-gacesa/gamatet/game/piece"

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
