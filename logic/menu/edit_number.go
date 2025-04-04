// Copyright (c) 2024,2025 by Marko Gaćeša

package menu

import (
	"strconv"
)

var _ Item = (*Number)(nil)

// Number is menu item that assigns a value to an int64 variable.
type Number struct {
	textBase
	ptr      *int64
	valueMin int64
	valueMax int64
}

// NewNumber creates new Number menu item.
func NewNumber(ptr *int64, vmin, vmax int64, label, description string, options ...func(Item)) *Number {
	if ptr == nil {
		panic("need non-nil pointer")
	}
	t := &Number{
		textBase: makeTextBase(20, label, description),
		ptr:      ptr,
		valueMin: vmin,
		valueMax: vmax,
	}
	t.textBase.converter = t
	t.fix()
	applyOptions(t, options...)
	return t
}

func (n *Number) fix() {
	if *n.ptr < n.valueMin {
		*n.ptr = n.valueMin
		n.markDirty()
	}
	if *n.ptr > n.valueMax {
		*n.ptr = n.valueMax
		n.markDirty()
	}
}

func (n *Number) getValueAsStr() string {
	return strconv.FormatInt(*n.ptr, 10)
}

func (n *Number) setValueFromStr(s string) {
	number, _ := strconv.ParseInt(s, 10, 64)
	*n.ptr = number
}

func (*Number) allowed(r rune) bool {
	return r == '-' || r >= '0' && r <= '9'
}

func (n *Number) allowedInsert(r rune, _ []rune, cursor int) bool {
	return r != '-' || r == '-' && cursor == 0 && n.valueMin < 0
}
