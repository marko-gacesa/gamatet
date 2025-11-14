// Copyright (c) 2024, 2025 by Marko Gaćeša

package config

import (
	"os"
	"strings"

	"github.com/marko-gacesa/gamatet/game/piece"
	"github.com/marko-gacesa/gamatet/game/setup"
	"github.com/marko-gacesa/gamatet/internal/config/key"
)

const (
	defaultLanguage = "en"
)

type Config struct {
	Language     string       `json:"language"`
	LocalPlayers LocalPlayers `json:"local_players"`
	Network      Network      `json:"network"`
	Presets      Presets      `json:"presets"`
}

func (cfg *Config) Sanitize() {
	cfg.SanitizeLanguage()
	cfg.LocalPlayers.Sanitize()
	cfg.Presets.Sanitize()
	cfg.Network.Sanitize()
}

func (cfg *Config) SanitizeLanguage() {
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

	// TODO: Check if lang is supported
}

type LocalPlayers struct {
	Infos []PlayerInfo `json:"infos"`
}

func (cfg *LocalPlayers) Sanitize() {
	cfg.Infos = SliceFixLen(cfg.Infos, setup.MaxLocalPlayers, func(idx int) PlayerInfo {
		return PlayerInfo{
			Name:  "",
			Input: key.DefaultInput[idx%len(key.DefaultInput)],
			GameConfig: PlayerConfig{
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
	Name       string       `json:"name"`
	Input      key.Input    `json:"input"`
	GameConfig PlayerConfig `json:"game_config"`
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

type Presets struct {
	Single       []setup.Setup `json:"single"`
	SingleCustom setup.Setup   `json:"single_custom"`
	Multi        []setup.Setup `json:"multi"`
	MultiCustom  setup.Setup   `json:"multi_custom"`
}

func (cfg *Presets) Sanitize() {
	cfg.Single = SliceFixLen(cfg.Single, setup.SinglePlayerPresetCount, setup.SinglePlayerPreset)
	for i := range cfg.Single {
		cfg.Single[i].SanitizeSingle()
		cfg.Single[i].SanitizeName()
	}

	if cfg.SingleCustom.Empty() {
		cfg.SingleCustom = setup.SinglePlayerPreset(0)
	}
	cfg.SingleCustom.SanitizeSingle()
	cfg.SingleCustom.Name = ""

	cfg.Multi = SliceFixLen(cfg.Multi, setup.MultiPlayerPresetCount, setup.MultiPlayerPreset)
	for i := range cfg.Multi {
		cfg.Multi[i].SanitizeMulti()
		cfg.Multi[i].SanitizeName()
	}

	if cfg.MultiCustom.Empty() {
		cfg.MultiCustom = setup.MultiPlayerPreset(0)
	}
	cfg.MultiCustom.SanitizeMulti()
	cfg.MultiCustom.Name = ""
}
