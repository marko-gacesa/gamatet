// Copyright (c) 2024 by Marko Gaćeša

package menu

import "fmt"

type Enum[T comparable] struct {
	editable
	value  *T
	values []T
}

func NewEnum[T comparable](v *T, values []T, label, description string) *Enum[T] {
	if v == nil {
		panic("need non-nil pointer")
	}
	if len(values) == 0 {
		panic("no enum values provided")
	}
	if index(*v, values) < 0 {
		*v = values[0]
	}
	return &Enum[T]{
		editable: editable{
			base:  base{description: description},
			label: label,
		},
		value:  v,
		values: values,
	}
}

func (e *Enum[T]) Text() string {
	if e.current != "" {
		return e.current
	}

	e.current = fmt.Sprintf("%s: %v", e.label, *e.value)
	return e.current
}

func (e *Enum[T]) Increase() {
	idx := index(*e.value, e.values)
	if idx < 0 {
		*e.value = e.values[0]
	} else {
		n := len(e.values)
		*e.value = e.values[(idx+1)%n]
	}
	e.dirty()
}

func (e *Enum[T]) Decrease() {
	idx := index(*e.value, e.values)
	if idx < 0 {
		*e.value = e.values[0]
	} else {
		n := len(e.values)
		*e.value = e.values[(idx-1+n)%n]
	}
	e.dirty()
}

func index[T comparable](v T, values []T) int {
	for i := 0; i < len(values); i++ {
		if v == values[i] {
			return i
		}
	}

	return -1
}
