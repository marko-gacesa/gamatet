// Copyright (c) 2020-2025 by Marko Gaćeša

package main

import (
	"context"
	"gamatet/graphics/loop"
	"gamatet/internal/app"
	"gamatet/internal/config"
	"gamatet/logic/appctx"
	"os"
)

func main() {
	globalCtx := appctx.Context

	logger := app.Logger()

	cfg, cfgPath := config.Load(logger)

	pid := os.Getpid()

	appCtx, appCtxStop := context.WithCancel(globalCtx)

	app := app.NewApp(appCtx, logger, cfg, cfgPath)
	app.Log().Info("Starting", "pid", pid, "cfgPath", cfgPath)

	err := loop.Loop(appCtx, app)
	if err != nil {
		app.Log().Error("Stopped", "error", err)
	}

	appCtxStop()

	app.WaitDone()
}
