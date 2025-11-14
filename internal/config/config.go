// Copyright (c) 2024, 2025 by Marko Gaćeša

package config

type Config struct {
	Locale       Locale       `json:"locale"`
	LocalPlayers LocalPlayers `json:"local_players"`
	Network      Network      `json:"network"`
	Presets      Presets      `json:"presets"`
}

func (cfg *Config) Sanitize() {
	cfg.Locale.Sanitize()
	cfg.LocalPlayers.Sanitize()
	cfg.Presets.Sanitize()
	cfg.Network.Sanitize()
}

func sliceFixLen[T any](a []T, desiredLen int, genFn func(idx int) T) []T {
	if len(a) > desiredLen {
		return a[:desiredLen]
	}
	for i := len(a); i < desiredLen; i++ {
		a = append(a, genFn(i))
	}
	return a[:desiredLen:desiredLen]
}
