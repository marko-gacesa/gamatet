// Copyright (c) 2024, 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	"strconv"

	"github.com/marko-gacesa/gamatet/game/setup"
	. "github.com/marko-gacesa/gamatet/internal/i18n"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) menuMultiPlayerDirectIPHostMenu(ctx screen.Context) *menu.Menu {
	var items []menu.Item
	items = append(items, app.menuItemEscape())

	presetRoutes := make([]route, len(app.cfg.Presets.Multi))
	for i := range app.cfg.Presets.Multi {
		presetRoutes[i] = routeMultiPlayerDirectIPHostPresetN + route(strconv.Itoa(i))
	}

	for i := range app.cfg.Presets.Multi {
		name := Tf(KeyMenuHostDirectIPPreset,
			i+1,
			app.cfg.Presets.Multi[i].Name,
			app.cfg.Presets.Multi[i].String(),
		)
		items = append(items, menu.NewCommand(&app.screenIDNext, presetRoutes[i], name, T(KeyMenuHostDirectIPPresetDesc)))
	}

	items = append(items, menu.NewCommand(&app.screenIDNext, routeMultiPlayerDirectIPHostCustomSetup,
		T(KeyMenuHostDirectIPCustom), T(KeyMenuHostDirectIPCustomDesc)))
	items = append(items, menu.NewCommand(&app.screenIDNext, routeMultiPlayerPresetEditMenu,
		T(KeyMenuHostDirectIPEditPresets), Tf(KeyMenuHostDirectIPEditPresetsDesc, setup.MultiPlayerPresetCount)))

	items = append(items, app.menuItemBack())

	return menu.New(T(KeyMenuHostDirectIPTitle), app.menuStopper(ctx), items...)
}
