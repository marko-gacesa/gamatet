// Copyright (c) 2024, 2025 by Marko Gaćeša

package config

import (
	"os"
	"strings"

	"github.com/marko-gacesa/gamatet/game/piece"
	"github.com/marko-gacesa/gamatet/game/setup"
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

	for i := range cfg.Infos {
		if cfg.Infos[i].Name == "" || len(cfg.Infos[i].Name) > setup.MaxLenName {
			switch i % 4 {
			case 0:
				cfg.Infos[i].Name = string('\u0394') // delta
			case 1:
				cfg.Infos[i].Name = string('\u0398') // theta
			case 2:
				cfg.Infos[i].Name = string('\u03A8') // psi
			case 3:
				cfg.Infos[i].Name = string('\u03A9') // omega
			}
		}

		cfg.Infos[i].Sanitize()
	}
}

type PlayerInfo struct {
	Name string `json:"name"`
	PlayerConfig
}

func (cfg *PlayerInfo) Sanitize() {
	cfg.PlayerConfig.Sanitize()
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
