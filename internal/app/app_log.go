// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"log/slog"
	"os"
)

func Logger() *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelInfo,
		ReplaceAttr: nil,
	}))
	return logger
}
