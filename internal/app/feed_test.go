package app

import (
	"testing"

	"github.com/marko-gacesa/gamatet/game/piece"
	"github.com/marko-gacesa/gamatet/game/setup"
)

func TestShapesWithShooters(t *testing.T) {
	for pieceType := setup.PieceType(0); pieceType <= 2; pieceType++ {
		for pieceSize := byte(setup.PieceSizeMin); pieceSize <= setup.PieceSizeMax; pieceSize++ {
			shapes := ShapesWithShooters(pieceType, pieceSize)
			var countShooters int
			var countRegular int
			for _, shape := range shapes {
				if _, ok := shape.(piece.ShapeShooter); ok {
					countShooters++
				} else {
					countRegular++
				}
			}

			ratio := float64(countRegular) / float64(countShooters)

			if ratio < 18.0 || ratio > 22.0 {
				t.Errorf("type=%s size=%d: Ratio=%0.3f out of range (regulars=%d shooters=%d)",
					pieceType, pieceSize,
					ratio, countRegular, countShooters)
			}
		}
	}
}
