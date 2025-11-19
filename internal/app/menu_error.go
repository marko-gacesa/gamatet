// Copyright (c) 2024, 2025 by Marko Gaćeša

package app

import (
	. "github.com/marko-gacesa/gamatet/internal/i18n"
	"github.com/marko-gacesa/gamatet/internal/values"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) menuShowError(m *menu.Menu, err error) {
	m.SetItems([]menu.Item{
		app.menuItemEscape(),
		menu.NewStatic(T(KeyError)+": "+err.Error(), "", nil),
		app.menuItemBack(),
	}...)
}

func (app *App) menuError(ctx screen.Context, err error) *menu.Menu {
	return app.menuErrorText(ctx, err.Error())
}

func (app *App) menuErrorText(ctx screen.Context, text string) *menu.Menu {
	return menu.New(values.ProgramName, app.menuStopper(ctx), []menu.Item{
		app.menuItemEscape(),
		menu.NewStatic(T(KeyError)+": "+text, "", nil),
		app.menuItemBack(),
	}...)
}
