// Copyright (c) 2024, 2025 by Marko Gaćeša

package app

import (
	"github.com/marko-gacesa/gamatet/internal/values"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) menuMultiPlayerLANMain(ctx screen.Context) *menu.Menu {
	return menu.New(values.ProgramName, app.menuStopper(ctx), []menu.Item{
		app.menuItemEscape(),
		menu.NewCommand(&app.screenIDNext, routeMultiPlayerLANHostMenu, itemTextPrefixForward+"Host LAN game", ""),
		menu.NewCommand(&app.screenIDNext, routeMultiPlayerLANJoinListen, itemTextPrefixForward+"Join LAN game", ""),
		app.menuItemBack(),
	}...)
}
