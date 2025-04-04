// Copyright (c) 2024,2025 by Marko Gaćeša

package menu

import (
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"strings"
	"unicode"
	"unicode/utf8"
)

var _ Item = (*Text)(nil)

// Text is menu item that assigns a value to a string variable.
type Text struct {
	textBase
	ptr *string
}

// NewText creates new Text menu item.
func NewText(ptr *string, maxLen int, label, description string, options ...func(Item)) *Text {
	if ptr == nil {
		panic("need non-nil pointer")
	}
	t := &Text{
		textBase: makeTextBase(maxLen, label, description),
		ptr:      ptr,
	}
	t.textBase.converter = t
	t.fix()
	applyOptions(t, options...)
	return t
}

func (t *Text) fix() {
	valid := true
	for i, r := range *t.ptr {
		if i >= t.maxLen || !t.allowed(r) {
			valid = false
			break
		}
	}
	if valid {
		return
	}

	var sb strings.Builder
	var l int
	for _, r := range *t.ptr {
		if t.allowed(r) {
			sb.WriteRune(r)
			l++
		}
		if l == t.maxLen {
			break
		}
	}

	*t.ptr = sb.String()
	t.markDirty()
}

func (t *Text) getValueAsStr() string {
	s, _, _ := transform.String(norm.NFC, *t.ptr)
	return s
}

func (t *Text) setValueFromStr(s string) {
	*t.ptr = s
}

func (*Text) allowed(r rune) bool {
	return utf8.ValidRune(r) && unicode.IsPrint(r)
}

func (*Text) allowedInsert(rune, []rune, int) bool {
	return true
}
