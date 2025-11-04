// Copyright (c) 2024, 2025 by Marko Gaćeša

package config

import (
	"gamatet/game/piece"
	"gamatet/game/setup"
	"os"
	"strings"
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
			Name: "",
			PlayerConfig: PlayerConfig{
				RotationDirectionCW: false,
				SlideDisabled:       false,
				WallKick:            piece.WallKickDefault,
			},
		}
	})

	for _, info := range cfg.Infos {
		info.Sanitize()
	}
}

type PlayerInfo struct {
	Name string `json:"name"`
	PlayerConfig
}

func (cfg *PlayerInfo) Sanitize() {
	if cfg.Name == "" || len(cfg.Name) > setup.MaxLenName {
		cfg.Name = "Player"
	}

	cfg.PlayerConfig.Sanitize()
}

type Presets struct {
	Single []setup.Setup `json:"single"`
	Multi  []setup.Setup `json:"multi"`
}

func (cfg *Presets) Sanitize() {
	cfg.Single = SliceFixLen(cfg.Single, setup.SinglePlayerPresetCount, setup.SinglePlayerPreset)
	for i := range cfg.Single {
		cfg.Single[i].SanitizeSingle()
	}

	cfg.Multi = SliceFixLen(cfg.Multi, setup.MultiPlayerPresetCount, setup.MultiPlayerPreset)
	for i := range cfg.Multi {
		cfg.Multi[i].SanitizeMulti()
	}
}
