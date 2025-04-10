// Copyright (c) 2024,2025 by Marko Gaćeša

package app

import (
	"gamatet/logic/menu"
	"gamatet/logic/values"
)

func (app *App) menuStopper(stopFn func()) func() {
	return func() {
		if app.screenIDNext != "" {
			stopFn()
		}
	}
}

func (app *App) menuMain(stopFn func()) *menu.Menu {
	return menu.New(values.ProgramName, app.menuStopper(stopFn), []menu.Item{
		menu.NewCancel(&app.screenIDNext, routeBack),
		menu.NewCommand(&app.screenIDNext, routeMenuSinglePlayer, "Single player", "Start demo fields"),
		menu.NewCommand(&app.screenIDNext, routeTestField, "Fields demo", "Start demo fields"),
		menu.NewCommand(&app.screenIDNext, routeTestBlocks, "Blocks", "Blocks demo"),
		menu.NewCommand(&app.screenIDNext, routeQuit, "Exit", "Exit application"),
	}...)
}

func (app *App) menuSinglePlayer(stopFn func()) *menu.Menu {
	return menu.New("Single Player", app.menuStopper(stopFn), []menu.Item{
		menu.NewCancel(&app.screenIDNext, routeBack),
		menu.NewCommand(&app.screenIDNext, routeGameSinglePlayNow, "Play Now!", "Start a classic game"),
		menu.NewCommand(&app.screenIDNext, routeGameDoublePlayNow, "Play Double Now!", "Start a double game"),
		menu.NewCommand(&app.screenIDNext, routeBack, "Back", "Back to main menu"),
	}...)
}
