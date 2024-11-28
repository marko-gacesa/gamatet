// Copyright (c) 2024 by Marko Gaćeša

package menu

type Bool struct {
	editable
	value *bool
}

func NewBool(b *bool, label, description string) *Bool {
	if b == nil {
		panic("need non-nil pointer")
	}
	return &Bool{
		editable: editable{
			base:  base{description: description},
			label: label,
		},
		value: b,
	}
}

func (b *Bool) Text() string {
	if b.current != "" {
		return b.current
	}

	if *b.value {
		b.current = b.label + ": " + "ON"
	} else {
		b.current = b.label + ": " + "OFF"
	}

	return b.current
}

func (b *Bool) Increase() {
	*b.value = !*b.value
	b.dirty()
}

func (b *Bool) Decrease() { b.Increase() }

func (b *Bool) Input(r rune) bool {
	if r == InputEnter {
		b.Increase()
		return true
	}
	return false
}
