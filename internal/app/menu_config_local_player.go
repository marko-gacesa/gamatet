// Copyright (c) 2025, 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	"fmt"

	"github.com/marko-gacesa/gamatet/game/setup"
	"github.com/marko-gacesa/gamatet/internal/config"
	"github.com/marko-gacesa/gamatet/internal/config/key"
	. "github.com/marko-gacesa/gamatet/internal/i18n"
	"github.com/marko-gacesa/gamatet/logic/gamepad"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) menuConfigLocalPlayer(ctx screen.Context, idx int) *menu.Menu {
	if idx < 0 || idx >= setup.MaxLocalPlayers {
		return app.menuError(ctx, errorPlayerIndexOutOfRange)
	}

	stringFn := func(k key.Key) string { return key.Map[k] }
	convertFn := func(b byte) key.Key { return key.Key(b) }

	info := &app.cfg.LocalPlayers.Infos[idx]

	section := newSetupKeySection()
	section.refresh(&info.Input.Keys)
	showKeysFn := func() bool {
		return info.Input.Source == config.InputSourceKeyboard && section.showKeys
	}

	items := make([]menu.Item, 0, 16)
	items = append(items,
		menu.NewText(&info.Name, setup.MaxLenName, setup.MaxLenName,
			T(KeyConfigPlayerName), T(KeyConfigPlayerNameDesc)),
		menu.NewEnum(&info.Input.Source,
			[]config.InputSource{config.InputSourceKeyboard, config.InputSourceGamepad},
			func(source config.InputSource) string {
				switch source {
				case config.InputSourceKeyboard:
					return T(KeyInputSourceKeyboard)
				case config.InputSourceGamepad:
					return T(KeyInputSourceGamepad)
				default:
					return "?"
				}
			},
			T(KeyConfigPlayerInputDevice), T(KeyConfigPlayerInputDeviceDesc)),
		menu.NewEnum(&section.showKeys, []bool{false, true}, section.showKeysStr,
			T(KeyConfigPlayerControl), T(KeyConfigPlayerControlDesc), menu.WithVisible(func() bool {
				return info.Input.Source == config.InputSourceKeyboard
			})),
		menu.NewKey(&info.Input.Keys.Left,
			"\t"+T(KeyConfigPlayerKeyLeft), T(KeyConfigPlayerKeyLeftDesc),
			stringFn, convertFn, menu.WithVisible(showKeysFn)),
		menu.NewKey(&info.Input.Keys.Right,
			"\t"+T(KeyConfigPlayerKeyRight), T(KeyConfigPlayerKeyRightDesc),
			stringFn, convertFn, menu.WithVisible(showKeysFn)),
		menu.NewKey(&info.Input.Keys.Activate,
			"\t"+T(KeyConfigPlayerKeyActivate), T(KeyConfigPlayerKeyActivateDesc),
			stringFn, convertFn, menu.WithVisible(showKeysFn)),
		menu.NewKey(&info.Input.Keys.Boost,
			"\t"+T(KeyConfigPlayerKeyBoost), T(KeyConfigPlayerKeyBoostDesc),
			stringFn, convertFn, menu.WithVisible(showKeysFn)),
		menu.NewKey(&info.Input.Keys.Drop,
			"\t"+T(KeyConfigPlayerKeyDrop), T(KeyConfigPlayerKeyDropDesc),
			stringFn, convertFn, menu.WithVisible(showKeysFn)),
		menu.NewEnum(&info.Input.Gamepad, gamepad.Indices,
			func(gamepadIdx int) string {
				if gamepadIdx < 0 || gamepadIdx >= gamepad.Count {
					return "?"
				}
				if gamepad.Gamepads[gamepadIdx].Connected {
					return fmt.Sprintf("#%d (%q)", gamepadIdx+1, gamepad.Gamepads[gamepadIdx].Name)
				}

				return fmt.Sprintf("#%d (%s)", gamepadIdx+1, T(KeyDeviceNotConnected))
			},
			T(KeyInputSourceGamepad), "",
			menu.WithVisible(func() bool {
				return info.Input.Source == config.InputSourceGamepad
			})),
		menu.NewEnum(&info.GameConfig.RotationDirectionCW,
			[]bool{false, true}, rotationDirCWStr,
			T(KeyConfigPlayerRotationDirection), T(KeyConfigPlayerRotationDirectionDesc)),
		app.menuItemEscape(),
		app.menuItemBack(),
	)
	return menu.New(Tf(KeyConfigPlayerTitle, idx+1), func(m *menu.Menu) {
		app.configStopper(ctx)(m)
		section.refresh(&info.Input.Keys)
	}, items...)
}
