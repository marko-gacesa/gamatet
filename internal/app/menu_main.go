// Copyright (c) 2024, 2025 by Marko Gaćeša

package app

import (
	"gamatet/internal/values"
	"gamatet/logic/menu"
	"gamatet/logic/screen"
)

func (app *App) menuMain(ctx screen.Context) *menu.Menu {
	return menu.New(values.ProgramName, app.menuStopper(ctx), []menu.Item{
		app.menuItemEscape(),
		menu.NewCommand(&app.screenIDNext, routeMenuSinglePlayer, "Single player", ""),
		menu.NewCommand(&app.screenIDNext, routeMenuLocalCreate, "Multiplayer local", "Multiplayer game on this machine"),
		menu.NewCommand(&app.screenIDNext, routeMenuLANMain, "Multiplayer LAN game", "Multiplayer game on the local area network"),
		menu.NewCommand(&app.screenIDNext, routeQuit, "Exit", "Exit application"),
	}...)
}
