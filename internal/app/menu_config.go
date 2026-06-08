// Copyright (c) 2025, 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	"strconv"

	"github.com/marko-gacesa/gamatet/game/setup"
	"github.com/marko-gacesa/gamatet/internal/config"
	. "github.com/marko-gacesa/gamatet/internal/i18n"
	"github.com/marko-gacesa/gamatet/logic/gamepad"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) menuConfig(ctx screen.Context) *menu.Menu {
	items := make([]menu.Item, 0, setup.MaxLocalPlayers+6)

	items = append(items,
		menu.NewCommand(&app.screenIDNext, routeConfigLanguage,
			T(KeyConfigLanguage), T(KeyConfigLanguageDesc)),
	)

	for i := range setup.MaxLocalPlayers {
		var info string

		name := app.cfg.LocalPlayers.Infos[i].Name
		source := "?"
		controls := "?"

		switch app.cfg.LocalPlayers.Infos[i].Input.Source {
		case config.InputSourceKeyboard:
			source = T(KeyInputSourceKeyboard)
			controls = app.cfg.LocalPlayers.Infos[i].Input.Keys.String()
		case config.InputSourceGamepad:
			gamepadIdx := app.cfg.LocalPlayers.Infos[i].Input.Gamepad
			if gamepadIdx >= 0 && gamepadIdx < gamepad.Count {
				source = T(KeyInputSourceGamepad) + " #" + strconv.Itoa(gamepadIdx+1)
				if gamepad.Gamepads[gamepadIdx].Connected {
					controls = "\"" + gamepad.Gamepads[gamepadIdx].Name + "\""
				} else {
					controls = T(KeyDeviceNotConnected)
				}
			}
		}

		if name != "" {
			info = Tf(KeyConfigEditPlayer, i+1, name, source, controls)
		} else {
			info = Tf(KeyConfigEditPlayerNoName, i+1, source, controls)
		}

		items = append(items, menu.NewCommand(&app.screenIDNext,
			route(routeConfigLocalPlayerSetupN+strconv.Itoa(i)),
			info,
			T(KeyConfigEditPlayerDesc)))
	}
	items = append(items,
		menu.NewCommand(&app.screenIDNext, routeConfigVideoSetup,
			T(KeyConfigVideoOptions), T(KeyConfigVideoOptionsDesc)),
		menu.NewCommand(&app.screenIDNext, routeSinglePlayerPresetEditMenu,
			T(KeyConfigSinglePresets), Tf(KeyConfigSinglePresetsDesc, setup.SinglePlayerPresetCount)),
		menu.NewCommand(&app.screenIDNext, routeMultiPlayerPresetEditMenu,
			T(KeyConfigMultiPresets), Tf(KeyConfigMultiPresetsDesc, setup.MultiPlayerPresetCount)),
		app.menuItemEscape(),
		app.menuItemBack(),
	)
	return menu.New(T(KeyConfigTitle), app.menuStopper(ctx), items...)
}
