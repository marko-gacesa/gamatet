// Copyright (c) 2024, 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	. "github.com/marko-gacesa/gamatet/internal/i18n"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) menuMultiPlayerLANMain(ctx screen.Context) *menu.Menu {
	return menu.New(T(KeyMenuLANTitle), app.menuStopper(ctx), []menu.Item{
		app.menuItemEscape(),
		menu.NewCommand(&app.screenIDNext, routeMultiPlayerLANHostMenu, T(KeyMenuLANHost), T(KeyMenuLANHostDesc)),
		menu.NewCommand(&app.screenIDNext, routeMultiPlayerLANJoinListen, T(KeyMenuLANJoin), T(KeyMenuLANJoinDesc)),
		app.menuItemBack(),
	}...)
}
