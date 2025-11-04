// Copyright (c) 2025 by Marko Gaćeša

package config

import (
	"gamatet/game/piece"
	"gamatet/game/setup"
)

type PlayerConfig struct {
	RotationDirectionCW bool `json:"rotation_direction_cw"`
	SlideDisabled       bool `json:"slide_disabled"`
	WallKick            byte `json:"wall_kick"`
}

func (cfg PlayerConfig) Serialize() []byte {
	return setup.Pack((*setup.PlayerConfig)(&cfg))
}

func (cfg *PlayerConfig) Sanitize() {
	if cfg.WallKick > piece.WallKickMax {
		cfg.WallKick = piece.WallKickDefault
	}
}
