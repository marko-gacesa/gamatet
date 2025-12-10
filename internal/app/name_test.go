// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	"testing"

	"github.com/marko-gacesa/gamatet/internal/i18n"
	"github.com/marko-gacesa/gamatet/logic/lang"
)

func TestName(t *testing.T) {
	tests := map[string]string{
		"":           "P 1",      // empty
		" ":          "P 1",      // whitespace only
		"A\t\n B":    "A B",      // eol, tab
		"abc":        "abc",      // latin
		"123":        "123",      // numbers
		"Δ":          "Δ",        // greek
		"Ћ":          "Ћ",        // cyrillic
		".:[]{}$@":   ".:[]{}$@", // ascii punctuation
		"世界123":    "123",      // remove japanese/chinese
		"عالم123":    "123",      // remove arabic
		"\xFBAB\xF0": "AB",       // remove invalid

		"0123456789-too-long-string-превише": "0123456789-too-long-string-преви",
	}

	lang.DefineFallback(map[string]string{
		i18n.KeyFieldPlayerNameDefault: "P %s",
	})

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			got := playerName(input, 0, 0, 0)
			if got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		})
	}
}
