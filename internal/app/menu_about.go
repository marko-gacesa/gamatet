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

	if values.BuildTime != "" {
		items = append(items, menu.NewStatic(T(KeyMenuAboutBuildTime)+": "+values.BuildTime, "", nil))
	}
	if values.GitSHA != "" {
		items = append(items, menu.NewStatic(T(KeyMenuAboutGitSHA)+": "+values.GitSHA, "", nil))
	}
	if values.VersionTag != "" {
		items = append(items, menu.NewStatic(T(KeyMenuAboutVersion)+": "+values.VersionTag, "", nil))
	}

	items = append(items, menu.NewStatic(fmt.Sprintf("Copyright (c) %s Marko Gaćeša", values.ProgramDate), "", nil))

	items = append(items, app.menuItemEscape())
	items = append(items, app.menuItemBack())

	return menu.New(T(KeyMenuAboutTitle), app.menuStopper(ctx), items...)
}
