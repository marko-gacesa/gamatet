// Copyright (c) 2024, 2025 by Marko Gaćeša

package menu

import (
	"fmt"
	"slices"
	"strings"
)

var _ Item = (*Enum[int])(nil)

// Enum is a menu item that manages a variable of any type.
// The variable can have only one value from the provided list.
type Enum[T comparable] struct {
	base
	ptr     *T
	values  []T
	nameMap map[T]string
}

// NewEnum creates new Enum menu item.
func NewEnum[T comparable](ptr *T, values []T, nameMap map[T]string, label, description string, options ...func(Item)) *Enum[T] {
	if ptr == nil {
		panic(strNilPointer)
	}
	if len(values) == 0 {
		panic("no enum values provided")
	}
	e := &Enum[T]{
		base:    makeBase(label, description),
		ptr:     ptr,
		values:  values,
		nameMap: nameMap,
	}
	e.fix()
	applyOptions(e, options...)
	return e
}

func (e *Enum[T]) Text() string {
	if e.current != "" {
		return e.current
	}

	var text string
	if e.nameMap != nil {
		text = e.nameMap[*e.ptr]
	} else {
		text = fmt.Sprintf("%v", *e.ptr)
	}

	n := len(e.values)
	idx := slices.Index(e.values, *e.ptr)

	sb := strings.Builder{}
	sb.WriteString(e.getLabel())
	sb.WriteString(": ")
	if idx > 0 {
		sb.WriteString("‹ ")
	}
	sb.WriteString(text)
	if idx < n-1 {
		sb.WriteString(" ›")
	}

	e.current = sb.String()
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
