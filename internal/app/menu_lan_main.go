// Copyright (c) 2024, 2025 by Marko Gaćeša

package app

import (
	"gamatet/internal/values"
	"gamatet/logic/menu"
	"gamatet/logic/screen"
)

func (app *App) menuLANMain(ctx screen.Context) *menu.Menu {
	return menu.New(values.ProgramName, app.menuStopper(ctx), []menu.Item{
		app.menuItemEscape(),
		menu.NewCommand(&app.screenIDNext, routeMenuLANServerCreate, itemTextPrefixForward+"Host LAN game", ""),
		menu.NewCommand(&app.screenIDNext, routeMenuLANClientJoin, itemTextPrefixForward+"Join LAN game", ""),
		app.menuItemBack(),
	}...)
}
