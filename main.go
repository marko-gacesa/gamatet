// Copyright (c) 2020-2025 by Marko Gaćeša

package main

import (
	"context"
	"github.com/marko-gacesa/gamatet/graphics/loop"
	"github.com/marko-gacesa/gamatet/internal/app"
	"github.com/marko-gacesa/gamatet/internal/config"
	"github.com/marko-gacesa/gamatet/internal/values"
	"github.com/marko-gacesa/gamatet/logic/appctx"
	"os"
)

func main() {
	pid := os.Getpid()

	logger := app.Logger()
	logger.Info(values.ProgramName,
		"version", values.VersionTag,
		"commit_sha", values.GitSHA,
		"build_time", values.BuildTime,
		"pid", pid)

	cfg, cfgPath := config.Load(logger)

	globalCtx := appctx.Context
	appCtx, appCtxStop := context.WithCancel(globalCtx)

	app := app.NewApp(appCtx, logger, cfg, cfgPath)
	app.Log().Info("Config", "cfg_path", cfgPath)

	err := loop.Loop(appCtx, app)
	if err != nil {
		app.Log().Error("Stopped", "error", err)
	}

	appCtxStop()

	app.WaitDone()
}
