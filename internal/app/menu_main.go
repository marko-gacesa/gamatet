// Copyright (c) 2024, 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	. "github.com/marko-gacesa/gamatet/internal/i18n"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) menuMain(ctx screen.Context) *menu.Menu {
	return menu.New("", app.menuStopper(ctx), []menu.Item{
		app.menuItemEscape(),
		menu.NewCommand(&app.screenIDNext, routeSinglePlayerMenu, T(KeyMenuMainSinglePlayer), T(KeyMenuMainSinglePlayerDesc)),
		menu.NewCommand(&app.screenIDNext, routeMultiPlayerLocalMenu, T(KeyMenuMainMultiplayerLocal), T(KeyMenuMainMultiplayerLocalDesc)),
		menu.NewCommand(&app.screenIDNext, routeMultiPlayerLANMenu, T(KeyMenuMainMultiplayerLAN), T(KeyMenuMainMultiplayerLANDesc)),
		menu.NewCommand(&app.screenIDNext, routeMultiPlayerDirectIPMenu, T(KeyMenuMainMultiPlayerDirectIP), T(KeyMenuMainMultiPlayerDirectIPDesc)),
		menu.NewCommand(&app.screenIDNext, routeConfigMenu, T(KeyMenuMainConfigure), T(KeyMenuMainConfigureDesc)),
		menu.NewCommand(&app.screenIDNext, routeAboutMenu, T(KeyMenuMainAbout), T(KeyMenuMainAboutDesc)),
		menu.NewCommand(&app.screenIDNext, routeQuit, T(KeyMenuMainExit), T(KeyMenuMainExitDesc)),
	}...)
}
