// Copyright (c) 2024 by Marko Gaćeša

package app

import (
	"gamatet/logic/menu"
)

func (app *App) routeTo(route route) func(*menu.Menu, *menu.Command) {
	return func(m *menu.Menu, cmd *menu.Command) {
		app.screenIDNext = route
		m.Finish()
		cmd.ClearFunction()
	}
}

func (app *App) menuMain() *menu.Menu {
	return menu.New("Gamatet", []menu.Item{
		menu.NewCommand("Fields demo", "Start demo fields", app.routeTo(routeTestField)),
		menu.NewCommand("Blocks", "Blocks demo", app.routeTo(routeTestBlocks)),
		menu.NewCommand("Exit", "Exit application", app.routeTo(routeQuit)),
	}...)
}
