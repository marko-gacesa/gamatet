// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	"strconv"

	"github.com/marko-gacesa/gamatet/game/setup"
	. "github.com/marko-gacesa/gamatet/internal/i18n"
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
		if name != "" {
			info = Tf(KeyConfigEditPlayer, i+1, name)
		} else {
			info = Tf(KeyConfigEditPlayerNoName, i+1)
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
