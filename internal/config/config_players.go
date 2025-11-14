// Copyright (c) 2025 by Marko Gaćeša

package config

import (
	"github.com/marko-gacesa/gamatet/game/piece"
	"github.com/marko-gacesa/gamatet/game/setup"
	"github.com/marko-gacesa/gamatet/internal/config/key"
)

type LocalPlayers struct {
	Infos []PlayerInfo `json:"infos"`
}

func (cfg *LocalPlayers) Sanitize() {
	cfg.Infos = sliceFixLen(cfg.Infos, setup.MaxLocalPlayers, func(idx int) PlayerInfo {
		return PlayerInfo{
			Name:  "",
			Input: key.DefaultInput[idx%len(key.DefaultInput)],
			GameConfig: Player{
				RotationDirectionCW: false,
				SlideDisabled:       false,
				WallKick:            piece.WallKickDefault,
			},
		}
	})

	for i := range cfg.Infos {
		cfg.Infos[i].Sanitize(i)
	}

}

type PlayerInfo struct {
	Name       string    `json:"name"`
	Input      key.Input `json:"input"`
	GameConfig Player    `json:"game_config"`
}

func (cfg *PlayerInfo) Sanitize(idx int) {
	if cfg.Name == "" || len(cfg.Name) > setup.MaxLenName {
		switch idx % 4 {
		case 0:
			cfg.Name = string('\u0394') // delta
		case 1:
			cfg.Name = string('\u0398') // theta
		case 2:
			cfg.Name = string('\u03A8') // psi
		case 3:
			cfg.Name = string('\u03A9') // omega
		}
	}
	cfg.Input.Sanitize(idx)
	cfg.GameConfig.Sanitize()
}

type Player struct {
	RotationDirectionCW bool `json:"rotation_direction_cw"`
	SlideDisabled       bool `json:"slide_disabled"`
	WallKick            byte `json:"wall_kick"`
}

func (cfg Player) Serialize() []byte {
	return setup.Pack((*setup.PlayerConfig)(&cfg))
}

func (cfg *Player) Sanitize() {
	if cfg.WallKick > piece.WallKickMax {
		cfg.WallKick = piece.WallKickDefault
	}
}
