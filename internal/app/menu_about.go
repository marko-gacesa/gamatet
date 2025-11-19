// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"fmt"

	. "github.com/marko-gacesa/gamatet/internal/i18n"
	"github.com/marko-gacesa/gamatet/internal/values"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) menuAbout(ctx screen.Context) *menu.Menu {
	var items []menu.Item

	items = append(items, menu.NewStatic(T(KeyMenuAboutAuthor)+fmt.Sprintf(": %s © Marko Gaćeša", values.ProgramDate), "", nil))
	if values.VersionTag != "" {
		items = append(items,
			menu.NewStatic(T(KeyMenuAboutVersion)+": "+values.VersionTag,
				T(KeyMenuAboutGitSHA)+": "+values.GitSHA+". "+T(KeyMenuAboutBuildTime)+": "+values.BuildTime,
				nil))
	}

	items = append(items, app.menuItemEscape())
	items = append(items, app.menuItemBack())

	return menu.New(T(KeyMenuAboutTitle), app.menuStopper(ctx), items...)
}
