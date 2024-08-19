// Copyright (c) 2024 by Marko Gaćeša

package app

import (
	"context"
	"gamatet/internal/config"
	"gamatet/logic/screen"
)

type App struct {
	cfg     config.Config
	cfgPath string

	screenIDHistory *routes // screen history, the last entry is the id currently active screen
	screenIDNext    route

	screener screen.Screener
}

func NewApp(cfg config.Config, cfgPath string) *App {
	return &App{
		cfg:             cfg,
		cfgPath:         cfgPath,
		screenIDHistory: (&routes{}).push(routeMain),
	}
}

func (app *App) SetScreener(screener screen.Screener) {
	app.screener = screener
}

func (app *App) MakeScreen(ctx context.Context) screen.Screen {
	id := app.screenIDHistory.curr()
	var data any

	switch id {
	case routeMain:
		data = app.menuMain()
	case "", routeQuit:
		data = nil
	case routeTestBlocks:
		data = "test-blocks"
	case routeTestField:
		data = "test-fields"
	}

	return app.screener.Screen(ctx, data)
}

func (app *App) ScreenFinish() {
	if app.screenIDNext != "" {
		app.screenIDHistory.push(app.screenIDNext)
		app.screenIDNext = ""
	} else {
		app.screenIDHistory.pop()
	}
}
