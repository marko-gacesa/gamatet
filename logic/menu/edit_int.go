// Copyright (c) 2024,2025 by Marko Gaćeša

package menu

import "strconv"

var _ Item = (*Integer)(nil)

// Integer is a menu item that manages an integer variable.
type Integer struct {
	base
	old      int
	ptr      *int
	valueMin int
	valueMax int
}

// NewInteger creates new Integer menu item.
func NewInteger(ptr *int, vmin, vmax int, label, description string, options ...func(Item)) *Integer {
	if ptr == nil {
		panic("need non-nil pointer")
	}
	if vmin > vmax {
		panic("invalid integer limits provided")
	}
	i := &Integer{
		base:     makeBase(label, description),
		old:      *ptr,
		ptr:      ptr,
		valueMin: vmin,
		valueMax: vmax,
	}
	i.fix()
	applyOptions(i, options...)
	return i
}

func (i *Integer) Text() string {
	if i.current != "" {
		return i.current
	}

	i.current = i.label + ": " + strconv.Itoa(*i.ptr)
	return i.current
}

func (i *Integer) fix() {
	if *i.ptr < i.valueMin {
		*i.ptr = i.valueMin
		i.markDirty()
	}
	if *i.ptr > i.valueMax {
		*i.ptr = i.valueMax
		i.markDirty()
	}
}

func (i *Integer) increase() {
	if *i.ptr < i.valueMax {
		*i.ptr++
		i.old = *i.ptr
		i.markDirty()
	}
}

func (i *Integer) decrease() {
	if *i.ptr > i.valueMin {
		*i.ptr--
		i.old = *i.ptr
		i.markDirty()
	}
}

func (i *Integer) input(r rune) bool {
	if r == InputEnter {
		i.increase()
		return true
	}
	return false
}
