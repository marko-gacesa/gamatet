// Copyright (c) 2024 by Marko Gaćeša

package menu

import (
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"strings"
	"unicode"
)

const cursor = "\x01"

type Text struct {
	editable
	value   *string
	editing bool
	editor  []rune
	cursor  int
	maxLen  int
}

func NewText(s *string, label, description string, maxLen int) *Text {
	if s == nil {
		panic("need non-nil pointer")
	}
	return &Text{
		editable: editable{
			base:  base{description: description},
			label: label,
		},
		value:  s,
		maxLen: maxLen,
	}
}

func (t *Text) Text() string {
	if t.current != "" {
		return t.current
	}

	if !t.editing {
		t.current = t.label + ": " + *t.value
	} else {
		t.current = t.label + ": " + string(t.editor[:t.cursor]) + cursor + string(t.editor[t.cursor:])
	}

	return t.current
}

func (t *Text) Increase() {
	if !t.editing {
		t.startEdit()
		return
	}

	c := t.cursor
	n := len(t.editor)
	if c < n {
		c++
	} else {
		c = n
	}

	if c != t.cursor {
		t.cursor = c
		t.dirty()
	}
}

func (t *Text) Decrease() {
	if !t.editing {
		t.startEdit()
		return
	}

	c := t.cursor
	if c > 0 {
		c--
	} else {
		c = 0
	}

	if c != t.cursor {
		t.cursor = c
		t.dirty()
	}
}

func (t *Text) Input(r rune) {
	switch {
	case r == '\n': // commit
		if t.editing {
			t.commit()
		} else {
			t.startEdit()
		}
		return
	case r == '\b': // backspace
		if t.cursor > 0 {
			t.editor = append(t.editor[:t.cursor-1], t.editor[t.cursor:]...)
			t.cursor--
			t.dirty()
		}
		return
	case r == '\xFF': // delete
		if t.cursor < len(t.editor) {
			t.editor = append(t.editor[:t.cursor], t.editor[t.cursor+1:]...)
			t.dirty()
		}
		return
	case unicode.IsPrint(r):
		if len(t.editor) >= t.maxLen {
			return
		}

		t.editor = append(t.editor[:t.cursor], append([]rune{r}, t.editor[t.cursor:]...)...)
		t.cursor++
		t.dirty()
	}
}

func (t *Text) FocusLost() {
	if !t.editing {
		return
	}

	t.editing = false
	t.dirty()
}

func (t *Text) startEdit() {
	s, _, _ := transform.String(norm.NFC, *t.value)
	t.editor = []rune(s)
	t.cursor = len(t.editor)
	t.editing = true
	t.dirty()
}

func (t *Text) commit() {
	*t.value = strings.TrimSpace(string(t.editor))
	t.editing = false
	t.dirty()
}
