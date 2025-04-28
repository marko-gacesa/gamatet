// Copyright (c) 2024, 2025 by Marko Gaćeša

package app

import (
	"gamatet/internal/values"
	"gamatet/logic/menu"
	"gamatet/logic/screen"
)

func (app *App) menuShowError(m *menu.Menu, err error) {
	m.SetItems([]menu.Item{
		app.menuItemEscape(),
		menu.NewStatic("Error: "+err.Error(), "", nil),
		app.menuItemBack(),
	}...)
}

func (app *App) menuError(ctx screen.Context, err error) *menu.Menu {
	return app.menuErrorText(ctx, err.Error())
}

func (app *App) menuErrorText(ctx screen.Context, text string) *menu.Menu {
	return menu.New(values.ProgramName, app.menuStopper(ctx), []menu.Item{
		app.menuItemEscape(),
		menu.NewStatic("Error: "+text, "", nil),
		app.menuItemBack(),
	}...)
}
