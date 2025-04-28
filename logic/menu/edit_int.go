// Copyright (c) 2024, 2025 by Marko Gaćeša

package menu

import (
	"fmt"
)

var _ Item = (*Integer[int])(nil)

// Integer is a menu item that manages an integer variable.
type Integer[T integer] struct {
	base
	old      T
	ptr      *T
	valueMin T
	valueMax T
}

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// NewInteger creates new Integer menu item.
func NewInteger[T integer](ptr *T, vmin, vmax T, label, description string, options ...func(Item)) *Integer[T] {
	if ptr == nil {
		panic("need non-nil pointer")
	}
	if vmin > vmax {
		panic("invalid integer limits provided")
	}
	i := &Integer[T]{
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

func (i *Integer[T]) Text() string {
	if i.current != "" {
		return i.current
	}

	i.current = i.getLabel() + ": " + fmt.Sprint(*i.ptr)
	return i.current
}

func (i *Integer[T]) fix() {
	if *i.ptr < i.valueMin {
		*i.ptr = i.valueMin
		i.markDirty()
	}
	if *i.ptr > i.valueMax {
		*i.ptr = i.valueMax
		i.markDirty()
	}
}

func (i *Integer[T]) increase() {
	if *i.ptr < i.valueMax {
		*i.ptr++
		i.old = *i.ptr
		i.markDirty()
	}
}

func (i *Integer[T]) decrease() {
	if *i.ptr > i.valueMin {
		*i.ptr--
		i.old = *i.ptr
		i.markDirty()
	}
}

func (i *Integer[T]) input(r rune) bool {
	if r == InputEnter {
		i.increase()
		return true
	}
	return false
}
