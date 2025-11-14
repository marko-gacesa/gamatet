// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"github.com/marko-gacesa/gamatet/internal/config"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) menuConfigVideo(ctx screen.Context) *menu.Menu {
	return menu.New("Video Options", app.menuStopper(ctx), []menu.Item{
		menu.NewBool(&app.cfg.Video.Fullscreen, "Fullscreen", ""),
		menu.NewNumber(&app.cfg.Video.WindowWidth, config.WindowWidthMin, config.WindowWidthMax, "Window width", ""),
		menu.NewNumber(&app.cfg.Video.WindowHeight, config.WindowHeightMin, config.WindowHeightMax, "Window height", ""),
		menu.NewEnum[float32](&app.cfg.Video.WindowOpacity,
			[]float32{1.0, 0.9, 0.8, 0.7, 0.6, 0.5, 0.4, 0.3, 0.2, 0.1},
			map[float32]string{
				1.0: "100% (Opaque)",
				0.9: "90%",
				0.8: "80%",
				0.7: "70%",
				0.6: "60%",
				0.5: "50%",
				0.4: "40%",
				0.3: "30%",
				0.2: "20%",
				0.1: "10%",
			},
			"Opacity", ""),
		app.menuItemEscape(),
		app.menuItemBack(),
	}...)
}
