// Copyright (c) 2024,2025 by Marko Gaćeša

package app

import (
	"gamatet/logic/menu"
	"gamatet/logic/screen"
	"gamatet/logic/values"
)

func (app *App) menuStopper(ctx screen.Context) func() {
	return func() {
		if app.screenIDNext != "" {
			ctx.Stop()
		}
	}
}

func (app *App) menuMain(ctx screen.Context) *menu.Menu {
	return menu.New(values.ProgramName, app.menuStopper(ctx), []menu.Item{
		menu.NewCancel(&app.screenIDNext, routeBack),
		menu.NewCommand(&app.screenIDNext, routeMenuSinglePlayer, "Single player", "Start demo fields"),
		menu.NewCommand(&app.screenIDNext, routeTestField, "Fields demo", "Start demo fields"),
		menu.NewCommand(&app.screenIDNext, routeTestBlocks, "Blocks", "Blocks demo"),
		menu.NewCommand(&app.screenIDNext, routeQuit, "Exit", "Exit application"),
	}...)
}

func (app *App) menuSinglePlayer(ctx screen.Context) *menu.Menu {
	return menu.New("Single Player", app.menuStopper(ctx), []menu.Item{
		menu.NewCancel(&app.screenIDNext, routeBack),
		menu.NewCommand(&app.screenIDNext, routeGameSinglePlayNow, "Play Now!", "Start a classic game"),
		menu.NewCommand(&app.screenIDNext, routeGameDoublePlayNow, "Play Double Now!", "Start a double game"),
		menu.NewCommand(&app.screenIDNext, routeBack, "Back", "Back to main menu"),
	}...)
}
