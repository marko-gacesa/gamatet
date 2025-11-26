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

func (app *App) menuSinglePlayer(ctx screen.Context) *menu.Menu {
	var items []menu.Item
	items = append(items, app.menuItemEscape())

	presetRoutes := make([]route, len(app.cfg.Presets.Single))
	for i := range app.cfg.Presets.Single {
		presetRoutes[i] = routeSinglePlayerPresetGameN + route(strconv.Itoa(i))
	}

	for i := range app.cfg.Presets.Single {
		name := Tf(KeyMenuSingleStartPreset,
			i+1,
			app.cfg.Presets.Single[i].Name,
			app.cfg.Presets.Single[i].String(),
		)
		items = append(items, menu.NewCommand(&app.screenIDNext, presetRoutes[i], name, T(KeyMenuSingleStartPresetDesc)))
	}
	items = append(items, menu.NewCommand(&app.screenIDNext, routeSinglePlayerCustomSetup,
		T(KeyMenuSingleCustom), T(KeyMenuSingleCustomDesc)))
	items = append(items, menu.NewCommand(&app.screenIDNext, routeSinglePlayerPresetEditMenu,
		T(KeyMenuSingleEditPresets), Tf(KeyMenuSingleEditPresetsDesc, setup.SinglePlayerPresetCount)))
	items = append(items, app.menuItemBack())

	return menu.New(T(KeyMenuSingleTitle), app.menuStopper(ctx), items...)
}

func (app *App) menuSingleEditPresets(ctx screen.Context) *menu.Menu {
	var items []menu.Item
	items = append(items, app.menuItemEscape())

	presetRoutes := make([]route, len(app.cfg.Presets.Single))
	for i := range app.cfg.Presets.Single {
		presetRoutes[i] = routeSinglePlayerPresetEditN + route(strconv.Itoa(i))
	}

	for i := range app.cfg.Presets.Single {
		name := Tf(KeyMenuSingleEditPreset,
			i+1,
			app.cfg.Presets.Single[i].Name,
			app.cfg.Presets.Single[i].String(),
		)
		items = append(items, menu.NewCommand(&app.screenIDNext, presetRoutes[i], name, T(KeyMenuSingleEditPresetDesc)))
	}
	items = append(items, app.menuItemBack())

	return menu.New(T(KeyMenuSingleEditPresetsTitle), app.menuStopper(ctx), items...)
}

func (app *App) menuSinglePlayerSetup(ctx screen.Context, presetIdx int, nextRoute route) *menu.Menu {
	app.resultSetup = nil // Clear result of this input

	if presetIdx >= len(app.cfg.Presets.Single) {
		return app.menuErrorText(ctx, T(KeyErrorPresetIndexOutOfRange))
	}

	var s setup.Setup
	if presetIdx >= 0 {
		s = app.cfg.Presets.Single[presetIdx]
		if s.Name == "" {
			s.Name = T(KeySetupDefaultGameName)
		}
	} else {
		s = app.cfg.Presets.SingleCustom
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
	items = append(items, setupSingle(&s, sections)...)
	items = append(items, setupResultSingle(&s, &app.resultSetup, presetIdx >= 0)...)
	items = append(items, app.menuItemBack())

	return menu.New(T(KeySetupTitleSingle), func(m *menu.Menu) {
		app.menuStopper(ctx)(m)
		sections.refresh(&s)

		if app.resultSetup != nil {
			if app.resultSetup.SanitizeSingle() {
				app.logger.Warn("sanitation required after setup")
			}
			if presetIdx >= 0 {
				app.cfg.Presets.Single[presetIdx] = *app.resultSetup
			} else {
				app.cfg.Presets.SingleCustom = *app.resultSetup
			}
			app.saveConfig()

			app.screenIDNext = nextRoute
			ctx.Stop()
			return
		}
	}, items...)
}
