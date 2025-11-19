// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"fmt"

	"github.com/marko-gacesa/gamatet/internal/config"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) menuConfigVideo(ctx screen.Context) *menu.Menu {
	return menu.New("Video Options", app.configStopper(ctx), []menu.Item{
		menu.NewBool(&app.cfg.Video.Fullscreen, "Fullscreen", ""),
		menu.NewNumber(&app.cfg.Video.WindowWidth, config.WindowWidthMin, config.WindowWidthMax, "Window width", ""),
		menu.NewNumber(&app.cfg.Video.WindowHeight, config.WindowHeightMin, config.WindowHeightMax, "Window height", ""),
		menu.NewEnum[float32](&app.cfg.Video.WindowOpacity,
			[]float32{1.0, 0.9, 0.8, 0.7, 0.6, 0.5, 0.4, 0.3, 0.2, 0.1},
			func(f float32) string {
				return fmt.Sprintf("%.0f%%", f*100)
			},
			"Opacity", ""),
		app.menuItemEscape(),
		app.menuItemBack(),
	}...)
}
