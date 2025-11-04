// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"gamatet/game/piece"
	"gamatet/game/setup"
)

func Feed(s setup.Setup) piece.Feed {
	if s.GameType != setup.GameTypeFallingPolyominoes {
		return piece.QFeed{}
	}

	badSize := int(s.PieceOptions.BagSize)
	seed := int(s.MiscOptions.Seed)
	isSingle := s.GameOptions.FieldCount == 1 && s.GameOptions.TeamSize == 1

	var color piece.Color
	if isSingle {
		color = piece.DefaultColor{}
	} else {
		color = piece.NewRandomColor(0.1, 0.5, 1.0, seed)
	}

	switch s.PieceOptions.PieceType {
	case setup.PieceTypeRotatingPolyominoes:
		return piece.NewRotTetrominoFeed(s.PieceOptions.PieceSize, badSize, seed, color)
	case setup.PieceTypeVMirroringPolyominoes:
		return piece.NewFlipVFeed(s.PieceOptions.PieceSize, badSize, seed, color)
	case setup.PieceTypeHMirroringPolyominoes:
		return piece.NewFlipHFeed(s.PieceOptions.PieceSize, badSize, seed, color)
	}

	return piece.QFeed{}
}
