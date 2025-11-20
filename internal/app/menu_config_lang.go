// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"slices"

	"github.com/marko-gacesa/gamatet/internal/config"
	. "github.com/marko-gacesa/gamatet/internal/i18n"
	"github.com/marko-gacesa/gamatet/logic/lang"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) menuConfigLanguage(ctx screen.Context) *menu.Menu {
	supportedMap := lang.StrInAll(KeyLanguageName)
	langs := make([]lang.Lang, 0, len(supportedMap))
	for l := range supportedMap {
		langs = append(langs, l)
	}
	slices.Sort(langs)

	items := make([]menu.Item, len(supportedMap)+1)
	for i, l := range langs {
		items[i] = menu.NewCommand(&app.cfg.Locale.Language, string(l), supportedMap[l], "")
	}
	items[len(supportedMap)] = app.menuItemEscape()

	return menu.New(T(KeyConfigLanguage), func(*menu.Menu) {
		if app.screenIDNext == "" {
			lang.Set(lang.Lang(app.cfg.Locale.Language))
			_ = config.Save(app.logger, app.cfgPath, &app.cfg)
		}
		app.screenIDNext = routeBack
		ctx.Stop()
	}, items...)
}
