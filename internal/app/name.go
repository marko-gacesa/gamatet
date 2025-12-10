// Copyright (c) 2025 by Marko Gaćeša
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
	. "github.com/marko-gacesa/gamatet/internal/i18n"
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
