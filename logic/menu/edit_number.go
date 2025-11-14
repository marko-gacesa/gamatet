// Copyright (c) 2024, 2025 by Marko Gaćeša

package menu

import (
	"strconv"
)

var _ Item = (*Number[int64])(nil)

// Number is menu item that assigns a value to an int64 variable.
type Number[T integer] struct {
	textBase
	ptr      *T
	valueMin T
	valueMax T
}

// NewNumber creates new Number menu item.
func NewNumber[T integer](ptr *T, vmin, vmax T, label, description string, options ...func(Item)) *Number[T] {
	if ptr == nil {
		panic(strNilPointer)
	}
	t := &Number[T]{
		textBase: makeTextBase(20, 20, label, description),
		ptr:      ptr,
		valueMin: vmin,
		valueMax: vmax,
	}
	t.textBase.converter = t
	t.fix()
	applyOptions(t, options...)
	return t
}

func (n *Number[T]) fix() {
	if *n.ptr < n.valueMin {
		*n.ptr = n.valueMin
		n.markDirty()
	}
	if *n.ptr > n.valueMax {
		*n.ptr = n.valueMax
		n.markDirty()
	}
}

func (n *Number[T]) getValueAsStr() string {
	return strconv.FormatInt(int64(*n.ptr), 10)
}

func (n *Number[T]) setValueFromStr(s string) {
	number, _ := strconv.ParseInt(s, 10, 64)
	*n.ptr = T(number)
}

func (*Number[T]) allowed(r rune) bool {
	return r == '-' || r >= '0' && r <= '9'
}

func (n *Number[T]) allowedInsert(r rune, _ []rune, cursor int) bool {
	return r == '-' && cursor == 0 && n.valueMin < 0 || r >= '0' && r <= '9'
}
