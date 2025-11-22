// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	"github.com/marko-gacesa/gamatet/game/setup"
	"github.com/marko-gacesa/gamatet/internal/config/key"
	. "github.com/marko-gacesa/gamatet/internal/i18n"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) menuConfigLocalPlayer(ctx screen.Context, idx int) *menu.Menu {
	if idx < 0 || idx >= setup.MaxLocalPlayers {
		return app.menuErrorText(ctx, T(KeyErrorPlayerIndexOutOfRange))
	}

	stringFn := func(k key.Key) string { return key.Map[k] }
	convertFn := func(b byte) key.Key { return key.Key(b) }

	info := &app.cfg.LocalPlayers.Infos[idx]

	section := newSetupKeySection()
	section.refresh(&info.Input)
	showKeysFn := func() bool {
		return section.showKeys
	}

	items := make([]menu.Item, 0, 10)
	items = append(items,
		menu.NewText(&info.Name, setup.MaxLenName, setup.MaxLenName,
			T(KeyConfigPlayerName), T(KeyConfigPlayerNameDesc)),
		menu.NewEnum(&section.showKeys, []bool{false, true}, section.showKeysStr,
			T(KeyConfigPlayerControl), T(KeyConfigPlayerControlDesc)),
		menu.NewKey(&info.Input.Left,
			"\t"+T(KeyConfigPlayerKeyLeft), T(KeyConfigPlayerKeyLeftDesc),
			stringFn, convertFn, menu.WithVisible(showKeysFn)),
		menu.NewKey(&info.Input.Right,
			"\t"+T(KeyConfigPlayerKeyRight), T(KeyConfigPlayerKeyRightDesc),
			stringFn, convertFn, menu.WithVisible(showKeysFn)),
		menu.NewKey(&info.Input.Activate,
			"\t"+T(KeyConfigPlayerKeyActivate), T(KeyConfigPlayerKeyActivateDesc),
			stringFn, convertFn, menu.WithVisible(showKeysFn)),
		menu.NewKey(&info.Input.Boost,
			"\t"+T(KeyConfigPlayerKeyBoost), T(KeyConfigPlayerKeyBoostDesc),
			stringFn, convertFn, menu.WithVisible(showKeysFn)),
		menu.NewKey(&info.Input.Drop,
			"\t"+T(KeyConfigPlayerKeyDrop), T(KeyConfigPlayerKeyDropDesc),
			stringFn, convertFn, menu.WithVisible(showKeysFn)),
		menu.NewEnum(&info.GameConfig.RotationDirectionCW,
			[]bool{false, true}, rotationDirCWStr,
			T(KeyConfigPlayerRotationDirection), T(KeyConfigPlayerRotationDirectionDesc)),
		app.menuItemEscape(),
		app.menuItemBack(),
	)
	return menu.New(Tf(KeyConfigPlayerTitle, idx+1), func(m *menu.Menu) {
		app.configStopper(ctx)(m)
		section.refresh(&info.Input)
	}, items...)
}
