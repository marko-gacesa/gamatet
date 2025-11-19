// Copyright (c) 2024, 2025 by Marko Gaćeša

package app

import (
	. "github.com/marko-gacesa/gamatet/internal/i18n"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) menuMain(ctx screen.Context) *menu.Menu {
	return menu.New(T(KeyProgramName), app.menuStopper(ctx), []menu.Item{
		app.menuItemEscape(),
		menu.NewCommand(&app.screenIDNext, routeSinglePlayerMenu, "Single player", ""),
		menu.NewCommand(&app.screenIDNext, routeMultiPlayerLocalMenu, "Multiplayer local", "Multiplayer game on this machine"),
		menu.NewCommand(&app.screenIDNext, routeMultiPlayerLANMenu, "Multiplayer LAN game", "Multiplayer game on the local area network"),
		menu.NewCommand(&app.screenIDNext, routeConfigMenu, "Configure", ""),
		menu.NewCommand(&app.screenIDNext, routeAboutMenu, "About", ""),
		menu.NewCommand(&app.screenIDNext, routeQuit, "Exit", "Exit game"),
	}...)
}
