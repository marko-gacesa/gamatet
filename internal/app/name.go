// Copyright (c) 2025, 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"

	"github.com/marko-gacesa/gamatet/game/setup"
	"github.com/marko-gacesa/gamatet/internal/config"
	"github.com/marko-gacesa/gamatet/internal/config/key"
	. "github.com/marko-gacesa/gamatet/internal/i18n"
	"github.com/marko-gacesa/gamatet/logic/gamepad"
)

func playerName(name string, fieldIdx, fieldPlayerIdx, playerIdx int) string {
	name, _, err := transform.String(
		transform.Chain(
			runes.ReplaceIllFormed(),
			runes.Remove(runes.Predicate(func(r rune) bool {
				return r == utf8.RuneError
			})),
			runes.Map(func(r rune) rune {
				if unicode.IsSpace(r) {
					return ' ' // convert all whitespace to ordinary space
				}
				return r
			}),
			runes.Remove(runes.Predicate(func(r rune) bool {
				return (r < ' ' || r > 127) &&
					!unicode.Is(unicode.Latin, r) &&
					!unicode.Is(unicode.Greek, r) &&
					!unicode.Is(unicode.Cyrillic, r)
			})),
		),
		name,
	)
	if err != nil {
		name = ""
	}

	name = strings.TrimSpace(name)

	for strings.Contains(name, "  ") {
		name = strings.Replace(name, "  ", " ", -1)
	}

	if name == "" {
		return Tf(KeyFieldPlayerNameDefault, strconv.Itoa(playerIdx+1))
	}

	if utf8.RuneCountInString(name) > setup.MaxLenName {
		name = string([]rune(name)[:setup.MaxLenName])
	}

	return name
}

func gameInput(input config.Input) []string {
	switch input.Source {
	case config.InputSourceKeyboard:
		return []string{
			T(KeyConfigPlayerKeyLeft) + ": " + key.Map[input.Keys.Left],
			T(KeyConfigPlayerKeyRight) + ": " + key.Map[input.Keys.Right],
			T(KeyConfigPlayerKeyActivate) + ": " + key.Map[input.Keys.Activate],
			T(KeyConfigPlayerKeyBoost) + ": " + key.Map[input.Keys.Boost],
			T(KeyConfigPlayerKeyDrop) + ": " + key.Map[input.Keys.Drop],
		}
	case config.InputSourceGamepad:
		idx := input.Gamepad
		if idx < 0 || idx > gamepad.Count {
			return nil
		}

		if !gamepad.Gamepads[idx].Connected {
			return []string{T(KeyInputSourceGamepad), T(KeyDeviceNotConnected)}
		}

		return splitTextToLines(gamepad.Gamepads[idx].Name, 20)
	}

	return nil
}

func splitTextToLines(s string, l int) (out []string) {
	words := strings.Split(s, " ")

	for i := 0; i < len(words); {
		words[i] = strings.TrimSpace(words[i])
		if words[i] == "" {
			words = append(words[:i], words[i+1:]...)
		} else {
			i++
		}
	}

	if len(words) == 0 {
		return nil
	}

	for i := 1; i < len(words); {
		if len(words[i-1])+len(words[i]) < l {
			words[i-1] = words[i-1] + " " + words[i]
			words = append(words[:i], words[i+1:]...)
		} else {
			i++
		}
	}

	return words
}
