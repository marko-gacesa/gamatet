// Copyright (c) 2025 by Marko Gaćeša

package setup

import "gamatet/game/piece"

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

var PieceTypeNameMap = map[PieceType]string{
	PieceTypeRotatingPolyominoes:   "Rotating Polyominoes",
	PieceTypeHMirroringPolyominoes: "Horizontal Mirroring Polyominoes",
	PieceTypeVMirroringPolyominoes: "Vertical Mirroring Polyominoes",
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

var PieceSizeNameMap = map[byte]string{
	PieceSize3: "Miniminoes",
	PieceSize4: "Tetrominoes",
	PieceSize5: "Pentominoes",
}

const (
	BagSizeDefault = 2
	BagSizeMax     = piece.MaxBagSize
)
