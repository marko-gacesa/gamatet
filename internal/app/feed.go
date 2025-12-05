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

	var shapes []any

	switch s.PieceOptions.PieceType {
	case setup.PieceTypeRotatingPolyominoes:
		for _, shape := range piece.GetRotShapes(s.PieceOptions.PieceSize) {
			shapes = append(shapes, shape)
		}
		shapes = append(shapes, piece.ShapeShooter{})
	case setup.PieceTypeVMirroringPolyominoes:
		for _, shape := range piece.GetFlipVShapes(s.PieceOptions.PieceSize) {
			shapes = append(shapes, shape)
		}
	case setup.PieceTypeHMirroringPolyominoes:
		for _, shape := range piece.GetFlipHShapes(s.PieceOptions.PieceSize) {
			shapes = append(shapes, shape)
		}
	}

	if len(shapes) == 0 {
		shapes = append(shapes, piece.QFeed{})
	}

	wb := piece.NewTypeWeights(1, 0, 0, 0, 0)
	ws := piece.NewTypeWeights(1, 1, 1, 1, 1)

	return piece.MixedFeed(wb, ws, bagSize, seed, color, shapes...)
}
