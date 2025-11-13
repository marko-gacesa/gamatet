// Copyright (c) 2024, 2025 by Marko Gaćeša

package app

import (
	"fmt"
	"strconv"

	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) menuMultiPlayerLocal(ctx screen.Context) *menu.Menu {
	var items []menu.Item
	items = append(items, app.menuItemEscape())

	presetRoutes := make([]route, len(app.cfg.Presets.Multi))
	for i := range app.cfg.Presets.Multi {
		presetRoutes[i] = routeMultiPlayerLocalPresetGameN + route(strconv.Itoa(i))
	}

	for i := range app.cfg.Presets.Multi {
		name := fmt.Sprintf("Start Preset %d: %s [%s]",
			i+1,
			app.cfg.Presets.Multi[i].Name,
			app.cfg.Presets.Multi[i].String(),
		)
		items = append(items, menu.NewCommand(&app.screenIDNext, presetRoutes[i], name, ""))
	}
	items = append(items, menu.NewCommand(&app.screenIDNext, routeMultiPayerLocalCustomSetup, "Custom game", ""))
	items = append(items, menu.NewCommand(&app.screenIDNext, routeMultiPlayerPresetEditMenu, "Edit presets", ""))
	items = append(items, app.menuItemBack())

	return menu.New("Multi Player Local", app.menuStopper(ctx), items...)
}
