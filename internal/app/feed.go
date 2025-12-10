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
	isBattle := s.GameOptions.FieldCount > 1

	var (
		color  piece.Color
		shapes []piece.ShapeAny
		wb     piece.TypeWeights
		ws     piece.TypeWeights
	)

	if isSingle {
		color = piece.DefaultColor{}
	} else {
		color = piece.NewRandomColor(setup.ColorRGB[:], seed)
	}

	if s.PieceOptions.Shooters {
		shapes = ShapesWithShooters(s.PieceOptions.PieceType, s.PieceOptions.PieceSize)
	} else {
		shapes = Shapes(s.PieceOptions.PieceType, s.PieceOptions.PieceSize)
	}

	if s.PieceOptions.SpecialBlocks {
		wb = piece.NewTypeWeights(17, 2, 1, 0, 0)
	} else {
		wb = piece.NewTypeWeights(1, 0, 0, 0, 0)
	}

	if isBattle {
		ws = piece.NewTypeWeights(0, 1, 1, 1, 1)
	} else {
		ws = piece.NewTypeWeights(0, 3, 2, 0, 0)
	}

	return piece.MixedFeed(wb, ws, bagSize, seed, color, shapes...)
}

func Shapes(pieceType setup.PieceType, pieceSize byte) []piece.ShapeAny {
	var shapes []piece.ShapeAny
	switch pieceType {
	case setup.PieceTypeRotatingPolyominoes:
		for _, shape := range piece.GetRotShapes(pieceSize) {
			shapes = append(shapes, shape)
		}
	case setup.PieceTypeVMirroringPolyominoes:
		for _, shape := range piece.GetFlipVShapes(pieceSize) {
			shapes = append(shapes, shape)
		}
	case setup.PieceTypeHMirroringPolyominoes:
		for _, shape := range piece.GetFlipHShapes(pieceSize) {
			shapes = append(shapes, shape)
		}
	}

	return shapes
}

func ShapesWithShooters(pieceType setup.PieceType, pieceSize byte) []piece.ShapeAny {
	var shapes []piece.ShapeAny
	switch pieceType {
	case setup.PieceTypeRotatingPolyominoes:
		switch pieceSize {
		case setup.PieceSize3: // 20:1
			for i := 0; i < 5; i++ {
				shapes = append(shapes, Shapes(pieceType, pieceSize)...)
			}
			shapes = append(shapes, piece.ShapeShooter{})
		case setup.PieceSize4: // 21:1
			for i := 0; i < 3; i++ {
				shapes = append(shapes, Shapes(pieceType, pieceSize)...)
			}
			shapes = append(shapes, piece.ShapeShooter{})
		case setup.PieceSize5: // 40:2
			shapes = append(shapes, Shapes(pieceType, pieceSize)...)
			shapes = append(shapes, piece.ShapeShooter{})
			shapes = append(shapes, piece.ShapeShooter{})
		}
	case setup.PieceTypeHMirroringPolyominoes, setup.PieceTypeVMirroringPolyominoes:
		switch pieceSize {
		case setup.PieceSize3: // 21:1
			for i := 0; i < 3; i++ {
				shapes = append(shapes, Shapes(pieceType, pieceSize)...)
			}
			shapes = append(shapes, piece.ShapeShooter{})
		case setup.PieceSize4: // 19:1
			shapes = append(shapes, Shapes(pieceType, pieceSize)...)
			shapes = append(shapes, piece.ShapeShooter{})
		case setup.PieceSize5: // 73:4
			shapes = append(shapes, Shapes(pieceType, pieceSize)...)
			shapes = append(shapes, piece.ShapeShooter{})
			shapes = append(shapes, piece.ShapeShooter{})
			shapes = append(shapes, piece.ShapeShooter{})
			shapes = append(shapes, piece.ShapeShooter{})
		}
	}

	return shapes
}
