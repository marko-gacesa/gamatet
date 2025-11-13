// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"math"

	"github.com/marko-gacesa/gamatet/game/setup"
	"github.com/marko-gacesa/gamatet/logic/menu"
)

func setupMultiPlayer(s *setup.Setup, sections *setupSections) []menu.Item {
	return []menu.Item{
		//menu.NewEnum(&s.GameOptions.GameType, setup.GameTypeAll, setup.GameTypeNameMap,
		//	"Game type", ""),
		menu.NewInteger(&s.GameOptions.FieldCount, 1, setup.MaxFieldCount,
			"Number of teams (game fields)", ""),
		menu.NewBool(&s.GameOptions.SamePiecesForAll,
			"All players get the same pieces", "",
			menu.WithVisible(func() bool {
				return s.FieldCount*s.TeamSize > 1
			})),
		menu.NewInteger(&s.GameOptions.TeamSize, 1, setup.MaxTeamSize,
			"Players per field (team size)", ""),
		menu.NewBool(&s.GameOptions.PlayerZones,
			"\tTeam member zones", "",
			menu.WithVisible(func() bool {
				return s.TeamSize > 1
			})),
		menu.NewBool(&s.GameOptions.PieceCollision,
			"\tTeam members' piece collision", "",
			menu.WithVisible(func() bool {
				return s.TeamSize > 1 && !s.GameOptions.PlayerZones
			}),
			menu.WithDisabled(func() bool {
				return s.GameOptions.PlayerZones
			})),

		menu.NewEnum(&sections.showField, []bool{false, true}, sections.showFieldMap,
			"Show field options", ""),
		menu.NewInteger(&s.FieldOptions.WidthSingle, setup.MinFieldWidthPerPlayer, setup.MaxFieldWidthSingle,
			"\tField width", "",
			menu.WithVisible(func() bool {
				return sections.showField && s.GameOptions.TeamSize == 1
			}),
		),
		menu.NewInteger(&s.FieldOptions.WidthPerPlayer, setup.MinFieldWidthPerPlayer, setup.MaxFieldWidthPerPlayer,
			"\tField width (per team member)", "",
			menu.WithVisible(func() bool {
				return sections.showField && s.GameOptions.TeamSize > 1
			}),
		),
		menu.NewInteger(&s.FieldOptions.Height, setup.MinFieldHeight, setup.MaxFieldHeight,
			"\tField height", "",
			menu.WithVisible(func() bool {
				return sections.showField
			})),
		menu.NewInteger(&s.FieldOptions.Speed, setup.MinSpeed, setup.MaxSpeed,
			"\tInitial speed", "",
			menu.WithVisible(func() bool {
				return sections.showField
			})),

		menu.NewEnum(&sections.showPiece, []bool{false, true}, sections.showPieceMap,
			"Show piece options", ""),
		menu.NewEnum(&s.PieceOptions.PieceType, setup.PieceTypeAll, setup.PieceTypeNameMap,
			"\tPieces type", "",
			menu.WithVisible(func() bool {
				return sections.showPiece
			})),
		menu.NewEnum(&s.PieceOptions.PieceSize, setup.PieceSizeAll, setup.PieceSizeNameMap,
			"\tPieces size", "",
			menu.WithVisible(func() bool {
				return sections.showPiece
			})),
		menu.NewInteger(&s.PieceOptions.BagSize, 1, setup.BagSizeMax,
			"\tBag size", "",
			menu.WithVisible(func() bool {
				return sections.showPiece
			})),

		menu.NewEnum(&sections.showMisc, []bool{false, true}, sections.showMiscMap,
			"Show misc options", ""),
		menu.NewBool(&s.MiscOptions.CustomSeed,
			"\tCustom random number seed", "",
			menu.WithVisible(func() bool {
				return sections.showMisc
			})),
		menu.NewNumber(&s.MiscOptions.Seed, math.MinInt64, math.MaxInt64,
			"\tRandom number seed", "",
			menu.WithVisible(func() bool {
				return sections.showMisc && s.MiscOptions.CustomSeed
			})),
	}
}

func setupResultMulti(s *setup.Setup, target **setup.Setup, maxPlayers byte, save bool) []menu.Item {
	var action string
	if save {
		action = itemTextPrefixBack + "Save"
	} else {
		action = itemTextPrefixForward + "Start"
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
			"Invalid setup: Too many players - Maximum is "+string('0'+rune(maxPlayers)), "", nil,
			menu.WithVisible(func() bool {
				return s.GameOptions.FieldCount*s.GameOptions.TeamSize > maxPlayers
			}),
			menu.WithDisabled(func() bool { return true })),
		menu.NewStatic(
			"Invalid setup: Need at least 2 players", "", nil,
			menu.WithVisible(func() bool {
				return s.GameOptions.FieldCount*s.GameOptions.TeamSize <= 1
			}),
			menu.WithDisabled(func() bool { return true })),
	}
}
