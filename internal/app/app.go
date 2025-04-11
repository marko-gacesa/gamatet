// Copyright (c) 2024 by Marko Gaćeša

package app

import (
	"context"
	"gamatet/graphics/scene"
	"gamatet/internal/config"
	"gamatet/logic/screen"
	"log/slog"
	"os"
)

type App struct {
	cfg     config.Config
	cfgPath string

	screenIDHistory *routes // screen history, the last entry is the id currently active screen
	screenIDNext    route

	screener screen.Screener

	logger *slog.Logger
}

func NewApp(cfg config.Config, cfgPath string) *App {
	return &App{
		cfg:             cfg,
		cfgPath:         cfgPath,
		screenIDHistory: (&routes{}).push(routeMain),
		logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   false,
			Level:       slog.LevelDebug,
			ReplaceAttr: nil,
		})),
	}
}

func (app *App) SetScreener(screener screen.Screener) {
	app.screener = screener
}

func (app *App) Log() *slog.Logger {
	return app.logger
}

func (app *App) MakeScreen(parentCtx context.Context) (screen.Screen, <-chan struct{}) {
	id := app.screenIDHistory.curr()
	var data any

	ctx := screen.NewContext(parentCtx)

	switch id {
	case "", routeQuit:
		data = nil
	case routeMain:
		data = app.menuMain(ctx)
	case routeMenuSinglePlayer:
		data = app.menuSinglePlayer(ctx)
	case routeTestBlocks:
		data = scene.Demo(scene.DemoBlocks)
	case routeTestField:
		data = scene.Demo(scene.DemoFields)
	case routeGameSinglePlayNow:
		data = app.gameOne(ctx)
	case routeGameDoublePlayNow:
		data = app.gameDouble(ctx)
	}

	return app.screener.Screen(ctx, data), ctx.Done()
}

func (app *App) ScreenFinish() {
	if app.screenIDNext == routeBack || app.screenIDNext == "" {
		app.screenIDHistory.pop()
	} else if app.screenIDNext != "" {
		app.screenIDHistory.push(app.screenIDNext)
		app.screenIDNext = ""
	}
}
