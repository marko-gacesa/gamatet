// Copyright (c) 2024, 2025 by Marko Gaćeša

package app

import (
	"strconv"

	"github.com/marko-gacesa/gamatet/game/setup"
	. "github.com/marko-gacesa/gamatet/internal/i18n"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) menuMultiPlayerLANHostMenu(ctx screen.Context) *menu.Menu {
	var items []menu.Item
	items = append(items, app.menuItemEscape())

	presetRoutes := make([]route, len(app.cfg.Presets.Multi))
	for i := range app.cfg.Presets.Multi {
		presetRoutes[i] = routeMultiPlayerLANHostPresetN + route(strconv.Itoa(i))
	}

	for i := range app.cfg.Presets.Multi {
		name := Tf(KeyMenuHostLANPreset,
			i+1,
			app.cfg.Presets.Multi[i].Name,
			app.cfg.Presets.Multi[i].String(),
		)
		items = append(items, menu.NewCommand(&app.screenIDNext, presetRoutes[i], name, T(KeyMenuHostLANPresetDesc)))
	}

	items = append(items, menu.NewCommand(&app.screenIDNext, routeMultiPayerLANHostCustomSetup,
		T(KeyMenuHostLANCustom), T(KeyMenuHostLANCustomDesc)))
	items = append(items, menu.NewCommand(&app.screenIDNext, routeMultiPlayerPresetEditMenu,
		T(KeyMenuHostLANEditPresets), Tf(KeyMenuHostLANEditPresetsDesc, setup.MultiPlayerPresetCount)))

	items = append(items, app.menuItemBack())

	return menu.New(T(KeyMenuHostLANTitle), app.menuStopper(ctx), items...)
}
