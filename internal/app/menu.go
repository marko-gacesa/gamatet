// Copyright (c) 2024, 2025 by Marko Gaćeša

package app

import (
	"gamatet/internal/values"
	"gamatet/logic/menu"
	"gamatet/logic/screen"
)

const (
	itemTextPrefixBack    = "← "
	itemTextPrefixForward = "→ "
)

func (app *App) menuStopper(ctx screen.Context) func(*menu.Menu) {
	return func(*menu.Menu) {
		if app.screenIDNext != "" {
			ctx.Stop()
		}
	}
}

func (app *App) menuItemEscape() *menu.Hidden[route] {
	return menu.NewHidden(menu.InputEscape, &app.screenIDNext, routeBack)
}

func (app *App) menuItemBack() *menu.Command[route] {
	return menu.NewCommand(&app.screenIDNext, routeBack, itemTextPrefixBack+"Back", "")
}

func (app *App) menuMain(ctx screen.Context) *menu.Menu {
	return menu.New(values.ProgramName, app.menuStopper(ctx), []menu.Item{
		app.menuItemEscape(),
		menu.NewCommand(&app.screenIDNext, routeMenuSinglePlayer, "Single player", ""),
		menu.NewCommand(&app.screenIDNext, routeMenuLocalCreate, "Multiplayer local", "Multiplayer game on this machine"),
		menu.NewCommand(&app.screenIDNext, routeMenuLANMain, "Multiplayer LAN game", "Multiplayer game on the local area network"),
		menu.NewCommand(&app.screenIDNext, routeQuit, "Exit", "Exit application"),
	}...)
}

func (app *App) menuSinglePlayer(ctx screen.Context) *menu.Menu {
	return menu.New("Single Player", app.menuStopper(ctx), []menu.Item{
		app.menuItemEscape(),
		menu.NewCommand(&app.screenIDNext, routeGameSinglePlayNow, "Play Now!", "Start a classic game"),
		menu.NewCommand(&app.screenIDNext, routeGameDoublePlayNow, "Play Double Now!", "Start a double game"),
		app.menuItemBack(),
	}...)
}
