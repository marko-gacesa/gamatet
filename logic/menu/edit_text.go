// Copyright (c) 2024, 2025 by Marko Gaćeša

package menu

import (
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var _ Item = (*Text)(nil)

// Text is menu item that assigns a value to a string variable.
type Text struct {
	textBase
	ptr *string
}

// NewText creates new Text menu item.
func NewText(ptr *string, maxLen, maxSize int, label, description string, options ...func(Item)) *Text {
	if ptr == nil {
		panic("need non-nil pointer")
	}
	t := &Text{
		textBase: makeTextBase(maxLen, maxSize, label, description),
		ptr:      ptr,
	}
	t.textBase.converter = t
	t.fix()
	applyOptions(t, options...)
	return t
}

func (t *Text) fix() {
	runeLen := 0
	for i, r := range *t.ptr {
		runeLen++
		byteLen := i + utf8.RuneLen(r)
		if runeLen > t.maxLen || byteLen > t.maxSize || !t.allowed(r) {
			*t.ptr = (*t.ptr)[:i]
			t.markDirty()
			return
		}
	}
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
