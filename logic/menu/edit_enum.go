// Copyright (c) 2024,2025 by Marko Gaćeša

package menu

import (
	"fmt"
	"slices"
)

var _ Item = (*Enum[int])(nil)

// Enum is a menu item that manages a variable of any type.
// The variable can have only one value from the provided list.
type Enum[T comparable] struct {
	base
	ptr    *T
	values []T
}

// NewEnum creates new Enum menu item.
func NewEnum[T comparable](ptr *T, values []T, label, description string, options ...func(Item)) *Enum[T] {
	if ptr == nil {
		panic("need non-nil pointer")
	}
	if len(values) == 0 {
		panic("no enum values provided")
	}
	e := &Enum[T]{
		base:   makeBase(label, description),
		ptr:    ptr,
		values: values,
	}
	e.fix()
	applyOptions(e, options...)
	return e
}

func (e *Enum[T]) Text() string {
	if e.current != "" {
		return e.current
	}

	e.current = fmt.Sprintf("%s: ‹%v›", e.label, *e.ptr)
	return e.current
}

func (e *Enum[T]) fix() {
	if slices.Index(e.values, *e.ptr) >= 0 {
		return
	}
	*e.ptr = e.values[0]
	e.markDirty()
}

func (e *Enum[T]) increase() {
	idx := slices.Index(e.values, *e.ptr)
	if idx < 0 {
		*e.ptr = e.values[0]
	} else {
		n := len(e.values)
		*e.ptr = e.values[(idx+1)%n]
	}
	e.markDirty()
}

func (e *Enum[T]) decrease() {
	idx := slices.Index(e.values, *e.ptr)
	if idx < 0 {
		*e.ptr = e.values[0]
	} else {
		n := len(e.values)
		*e.ptr = e.values[(idx-1+n)%n]
	}
	e.markDirty()
}

func (e *Enum[T]) input(r rune) bool {
	if r == InputEnter {
		e.increase()
		return true
	}
	return false
}
