// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"fmt"
	"time"

	"github.com/marko-gacesa/gamatet/internal/values"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) menuAbout(ctx screen.Context) *menu.Menu {
	t, _ := time.Parse(time.RFC3339Nano, values.BuildTime)

	var createdBy string

	if y := t.Year(); y < 2000 {
		createdBy = "by Marko Gaćeša"
	} else {
		createdBy = fmt.Sprintf("Copyright (c) %d Marko Gaćeša", t.Year())
	}

	return menu.New("About", app.menuStopper(ctx), []menu.Item{
		menu.NewStatic(createdBy, "", nil),
		app.menuItemEscape(),
		app.menuItemBack(),
	}...)
}
