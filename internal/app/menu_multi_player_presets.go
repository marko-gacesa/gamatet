// Copyright (c) 2024, 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	"math/rand/v2"
	"strconv"

	"github.com/marko-gacesa/gamatet/game/setup"
	. "github.com/marko-gacesa/gamatet/internal/i18n"
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
		name := Tf(KeyMenuMultiEditPreset,
			i+1,
			app.cfg.Presets.Multi[i].Name,
			app.cfg.Presets.Multi[i].String(),
		)
		items = append(items, menu.NewCommand(&app.screenIDNext, presetRoutes[i], name, T(KeyMenuMultiEditPresetDesc)))
	}
	items = append(items, app.menuItemBack())

	return menu.New(T(KeyMenuMultiEditPresetsTitle), app.menuStopper(ctx), items...)
}

func (app *App) menuMultiPlayerSetup(ctx screen.Context, maxPlayers byte, presetIdx int, nextRoute route) *menu.Menu {
	app.resultSetup = nil // Clear result of this input

	if presetIdx >= len(app.cfg.Presets.Multi) {
		return app.menuError(ctx, errorPresetIndexOutOfRange)
	}

	var s setup.Setup
	if presetIdx >= 0 {
		s = app.cfg.Presets.Multi[presetIdx]
		if s.Name == "" {
			s.Name = T(KeySetupDefaultGameName)
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
			T(KeySetupGameName), T(KeySetupGameNameDesc)))
	}
	items = append(items, setupMultiPlayer(&s, sections)...)
	items = append(items, setupResultMulti(&s, &app.resultSetup, maxPlayers, presetIdx >= 0)...)
	items = append(items, app.menuItemBack())

	return menu.New(T(KeySetupTitleMulti), func(m *menu.Menu) {
		app.menuStopper(ctx)(m)
		sections.refresh(&s)

		if app.resultSetup != nil {
			if app.resultSetup.SanitizeMulti() {
				app.logger.Warn("sanitation required after setup")
			}
			if presetIdx >= 0 {
				app.cfg.Presets.Multi[presetIdx] = *app.resultSetup
			} else {
				app.cfg.Presets.MultiCustom = *app.resultSetup
			}
			app.saveConfig()

			app.screenIDNext = nextRoute
			ctx.Stop()
			return
		}
	}, items...)
}
