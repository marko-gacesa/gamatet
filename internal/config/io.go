// Copyright (c) 2024, 2025 by Marko Gaćeša

package config

import (
	"encoding/json"
	"log/slog"
	"os"
	"path"
)

const filename = ".gamatet.config.json"

func Load(log *slog.Logger) (Config, string) {
	var cfg Config

	dirs := getDirList()
	for _, dir := range dirs {
		filePath := path.Join(dir, filename)
		f, err := os.Open(filePath)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			log.Error("failed to open config file", "error", err, "path", filePath)
			continue
		}

		if err := func() error {
			defer f.Close()
			return json.NewDecoder(f).Decode(&cfg)
		}(); err != nil {
			log.Error("failed to load config", "error", err, "path", filePath)
		}

		return cfg, dir
	}

	cfg.Sanitize()

	if len(dirs) > 0 {
		return cfg, path.Join(dirs[0], filename)
	}

	return cfg, filename
}

func getDirList() []string {
	var dirs []string

	// home dir
	dir, _ := os.UserHomeDir()
	if dir != "" {
		dirs = append(dirs, dir)
	}

	// exec's dir
	if len(os.Args) > 0 {
		dir = path.Dir(os.Args[0])
		if dir != "" {
			dirs = append(dirs, dir)
		}
	}

	return dirs
}
