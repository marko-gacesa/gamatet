// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"fmt"

	"github.com/marko-gacesa/gamatet/game/setup"
	"github.com/marko-gacesa/gamatet/internal/config/key"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) setupLocalPlayer(ctx screen.Context, idx int) *menu.Menu {
	if idx < 0 || idx >= setup.MaxLocalPlayers {
		return app.menuErrorText(ctx, "player index out of range")
	}

	dirMap := map[bool]string{
		false: "Counter clockwise",
		true:  "Clockwise",
	}

	stringFn := func(k key.Key) string { return key.Map[k] }
	convertFn := func(b byte) key.Key { return key.Key(b) }

	info := &app.cfg.LocalPlayers.Infos[idx]

	section := newSetupKeySection()
	section.refresh(&info.Input)
	showKeysFn := func() bool {
		return section.showKeys
	}

	items := make([]menu.Item, 0, 10)
	items = append(items,
		menu.NewText(&info.Name, setup.MaxLenName, setup.MaxLenName, "Name", ""),
		menu.NewEnum(&section.showKeys, []bool{false, true}, section.showKeysMap,
			"Control", ""),
		menu.NewKey(&info.Input.Left,
			"\tMove piece left key", "", stringFn, convertFn, menu.WithVisible(showKeysFn)),
		menu.NewKey(&info.Input.Right,
			"\tMove piece right key", "", stringFn, convertFn, menu.WithVisible(showKeysFn)),
		menu.NewKey(&info.Input.Activate,
			"\tActivate piece (rotate/transform) key", "", stringFn, convertFn, menu.WithVisible(showKeysFn)),
		menu.NewKey(&info.Input.Drop,
			"\tDrop piece down key", "", stringFn, convertFn, menu.WithVisible(showKeysFn)),
		menu.NewEnum(&info.GameConfig.RotationDirectionCW,
			[]bool{false, true}, dirMap, "Rotation Direction", ""),
		app.menuItemEscape(),
		app.menuItemBack(),
	)
	return menu.New(fmt.Sprintf("Player: %d", idx+1), func(m *menu.Menu) {
		if app.screenIDNext != "" {
			ctx.Stop()
			return
		}

		section.refresh(&info.Input)
	}, items...)
}
