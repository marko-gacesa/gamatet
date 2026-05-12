// Copyright (c) 2025, 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	"slices"
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
		"世界123":      "123",      // remove japanese/chinese
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

func TestSplitTextToLines(t *testing.T) {
	tests := []struct {
		text string
		l    int
		exp  []string
	}{
		{
			text: "",
			l:    10,
			exp:  nil,
		},
		{
			text: "one_word",
			l:    10,
			exp:  []string{"one_word"},
		},
		{
			text: "one_word_small_l",
			l:    5,
			exp:  []string{"one_word_small_l"},
		},
		{
			text: "two words",
			l:    20,
			exp:  []string{"two words"},
		},
		{
			text: "  with    spaces ",
			l:    20,
			exp:  []string{"with spaces"},
		},
		{
			text: "several long words",
			l:    15,
			exp:  []string{"several long", "words"},
		},
		{
			text: "" +
				"Laborum rerum sed iure aspernatur sed qui voluptas. Error ut magni ipsum itaque unde veritatis. " +
				"Molestiae similique repellendus nostrum in repellat enim sequi rerum. Aut autem totam ut.",
			l: 20,
			exp: []string{
				"Laborum rerum sed", "iure aspernatur sed", "qui voluptas. Error", "ut magni ipsum", "itaque unde",
				"veritatis. Molestiae", "similique", "repellendus nostrum", "in repellat enim", "sequi rerum. Aut",
				"autem totam ut.",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.text, func(t *testing.T) {
			if want, got := test.exp, splitTextToLines(test.text, test.l); !slices.Equal(want, got) {
				t.Errorf("want=%v got=%v", want, got)
			}
		})
	}
}
