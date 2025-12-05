// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package piece

import (
	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/logic/random"
)

func MixedFeed(wb, ws TypeWeights, bagSize int, seed int, c Color, shapes ...any) Feed {
	return NewGenericFeed(bagSize, seed, len(shapes), func(idx, shapeIdx uint, playerIndex byte) Piece {
		switch s := shapes[shapeIdx].(type) {
		case shapeSquare:
			b := wb.Block(c, idx, shapeIdx, seed, playerIndex)
			return &polyominoRot{shapeSquare: s, block: b}
		case shapeRectV:
			b := wb.Block(c, idx, shapeIdx, seed, playerIndex)
			return &polyominoFlipV{shapeRect: shapeRect(s), block: b}
		case shapeRectH:
			b := wb.Block(c, idx, shapeIdx, seed, playerIndex)
			return &polyominoFlipH{shapeRect: shapeRect(s), block: b}
		case ShapeShooter:
			b := ws.Block(RockColor{}, idx, shapeIdx, seed, playerIndex)
			return Shooter(5, b.Type)
		default:
			return NewQ(block.Rock)
		}
	})
}

type TypeWeights struct {
	w random.Weights[byte]
}

func NewTypeWeights(rock, lava, acid, wave, curl byte) TypeWeights {
	var w random.Weights[byte]
	w.Add(rock)
	w.Add(lava)
	w.Add(acid)
	w.Add(wave)
	w.Add(curl)
	return TypeWeights{w}
}

func (w TypeWeights) Block(c Color, idx, shapeIdx uint, seed int, playerIndex byte) block.Block {
	r := random.New(idx*13, uint(seed))
	switch w.w.Random(r) {
	case 1:
		return block.Lava
	case 2:
		return block.Acid
	case 3:
		return block.Wave
	case 4:
		return block.Curl
	default:
		return block.Block{Type: block.TypeRock, Hardness: 0, Color: c.Color(shapeIdx, playerIndex)}
	}
}
