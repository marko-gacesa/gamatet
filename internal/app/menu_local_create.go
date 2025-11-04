// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"gamatet/game/setup"
	"gamatet/internal/values"
	"gamatet/logic/menu"
	"gamatet/logic/screen"
	"math/rand/v2"
)

func (app *App) menuLocalCreateGame(ctx screen.Context, idx int) *menu.Menu {
	app.resultSetup = nil // Clear result of this input

	s := app.cfg.Presets.Multi[idx]
	if !s.MiscOptions.CustomSeed {
		s.MiscOptions.Seed = rand.Int64()
	}

	if s.Name == "" {
		s.Name = "Game"
	}

	sections := newSetupSections()
	sections.refresh(&s)

	items := make([]menu.Item, 0, 32)
	items = append(items, app.menuItemEscape())
	items = append(items, setupMulti(&s, sections)...)
	items = append(items, setupResult(&s, &app.resultSetup, setup.MaxLocalPlayers)...)
	items = append(items, app.menuItemBack())

	return menu.New(values.ProgramName, func(m *menu.Menu) {
		if app.screenIDNext != "" {
			ctx.Stop()
			return
		}

		if app.resultSetup != nil {
			app.screenIDNext = routeGame
			ctx.Stop()
			return
		}

		sections.refresh(&s)
	}, items...)
}
