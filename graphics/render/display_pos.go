// Copyright (c) 2024, 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package render

import (
	"github.com/marko-gacesa/gamatet/game/field"
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
	PreferredSideTopL2R PreferredSide = iota
	PreferredSideTopR2L
	PreferredSideLeftT2B
	PreferredSideLeftB2T
	PreferredSideRightT2B
	PreferredSideRightB2T
	PreferredSideBottomL2R
	PreferredSideBottomR2L

	PreferredSideCount
)

func (side PreferredSide) String() string {
	switch side {
	case PreferredSideTopL2R:
		return "T-L2R"
	case PreferredSideTopR2L:
		return "T-R2L"
	case PreferredSideLeftT2B:
		return "L-T2B"
	case PreferredSideLeftB2T:
		return "L-B2T"
	case PreferredSideRightT2B:
		return "R-T2B"
	case PreferredSideRightB2T:
		return "R-B2T"
	case PreferredSideBottomL2R:
		return "B-L2R"
	case PreferredSideBottomR2L:
		return "B-R2L"
	default:
		return "invalid"
	}
}

func (side PreferredSide) PieceCorners(playerCount int) PieceCornerList {
	a := PieceCornerLists[side%PreferredSideCount]
	for i := playerCount; i < field.MaxPieces; i++ {
		a[i] = DisplayPositionOff
	}
	return a
}

type PieceCornerList [field.MaxPieces]DisplayPosition

var PieceCornerLists = [PreferredSideCount]PieceCornerList{
	{
		DisplayPositionTopLeft,
		DisplayPositionTopRight,
		DisplayPositionBottomLeft,
		DisplayPositionBottomRight,
	},
	{
		DisplayPositionTopRight,
		DisplayPositionTopLeft,
		DisplayPositionBottomRight,
		DisplayPositionBottomLeft,
	},
	{
		DisplayPositionTopLeft,
		DisplayPositionBottomLeft,
		DisplayPositionTopRight,
		DisplayPositionBottomRight,
	},
	{
		DisplayPositionBottomLeft,
		DisplayPositionTopLeft,
		DisplayPositionBottomRight,
		DisplayPositionTopRight,
	},
	{
		DisplayPositionTopRight,
		DisplayPositionBottomRight,
		DisplayPositionTopLeft,
		DisplayPositionBottomLeft,
	},
	{
		DisplayPositionBottomRight,
		DisplayPositionTopRight,
		DisplayPositionBottomLeft,
		DisplayPositionTopLeft,
	},
	{
		DisplayPositionBottomLeft,
		DisplayPositionBottomRight,
		DisplayPositionTopLeft,
		DisplayPositionTopRight,
	},
	{
		DisplayPositionBottomRight,
		DisplayPositionBottomLeft,
		DisplayPositionTopRight,
		DisplayPositionTopLeft,
	},
}
