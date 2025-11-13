// Copyright (c) 2024, 2025 by Marko Gaćeša

package app

import (
	"fmt"
	"math/rand/v2"
	"strconv"

	"github.com/marko-gacesa/gamatet/game/setup"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) menuMultiPlayerEditPresets(ctx screen.Context) *menu.Menu {
	var items []menu.Item
	items = append(items, app.menuItemEscape())

	presetRoutes := make([]route, len(app.cfg.Presets.Multi))
	for i := range app.cfg.Presets.Multi {
		presetRoutes[i] = routeMultiPlayerPresetEditN + route(strconv.Itoa(i))
	}

	for i := range app.cfg.Presets.Multi {
		name := fmt.Sprintf("Edit Preset %d: %s [%s]",
			i+1,
			app.cfg.Presets.Multi[i].Name,
			app.cfg.Presets.Multi[i].String(),
		)
		items = append(items, menu.NewCommand(&app.screenIDNext, presetRoutes[i], name, ""))
	}
	items = append(items, app.menuItemBack())

	return menu.New("Multi Player: Edit presets", app.menuStopper(ctx), items...)
}

func (app *App) menuMultiPlayerSetup(ctx screen.Context, maxPlayers byte, presetIdx int, nextRoute route) *menu.Menu {
	app.resultSetup = nil // Clear result of this input

	if presetIdx >= len(app.cfg.Presets.Multi) {
		return app.menuErrorText(ctx, "Preset index out of range")
	}

	var s setup.Setup
	if presetIdx >= 0 {
		s = app.cfg.Presets.Multi[presetIdx]
		if s.Name == "" {
			s.Name = "Game"
		}
	} else {
		s = app.cfg.Presets.MultiCustom
	}

	if !s.MiscOptions.CustomSeed {
		s.MiscOptions.Seed = rand.Int64()
	}

	sections := newSetupSections()
	sections.refresh(&s)

	items := make([]menu.Item, 0, 32)
	items = append(items, app.menuItemEscape())
	if presetIdx >= 0 {
		items = append(items, menu.NewText(&s.Name, setup.MaxLenName, setup.MaxLenName,
			"Game name", ""))
	}
	items = append(items, setupMultiPlayer(&s, sections)...)
	items = append(items, setupResultMulti(&s, &app.resultSetup, maxPlayers, presetIdx >= 0)...)
	items = append(items, app.menuItemBack())

	return menu.New("Multi Player: Setup", func(m *menu.Menu) {
		if app.screenIDNext != "" {
			ctx.Stop()
			return
		}

		if app.resultSetup != nil {
			if app.resultSetup.SanitizeMulti() {
				app.logger.Warn("sanitation required after setup")
			}
			if presetIdx >= 0 {
				app.cfg.Presets.Multi[presetIdx] = *app.resultSetup
			} else {
				app.cfg.Presets.MultiCustom = *app.resultSetup
			}

			app.screenIDNext = nextRoute
			ctx.Stop()
			return
		}

		sections.refresh(&s)
	}, items...)
}
