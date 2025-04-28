// Copyright (c) 2024, 2025 by Marko Gaćeša

package app

import (
	"gamatet/game/setup"
	"gamatet/internal/values"
	"gamatet/logic/menu"
	"gamatet/logic/screen"
	"math/rand/v2"
	"os"
)

func (app *App) menuLANServerCreate(ctx screen.Context) *menu.Menu {
	app.resultSetup = nil // Clear result of this input

	s := app.cfg.LANGameDefaults
	if !s.MiscOptions.CustomSeed {
		s.MiscOptions.Seed = rand.Int64()
	}

	if s.Name == "" {
		hostname, _ := os.Hostname()
		if hostname == "" {
			hostname = "LAN"
		}
		s.Name = hostname + " game"
	}

	sections := newSetupSections()
	sections.refresh(&s)

	items := make([]menu.Item, 0, 32)
	items = append(items, app.menuItemEscape())
	items = append(items, setupMulti(&s, sections)...)
	items = append(items, setupResult(&s, &app.resultSetup, setup.MaxPlayers)...)
	items = append(items, app.menuItemBack())

	return menu.New(values.ProgramName, func(m *menu.Menu) {
		if app.screenIDNext != "" {
			ctx.Stop()
			return
		}

		if app.resultSetup != nil {
			app.screenIDNext = routeMenuLANServerLobby
			ctx.Stop()
			return
		}

		sections.refresh(&s)
	}, items...)
}
