// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"fmt"
	"strconv"

	"github.com/marko-gacesa/gamatet/game/setup"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) menuConfig(ctx screen.Context) *menu.Menu {
	items := make([]menu.Item, 0, setup.MaxLocalPlayers+6)
	for i := range setup.MaxLocalPlayers {
		info := app.cfg.LocalPlayers.Infos[i].Name
		items = append(items, menu.NewCommand(&app.screenIDNext,
			route(routeConfigLocalPlayerN+strconv.Itoa(i)),
			fmt.Sprintf("Edit local player %d: %s", i+1, info),
			""))
	}
	items = append(items,
		menu.NewCommand(&app.screenIDNext, routeConfigVideoSetup, "Video options", ""),
		menu.NewCommand(&app.screenIDNext, routeSinglePlayerPresetEditMenu, "Edit single game presets", ""),
		menu.NewCommand(&app.screenIDNext, routeMultiPlayerPresetEditMenu, "Edit multiplayer game presets", ""),
		app.menuItemEscape(),
		app.menuItemBack(),
	)
	return menu.New("Configuration", app.menuStopper(ctx), items...)
}
