// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package config

import (
	"github.com/marko-gacesa/gamatet/internal/config/key"
	"github.com/marko-gacesa/gamatet/logic/gamepad"
)

type Input struct {
	Source  InputSource `json:"source"`
	Keys    key.Input   `json:"keys"`
	Gamepad int         `json:"gamepad"`
}

type InputSource string

const (
	InputSourceKeyboard InputSource = "keyboard"
	InputSourceGamepad  InputSource = "gamepad"
)

func DefaultInput(idx int) Input {
	return Input{
		Source:  InputSourceKeyboard,
		Keys:    key.DefaultInput[idx%len(key.DefaultInput)],
		Gamepad: idx,
	}
}

func (cfg *Input) Sanitize(idx int) {
	if cfg.Source != InputSourceKeyboard && cfg.Source != InputSourceGamepad {
		cfg.Source = InputSourceKeyboard
	}
	cfg.Keys.Sanitize(idx)
	if cfg.Gamepad < 0 || cfg.Gamepad >= gamepad.Count {
		cfg.Gamepad = idx
	}
}
