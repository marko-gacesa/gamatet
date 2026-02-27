// Copyright (c) 2025, 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	"fmt"

	"github.com/marko-gacesa/gamatet/graphics/monitor"
	"github.com/marko-gacesa/gamatet/internal/config"
	. "github.com/marko-gacesa/gamatet/internal/i18n"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) menuConfigVideo(ctx screen.Context) *menu.Menu {
	monitors := monitor.GetMonitors()
	names := make([]string, len(monitors)+1)
	names[0] = ""
	displayMap := make(map[string]string)
	displayMap[""] = T(KeyConfigVideoMonitorPrimary)
	for i, m := range monitors {
		names[i+1] = m.Name
		displayMap[m.Name] = m.String()
	}

	return menu.New(T(KeyConfigVideoTitle), app.configStopper(ctx), []menu.Item{
		menu.NewBool(&app.cfg.Video.Fullscreen,
			T(KeyConfigVideoFullscreen), T(KeyConfigVideoFullscreenDesc),
			withBoolStr()),
		menu.NewEnum[string](&app.cfg.Video.Monitor,
			names,
			func(s string) string {
				return displayMap[s]
			},
			"\t"+T(KeyConfigVideoMonitor), T(KeyConfigVideoMonitorDesc),
			menu.WithVisible(func() bool {
				return app.cfg.Video.Fullscreen
			})),
		menu.NewNumber(&app.cfg.Video.WindowWidth, config.WindowWidthMin, config.WindowWidthMax,
			"\t"+T(KeyConfigVideoWindowWidth), T(KeyConfigVideoWindowWidthDesc),
			menu.WithVisible(func() bool {
				return !app.cfg.Video.Fullscreen
			})),
		menu.NewNumber(&app.cfg.Video.WindowHeight, config.WindowHeightMin, config.WindowHeightMax,
			"\t"+T(KeyConfigVideoWindowHeight), T(KeyConfigVideoWindowHeightDesc),
			menu.WithVisible(func() bool {
				return !app.cfg.Video.Fullscreen
			})),
		menu.NewEnum[float32](&app.cfg.Video.WindowOpacity,
			[]float32{1.0, 0.9, 0.8, 0.7, 0.6, 0.5, 0.33},
			func(f float32) string {
				return fmt.Sprintf("%.0f%%", f*100)
			},
			"\t"+T(KeyConfigVideoOpacity), T(KeyConfigVideoOpacityDesc),
			menu.WithVisible(func() bool {
				return !app.cfg.Video.Fullscreen
			})),
		app.menuItemEscape(),
		app.menuItemBack(),
	}...)
}
