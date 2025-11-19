// Copyright (c) 2025 by Marko Gaćeša

package config

import (
	"os"
	"strings"

	"github.com/marko-gacesa/gamatet/logic/lang"
)

const (
	defaultLanguage = "en"
)

type Locale struct {
	Language string
}

func (cfg *Locale) Sanitize() {
	if cfg.Language == "" {
		const envLang = "LANG"
		lang := os.Getenv(envLang)
		idx := strings.IndexAny(lang, "._-")
		if idx < 0 {
			cfg.Language = defaultLanguage
			return
		}

		cfg.Language = lang[:idx]
	}

	cfg.Language = strings.ToLower(cfg.Language)

	var found bool
	for _, l := range lang.Supported() {
		if cfg.Language == string(l) {
			found = true
			break
		}
	}
	if !found {
		cfg.Language = defaultLanguage
	}
}
