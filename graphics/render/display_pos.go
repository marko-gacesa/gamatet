// Copyright (c) 2024, 2025 by Marko Gaćeša

package render

import (
	"gamatet/game/field"
)

type DisplayPosition byte

const (
	DisplayPositionOff         DisplayPosition = 0
	DisplayPositionTopLeft     DisplayPosition = 1
	DisplayPositionTopRight    DisplayPosition = 2
	DisplayPositionBottomLeft  DisplayPosition = 3
	DisplayPositionBottomRight DisplayPosition = 4
)

type PreferredSide byte

const (
	PreferredSideTop PreferredSide = iota
	PreferredSideLeft
	PreferredSideRight
	PreferredSideBottom
)

func (side PreferredSide) Pos() [field.MaxPieces]DisplayPosition {
	switch side % field.MaxPieces {
	case PreferredSideTop:
		return [4]DisplayPosition{
			DisplayPositionTopLeft,
			DisplayPositionTopRight,
			DisplayPositionBottomLeft,
			DisplayPositionBottomRight,
		}
	case PreferredSideLeft:
		return [4]DisplayPosition{
			DisplayPositionTopLeft,
			DisplayPositionBottomLeft,
			DisplayPositionTopRight,
			DisplayPositionBottomRight,
		}
	case PreferredSideRight:
		return [4]DisplayPosition{
			DisplayPositionTopRight,
			DisplayPositionBottomRight,
			DisplayPositionTopLeft,
			DisplayPositionBottomLeft,
		}
	case PreferredSideBottom:
		return [4]DisplayPosition{
			DisplayPositionBottomLeft,
			DisplayPositionBottomRight,
			DisplayPositionTopLeft,
			DisplayPositionTopRight,
		}
	}
	panic("unreachable")
}

func (side PreferredSide) PosN(n int) [field.MaxPieces]DisplayPosition {
	a := side.Pos()
	for i := n; i < field.MaxPieces; i++ {
		a[i] = DisplayPositionOff
	}
	return a
}
