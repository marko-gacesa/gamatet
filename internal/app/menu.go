// Copyright (c) 2024 by Marko Gaćeša

package app

import (
	"context"
	"gamatet/logic/menu"
)

func (app *App) routeTo(route route, cancelFn context.CancelFunc) func(*menu.Menu, *menu.Command) {
	return func(m *menu.Menu, cmd *menu.Command) {
		app.screenIDNext = route
		cancelFn()
		cmd.ClearFunction()
	}
}

func (app *App) menuMain(ctx context.Context) (*menu.Menu, context.Context) {
	ctx, cancelCtx := context.WithCancel(ctx)
	return menu.New("Gamatet", cancelCtx, []menu.Item{
		menu.NewCommand("Single player", "Start demo fields", app.routeTo(routeMenuSinglePlayer, cancelCtx)),
		menu.NewCommand("Fields demo", "Start demo fields", app.routeTo(routeTestField, cancelCtx)),
		menu.NewCommand("Blocks", "Blocks demo", app.routeTo(routeTestBlocks, cancelCtx)),
		menu.NewCommand("Exit", "Exit application", app.routeTo(routeQuit, cancelCtx)),
	}...), ctx
}

func (app *App) menuSinglePlayer(ctx context.Context) (*menu.Menu, context.Context) {
	ctx, cancelCtx := context.WithCancel(ctx)
	return menu.New("Single Player", cancelCtx, []menu.Item{
		menu.NewCommand("Play Now!", "Start a classic game", app.routeTo(routeGameSinglePlayNow, cancelCtx)),
		menu.NewCommand("Play Double Now!", "Start a double game", app.routeTo(routeGameDoublePlayNow, cancelCtx)),
		menu.NewCommand("Back", "Back to main menu", app.routeTo(routeBack, cancelCtx)),
	}...), ctx
}
