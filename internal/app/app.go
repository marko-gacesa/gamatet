// Copyright (c) 2024 by Marko Gaćeša

package app

import (
	"context"
	"gamatet/graphics/scene"
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

func (app *App) MakeScreen(ctx context.Context) (screen.Screen, context.Context) {
	id := app.screenIDHistory.curr()
	var data any

	switch id {
	case "", routeQuit:
		data = nil
	case routeMain:
		var cancelCtx context.CancelFunc
		ctx, cancelCtx = context.WithCancel(ctx)
		data = app.menuMain(cancelCtx)
	case routeMenuSinglePlayer:
		var cancelCtx context.CancelFunc
		ctx, cancelCtx = context.WithCancel(ctx)
		data = app.menuSinglePlayer(cancelCtx)
	case routeTestBlocks:
		data, ctx = scene.Demo(ctx, scene.DemoBlocks)
	case routeTestField:
		data, ctx = scene.Demo(ctx, scene.DemoFields)
	case routeGameSinglePlayNow:
		data, ctx = app.gameOne(ctx)
	case routeGameDoublePlayNow:
		data, ctx = app.gameDouble(ctx)
	}

	return app.screener.Screen(ctx, data), ctx
}

func (app *App) ScreenFinish() {
	if app.screenIDNext == routeBack || app.screenIDNext == "" {
		app.screenIDHistory.pop()
	} else if app.screenIDNext != "" {
		app.screenIDHistory.push(app.screenIDNext)
		app.screenIDNext = ""
	}
}
