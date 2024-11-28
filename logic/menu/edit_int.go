// Copyright (c) 2024 by Marko Gaćeša

package menu

import "strconv"

type Integer struct {
	editable
	value    *int
	valueMin int
	valueMax int
}

func NewInteger(v *int, vmin, vmax int, label, description string) *Integer {
	if v == nil {
		panic("need non-nil pointer")
	}
	if vmin > vmax {
		panic("invalid integer limits provided")
	}
	if *v < vmin {
		*v = vmin
	}
	if *v > vmax {
		*v = vmax
	}
	return &Integer{
		editable: editable{
			base:  base{description: description},
			label: label,
		},
		value:    v,
		valueMin: vmin,
		valueMax: vmax,
	}
}

func (i *Integer) Text() string {
	if i.current != "" {
		return i.current
	}

	i.current = i.label + ": " + strconv.Itoa(*i.value)
	return i.current
}

func (i *Integer) Increase() {
	if *i.value < i.valueMax {
		*i.value++
		i.dirty()
	}
}

func (i *Integer) Decrease() {
	if *i.value > i.valueMin {
		*i.value--
		i.dirty()
	}
}

func (i *Integer) Input(r rune) bool {
	if r == InputEnter {
		i.Increase()
		return true
	}
	return false
}
