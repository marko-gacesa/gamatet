// Copyright (c) 2025 by Marko Gaćeša

package config

const (
	WindowWidthMin  = 768
	WindowWidthMax  = 7680
	WindowHeightMin = 432
	WindowHeightMax = 4320
)

type Video struct {
	Fullscreen    bool    `json:"fullscreen"`
	WindowWidth   int     `json:"window_width"`
	WindowHeight  int     `json:"window_height"`
	WindowOpacity float32 `json:"window_opacity"`
}

func (cfg *Video) Sanitize() {
	if cfg.WindowWidth == 0 {
		cfg.WindowWidth = 1400
	} else if cfg.WindowWidth < WindowWidthMin {
		cfg.WindowWidth = WindowWidthMin
	} else if cfg.WindowWidth > WindowWidthMax {
		cfg.WindowWidth = WindowWidthMax
	}

	if cfg.WindowHeight == 0 {
		cfg.WindowHeight = 800
	} else if cfg.WindowHeight < WindowHeightMin {
		cfg.WindowHeight = WindowHeightMin
	} else if cfg.WindowHeight > WindowHeightMax {
		cfg.WindowHeight = WindowHeightMax
	}

	if cfg.WindowOpacity == 0 {
		cfg.WindowOpacity = 1
	} else if cfg.WindowOpacity <= 0.1 {
		cfg.WindowOpacity = 0.1
	} else if cfg.WindowOpacity > 1.0 {
		cfg.WindowOpacity = 1.0
	}
}
