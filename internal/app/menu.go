// Copyright (c) 2024, 2025 by Marko Gaćeša

package app

import (
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
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
