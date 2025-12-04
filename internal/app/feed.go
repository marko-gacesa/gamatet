// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	"github.com/marko-gacesa/gamatet/game/piece"
	"github.com/marko-gacesa/gamatet/game/setup"
)

func Feed(s setup.Setup) piece.Feed {
	if s.GameType != setup.GameTypeFallingPolyominoes {
		return piece.QFeed{}
	}

	bagSize := int(s.PieceOptions.BagSize)
	seed := int(s.MiscOptions.Seed)
	isSingle := s.GameOptions.FieldCount == 1 && s.GameOptions.TeamSize == 1

	var color piece.Color
	if isSingle {
		color = piece.DefaultColor{}
	} else {
		color = piece.NewRandomColor(setup.ColorRGB[:], seed)
	}

	switch s.PieceOptions.PieceType {
	case setup.PieceTypeRotatingPolyominoes:
		return piece.NewRotTetrominoFeed(s.PieceOptions.PieceSize, bagSize, seed, color)
	case setup.PieceTypeVMirroringPolyominoes:
		return piece.NewFlipVFeed(s.PieceOptions.PieceSize, bagSize, seed, color)
	case setup.PieceTypeHMirroringPolyominoes:
		return piece.NewFlipHFeed(s.PieceOptions.PieceSize, bagSize, seed, color)
	}

	return piece.QFeed{}
}
