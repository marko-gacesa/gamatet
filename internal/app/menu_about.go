// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"fmt"

	"github.com/marko-gacesa/gamatet/internal/values"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) menuAbout(ctx screen.Context) *menu.Menu {
	var createdBy string

	if t := values.BuildTime; len(t) >= 4 {
		createdBy = fmt.Sprintf("Copyright (c) %s Marko Gaćeša", t[:4])
	} else {
		createdBy = "by Marko Gaćeša"
	}

	return menu.New("About", app.menuStopper(ctx), []menu.Item{
		menu.NewStatic(createdBy, "", nil),
		app.menuItemEscape(),
		app.menuItemBack(),
	}...)
}
