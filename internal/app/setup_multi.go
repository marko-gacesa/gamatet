// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"math"

	"github.com/marko-gacesa/gamatet/game/setup"
	. "github.com/marko-gacesa/gamatet/internal/i18n"
	"github.com/marko-gacesa/gamatet/logic/menu"
)

func setupMultiPlayer(s *setup.Setup, sections *setupSections) []menu.Item {
	return []menu.Item{
		menu.NewInteger(&s.GameOptions.FieldCount, 1, setup.MaxFieldCount,
			T(KeySetupTeamCount), T(KeySetupTeamCountDesc)),
		menu.NewBool(&s.GameOptions.SamePiecesForAll,
			T(KeySetupSamePieces), T(KeySetupSamePiecesDesc),
			menu.WithVisible(func() bool {
				return s.FieldCount*s.TeamSize > 1
			})),
		menu.NewInteger(&s.GameOptions.TeamSize, 1, setup.MaxTeamSize,
			T(KeySetupTeamSize), T(KeySetupTeamSizeDesc)),
		menu.NewBool(&s.GameOptions.PlayerZones,
			"\t"+T(KeySetupPlayerZones), T(KeySetupPlayerZonesDesc),
			menu.WithVisible(func() bool {
				return s.TeamSize > 1
			})),
		menu.NewBool(&s.GameOptions.PieceCollision,
			"\t"+T(KeySetupPieceCollision), T(KeySetupPieceCollisionDesc),
			menu.WithVisible(func() bool {
				return s.TeamSize > 1 && !s.GameOptions.PlayerZones
			}),
			menu.WithDisabled(func() bool {
				return s.GameOptions.PlayerZones
			})),

		menu.NewEnum(&sections.showField, []bool{false, true}, sections.showFieldsStr,
			T(KeySetupShowFieldOptions), T(KeySetupShowFieldOptionsDesc)),
		menu.NewInteger(&s.FieldOptions.WidthSingle, setup.MinFieldWidthPerPlayer, setup.MaxFieldWidthSingle,
			"\t"+T(KeySetupFieldWidth), T(KeySetupFieldWidthDesc),
			menu.WithVisible(func() bool {
				return sections.showField && s.GameOptions.TeamSize == 1
			}),
		),
		menu.NewInteger(&s.FieldOptions.WidthPerPlayer, setup.MinFieldWidthPerPlayer, setup.MaxFieldWidthPerPlayer,
			"\t"+T(KeySetupFieldWidthPerPlayer), T(KeySetupFieldWidthPerPlayerDesc),
			menu.WithVisible(func() bool {
				return sections.showField && s.GameOptions.TeamSize > 1
			}),
		),
		menu.NewInteger(&s.FieldOptions.Height, setup.MinFieldHeight, setup.MaxFieldHeight,
			"\t"+T(KeySetupFieldHeight), T(KeySetupFieldHeightDesc),
			menu.WithVisible(func() bool {
				return sections.showField
			})),
		menu.NewInteger(&s.FieldOptions.Speed, setup.MinSpeed, setup.MaxSpeed,
			"\t"+T(KeySetupFieldSpeed), T(KeySetupFieldSpeedDesc),
			menu.WithVisible(func() bool {
				return sections.showField
			})),

		menu.NewEnum(&sections.showPiece, []bool{false, true}, sections.showPieceStr,
			T(KeySetupShowPieceOptions), T(KeySetupShowPieceOptionsDesc)),
		menu.NewEnum(&s.PieceOptions.PieceType, setup.PieceTypeAll, pieceTypeStr,
			"\t"+T(KeySetupPieceType), T(KeySetupPieceTypeDesc),
			menu.WithVisible(func() bool {
				return sections.showPiece
			})),
		menu.NewEnum(&s.PieceOptions.PieceSize, setup.PieceSizeAll, pieceSizeStr,
			"\t"+T(KeySetupPieceSize), T(KeySetupPieceSizeDesc),
			menu.WithVisible(func() bool {
				return sections.showPiece
			})),
		menu.NewInteger(&s.PieceOptions.BagSize, 1, setup.BagSizeMax,
			"\t"+T(KeySetupBagSize), T(KeySetupBagSizeDesc),
			menu.WithVisible(func() bool {
				return sections.showPiece
			})),

		menu.NewEnum(&sections.showMisc, []bool{false, true}, sections.showMiscStr,
			T(KeySetupShowMiscOptions), T(KeySetupShowMiscOptionsDesc)),
		menu.NewBool(&s.MiscOptions.CustomSeed,
			"\t"+T(KeySetupCustomRandomSeed), T(KeySetupCustomRandomSeedDesc),
			menu.WithVisible(func() bool {
				return sections.showMisc
			})),
		menu.NewNumber(&s.MiscOptions.Seed, math.MinInt64, math.MaxInt64,
			"\t"+T(KeySetupRandomSeed), T(KeySetupRandomSeedDesc),
			menu.WithVisible(func() bool {
				return sections.showMisc && s.MiscOptions.CustomSeed
			})),
	}
}

func setupResultMulti(s *setup.Setup, target **setup.Setup, maxPlayers byte, save bool) []menu.Item {
	var action string
	if save {
		action = T(KeySetupSave)
	} else {
		action = T(KeySetupStart)
	}
	return []menu.Item{
		menu.NewCommand(target, s,
			"", "",
			menu.WithVisible(func() bool {
				playerCount := s.GameOptions.PlayerCount()
				return playerCount > 1 && playerCount <= maxPlayers
			}),
			menu.WithLabelFn(func() string {
				return action + " (" + s.String() + ")"
			})),
		menu.NewStatic(
			Tf(KeySetupIssueTooManyPlayers, string('0'+rune(maxPlayers))), T(KeySetupIssueTooManyPlayersDesc), nil,
			menu.WithVisible(func() bool {
				return s.GameOptions.FieldCount*s.GameOptions.TeamSize > maxPlayers
			}),
			menu.WithDisabled(func() bool { return true })),
		menu.NewStatic(
			T(KeySetupIssueNeedAtLeast2), T(KeySetupIssueNeedAtLeast2Desc), nil,
			menu.WithVisible(func() bool {
				return s.GameOptions.FieldCount*s.GameOptions.TeamSize <= 1
			}),
			menu.WithDisabled(func() bool { return true })),
	}
}
