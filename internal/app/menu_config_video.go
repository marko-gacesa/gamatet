// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) menuConfigVideo(ctx screen.Context) *menu.Menu {
	return menu.New("Video Options", app.menuStopper(ctx), []menu.Item{
		app.menuItemEscape(),
		app.menuItemBack(),
	}...)
}
