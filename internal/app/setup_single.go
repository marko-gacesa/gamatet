// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	"math"

	"github.com/marko-gacesa/gamatet/game/setup"
	. "github.com/marko-gacesa/gamatet/internal/i18n"
	"github.com/marko-gacesa/gamatet/logic/menu"
)

func setupSingle(s *setup.Setup, sections *setupSections) []menu.Item {
	return []menu.Item{
		menu.NewEnum(&sections.showField, []bool{false, true}, sections.showFieldsStr,
			T(KeySetupShowFieldOptions), T(KeySetupShowFieldOptionsDesc)),
		menu.NewInteger(&s.FieldOptions.WidthSingle, setup.MinFieldWidthPerPlayer, setup.MaxFieldWidthSingle,
			"\t"+T(KeySetupFieldWidth), T(KeySetupFieldWidthDesc),
			menu.WithVisible(func() bool {
				return sections.showField
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
			}),
			withBoolStr()),
		menu.NewNumber(&s.MiscOptions.Seed, math.MinInt64, math.MaxInt64,
			"\t"+T(KeySetupRandomSeed), T(KeySetupRandomSeedDesc),
			menu.WithVisible(func() bool {
				return sections.showMisc && s.MiscOptions.CustomSeed
			})),
	}
}

func setupResultSingle(s *setup.Setup, target **setup.Setup, save bool) []menu.Item {
	var action string
	if save {
		action = T(KeySetupSave)
	} else {
		action = T(KeySetupStart)
	}
	return []menu.Item{
		menu.NewCommand(target, s,
			"", "",
			menu.WithLabelFn(func() string {
				return action + " (" + s.String() + ")"
			})),
	}
}
