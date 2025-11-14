// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"fmt"

	"github.com/marko-gacesa/gamatet/game/setup"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) menuConfigLocalPlayer(ctx screen.Context, idx int) *menu.Menu {
	if idx < 0 || idx >= setup.MaxLocalPlayers {
		return app.menuErrorText(ctx, "player index out of range")
	}

	dirMap := map[bool]string{
		false: "Counter clockwise",
		true:  "Clockwise",
	}

	items := make([]menu.Item, 0, 2)
	items = append(items,
		menu.NewText(&app.cfg.LocalPlayers.Infos[idx].Name, setup.MaxLenName, setup.MaxLenName, "Name", ""),
		menu.NewEnum(&app.cfg.LocalPlayers.Infos[idx].GameConfig.RotationDirectionCW, []bool{false, true}, dirMap, "Rotation Direction", ""),
		app.menuItemEscape(),
		app.menuItemBack(),
	)
	return menu.New(fmt.Sprintf("Player: %d", idx+1), app.menuStopper(ctx), items...)
}
