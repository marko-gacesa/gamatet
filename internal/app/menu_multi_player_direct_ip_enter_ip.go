// Copyright (c) 2024, 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	"math"
	"math/rand/v2"

	"github.com/marko-gacesa/gamatet/game/setup"
	. "github.com/marko-gacesa/gamatet/internal/i18n"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
	"github.com/marko-gacesa/udpstar/udpstar/message"
)

func (app *App) menuMultiPlayerDirectIPHostEnterIP(ctx screen.Context) *menu.Menu {
	app.resultToken = message.Token(rand.Uint32())
	return app.menuEnterIP(ctx, routeMultiPlayerDirectIPHostLobby)
}

func (app *App) menuMultiPlayerDirectIPJoinEnterIP(ctx screen.Context) *menu.Menu {
	return app.menuEnterIP(ctx, routeMultiPlayerDirectIPJoinLobby)
}

func (app *App) menuEnterIP(ctx screen.Context, r route) *menu.Menu {
	var items []menu.Item
	items = append(items, app.menuItemEscape())

	app.resultToken = message.RandomToken()%99 + 1

	items = append(items,
		menu.NewIP(&app.cfg.Network.DirectIPAddress,
			T(KeyMenuDirectIPEnterIP), T(KeyMenuDirectIPEnterIPDesc)),
		menu.NewNumber(&app.resultToken, 1, math.MaxUint32,
			T(KeyMenuDirectIPEnterToken), T(KeyMenuDirectIPEnterTokenDesc)),
		menu.NewText(&app.cfg.LocalPlayers.Infos[0].Name, setup.MaxLenName, setup.MaxLenName,
			T(KeyConfigPlayerName), T(KeyConfigPlayerNameDesc)),
		menu.NewCommand(&app.screenIDNext, r,
			T(KeyMenuDirectIPProceedToLobby), T(KeyMenuDirectIPProceedToLobbyDesc),
			menu.WithVisible(func() bool {
				return app.cfg.Network.DirectIPAddress != ""
			}),
		),
		menu.NewStatic(
			T(KeyMenuDirectIPInvalidIP), T(KeyMenuDirectIPInvalidIPDesc),
			nil,
			menu.WithVisible(func() bool {
				return app.cfg.Network.DirectIPAddress == ""
			}),
			menu.WithDisabled(func() bool { return true }),
		),
	)

	items = append(items, app.menuItemBack())

	return menu.New(T(KeyMenuHostDirectIPTitle), app.configStopper(ctx), items...)
}
