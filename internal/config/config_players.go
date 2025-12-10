// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

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

func (cfg *LocalPlayers) Inputs() [setup.MaxLocalPlayers]key.Input {
	var inputs [setup.MaxLocalPlayers]key.Input
	for i := range min(len(cfg.Infos), setup.MaxLocalPlayers) {
		inputs[i] = cfg.Infos[i].Input
	}
	return inputs
}

type PlayerInfo struct {
	Name       string    `json:"name"`
	Input      key.Input `json:"input"`
	GameConfig Player    `json:"game_config"`
}

func (cfg *PlayerInfo) Sanitize(idx int) {
	if len(cfg.Name) > setup.MaxLenName {
		cfg.Name = string([]rune(cfg.Name)[:setup.MaxLenName])
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
