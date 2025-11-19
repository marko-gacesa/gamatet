// Copyright (c) 2025 by Marko Gaćeša

package lang

import (
	"maps"
	"slices"
	"testing"
)

func TestDefine(t *testing.T) {
	const (
		key1 = "key1"
		key2 = "key2"
	)

	Define("en", map[string]string{
		key1: "value1",
		key2: "value2",
	})
	Define("sr", map[string]string{
		key1: "vrednost1",
	})

	Set("en")

	if want, got := "value1", Str(key1); want != got {
		t.Errorf("%s: want %q, got %q", key1, want, got)
	}

	Set("sr")

	if want, got := key2, Str(key2); want != got {
		t.Errorf("%s: want %q, got %q", key2, want, got)
	}

	DefineFallbackFromExisting("en")

	if want, got := "value2", Str(key2); want != got {
		t.Errorf("%s, after fallback: want %q, got %q", key2, want, got)
	}

	if want, got := []Lang{"en", "sr"}, Supported(); !slices.Equal(want, got) {
		t.Errorf("supported: want %q, got %q", want, got)
	}

	if want, got := map[Lang]string{"en": "value1", "sr": "vrednost1"}, StrInAll(key1); !maps.Equal(want, got) {
		t.Errorf("str-in-all: want %q, got %q", want, got)
	}
}
