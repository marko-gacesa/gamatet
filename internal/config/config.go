// Copyright (c) 2024 by Marko Gaćeša

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

const (
	filename = ".gamatet.config.json"

	defaultLanguage = "en"
)

type Config struct {
	Language string `json:"language"`
}

func Default() Config {
	return Config{
		Language: defaultLanguage,
	}
}

func Load() (Config, string) {
	var cfg Config

	dirs := getDirList()
	for _, dir := range dirs {
		fn := path.Join(dir, filename)
		f, err := os.Open(fn)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			fmt.Printf("failed to open config file %s: %s\n", fn, err.Error())
			continue
		}

		if err := func() error {
			defer f.Close()
			return json.NewDecoder(f).Decode(&cfg)
		}(); err != nil {
			fmt.Printf("failed to load config from %s: %s\n", fn, err.Error())
		}

		return cfg, dir
	}

	if len(dirs) > 0 {
		return Default(), path.Join(dirs[0], filename)
	}

	return Default(), filename
}

func getDirList() []string {
	var dirs []string

	dir, _ := os.UserHomeDir()
	if dir != "" {
		dirs = append(dirs, dir)
	}

	if len(os.Args) > 0 {
		dir = path.Dir(os.Args[0])
		if dir != "" {
			dirs = append(dirs, dir)
		}
	}

	return dirs
}
