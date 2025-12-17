// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

const (
	LevelNone  = "none"
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

func Logger(logLevel string) *slog.Logger {
	var level slog.Level

	switch strings.ToLower(logLevel) {
	case LevelNone:
		return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
			AddSource:   false,
			Level:       slog.LevelError,
			ReplaceAttr: nil,
		}))
	case LevelDebug:
		level = slog.LevelDebug
	case LevelInfo, "":
		level = slog.LevelInfo
	case LevelWarn:
		level = slog.LevelWarn
	case LevelError:
		level = slog.LevelError
	default:
		fmt.Printf("Invalid log level %q. Valid values are %q, %q, %q, %q, or %q.\n",
			logLevel, LevelNone, LevelDebug, LevelInfo, LevelWarn, LevelError)
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       level,
		ReplaceAttr: nil,
	}))

	slog.SetDefault(logger)

	return logger
}
