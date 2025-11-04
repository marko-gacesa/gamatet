// Copyright (c) 2024, 2025 by Marko Gaćeša

package app

import (
	"fmt"
	"gamatet/game/setup"
	"gamatet/logic/menu"
	"gamatet/logic/screen"
)

func (app *App) menuSinglePlayer(ctx screen.Context) *menu.Menu {
	var items []menu.Item
	items = append(items, app.menuItemEscape())
	presetRoutes := []route{
		routeGameSinglePreset1,
		routeGameSinglePreset2,
		routeGameSinglePreset3,
		routeGameSinglePreset4,
		routeGameSinglePreset5,
	}
	for i := range setup.SinglePlayerPresetCount {
		name := fmt.Sprintf("Preset %d: %s [%s]",
			i+1,
			app.cfg.Presets.Single[i].Name,
			app.cfg.Presets.Single[i].String(),
		)
		items = append(items, menu.NewCommand(&app.screenIDNext, presetRoutes[i], name, ""))
	}
	items = append(items, menu.NewCommand(&app.screenIDNext, routeGameSinglePreset1, "Custom game", ""))
	items = append(items, menu.NewCommand(&app.screenIDNext, routeGameSinglePreset1, "Edit presets", ""))
	items = append(items, app.menuItemBack())

	return menu.New("Single Player", app.menuStopper(ctx), items...)
}
