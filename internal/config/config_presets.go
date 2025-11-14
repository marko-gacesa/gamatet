// Copyright (c) 2025 by Marko Gaćeša

package config

import "github.com/marko-gacesa/gamatet/game/setup"

type Presets struct {
	Single       []setup.Setup `json:"single"`
	SingleCustom setup.Setup   `json:"single_custom"`
	Multi        []setup.Setup `json:"multi"`
	MultiCustom  setup.Setup   `json:"multi_custom"`
}

func (cfg *Presets) Sanitize() {
	cfg.Single = sliceFixLen(cfg.Single, setup.SinglePlayerPresetCount, setup.SinglePlayerPreset)
	for i := range cfg.Single {
		cfg.Single[i].SanitizeSingle()
		cfg.Single[i].SanitizeName()
	}

	if cfg.SingleCustom.Empty() {
		cfg.SingleCustom = setup.SinglePlayerPreset(0)
	}
	cfg.SingleCustom.SanitizeSingle()
	cfg.SingleCustom.Name = ""

	cfg.Multi = sliceFixLen(cfg.Multi, setup.MultiPlayerPresetCount, setup.MultiPlayerPreset)
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
