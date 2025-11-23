// Copyright (c) 2024, 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package menu

import (
	"strings"
)

const (
	cursor      = '\x01'
	markerStart = '\x02'
	markerEnd   = '\x03'
)

type textBase struct {
	base
	editing   bool
	marked    bool
	editor    []rune
	cursor    int
	maxLen    int
	maxSize   int
	converter valueConverter
}

type valueConverter interface {
	getValueAsStr() string
	setValueFromStr(string)
	allowed(rune) bool
	allowedInsert(rune, []rune, int) bool
}

func makeTextBase(maxLen, maxSize int, label, description string) textBase {
	return textBase{
		base:    makeBase(label, description),
		maxLen:  maxLen,
		maxSize: maxSize,
	}
}

func (t *textBase) Text() string {
	if t.current != "" {
		return t.current
	}

	sb := strings.Builder{}
	sb.WriteString(t.getLabel())
	sb.WriteString(": ")
	if !t.editing {
		sb.WriteString(t.converter.getValueAsStr())
	} else if t.marked {
		if s := string(t.editor[:t.cursor]); len(s) > 0 {
			sb.WriteByte(markerStart)
			sb.WriteString(s)
			sb.WriteByte(markerEnd)
		}
		sb.WriteByte(cursor)
		if s := string(t.editor[t.cursor:]); len(s) > 0 {
			sb.WriteByte(markerStart)
			sb.WriteString(s)
			sb.WriteByte(markerEnd)
		}
	} else {
		sb.WriteString(string(t.editor[:t.cursor]))
		sb.WriteByte(cursor)
		sb.WriteString(string(t.editor[t.cursor:]))
	}
	t.current = sb.String()

	return t.current
}

func (t *textBase) increase() {
	if !t.editing {
		t.startEdit()
		return
	}

	if t.marked {
		t.marked = false
		t.markDirty()
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
		t.markDirty()
	}
}

func (t *textBase) decrease() {
	if !t.editing {
		t.startEdit()
		return
	}

	if t.marked {
		t.marked = false
		t.markDirty()
	}

	c := t.cursor
	if c > 0 {
		c--
	} else {
		c = 0
	}

	if c != t.cursor {
		t.cursor = c
		t.markDirty()
	}
}

func (t *textBase) input(r rune) bool {
	switch {
	case r == InputEnter:
		if t.editing {
			t.converter.setValueFromStr(string(t.editor))
			t.stopEdit()
		} else {
			t.startEdit()
		}
		return true
	case r == InputBackspace:
		if !t.editing {
			t.editing = true
			t.marked = false
			t.editor = nil
			t.cursor = 0
			t.markDirty()
		} else if t.marked {
			t.marked = false
			t.editor = t.editor[:0]
			t.cursor = 0
			t.markDirty()
		} else if t.cursor > 0 {
			t.editor = append(t.editor[:t.cursor-1], t.editor[t.cursor:]...)
			t.cursor--
			t.markDirty()
		}
		return true
	case r == InputDelete:
		if !t.editing {
			t.editing = true
			t.marked = false
			t.editor = nil
			t.cursor = 0
			t.markDirty()
		} else if t.marked {
			t.marked = false
			t.editor = t.editor[:0]
			t.cursor = 0
			t.markDirty()
		} else if t.cursor < len(t.editor) {
			t.editor = append(t.editor[:t.cursor], t.editor[t.cursor+1:]...)
			t.markDirty()
		}
		return true
	case r == InputEscape:
		if t.editing {
			t.stopEdit()
			return true
		}
	case t.converter.allowed(r):
		if !t.editing {
			t.startEdit()
			t.marked = false
		}

		if t.marked && t.converter.allowedInsert(r, nil, 0) {
			t.marked = false
			t.editor = append(t.editor[:0], r)
			t.cursor = 1
			t.markDirty()
		} else if len(t.editor) < t.maxLen && t.converter.allowedInsert(r, t.editor, t.cursor) {
			t.editor = append(t.editor[:t.cursor], append([]rune{r}, t.editor[t.cursor:]...)...)
			t.cursor++
			t.markDirty()
		}

		return true
	}

	return false
}

func (t *textBase) focusLost() {
	if !t.editing {
		return
	}

	t.stopEdit()
}

func (t *textBase) startEdit() {
	s := t.converter.getValueAsStr()
	t.editor = []rune(s)
	t.cursor = len(t.editor)
	t.editing = true
	if len(s) > 0 {
		t.marked = true
	}
	t.markDirty()
}

func (t *textBase) stopEdit() {
	t.editing = false
	t.markDirty()
}
