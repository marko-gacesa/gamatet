// Copyright (c) 2020-2025 by Marko Gaćeša

package main

import (
	"context"
	"os"

	"github.com/marko-gacesa/gamatet/internal/i18n"

	"github.com/marko-gacesa/gamatet/graphics/loop"
	"github.com/marko-gacesa/gamatet/internal/app"
	"github.com/marko-gacesa/gamatet/internal/config"
	"github.com/marko-gacesa/gamatet/internal/values"
	"github.com/marko-gacesa/gamatet/logic/appctx"
)

func main() {
	pid := os.Getpid()

	logger := app.Logger()
	logger.Info(values.ProgramName,
		"version", values.VersionTag,
		"commit_sha", values.GitSHA,
		"build_time", values.BuildTime,
		"pid", pid)

	i18n.ParseEmbeddedLanguages(logger)

	cfg, cfgPath := config.Load(logger)
	cfg.Sanitize()

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
