// Copyright (c) 2024,2025 by Marko Gaćeša

package menu

var _ Item = (*Bool)(nil)

// Bool is a menu item that manages a boolean variable.
type Bool struct {
	base
	ptr *bool
}

// NewBool creates new Bool menu item.
func NewBool(ptr *bool, label, description string, options ...func(Item)) *Bool {
	if ptr == nil {
		panic("need non-nil pointer")
	}
	b := &Bool{
		base: makeBase(label, description),
		ptr:  ptr,
	}
	applyOptions(b, options...)
	return b
}

func (b *Bool) Text() string {
	if b.current != "" {
		return b.current
	}

	if *b.ptr {
		b.current = b.label + ": " + "ON"
	} else {
		b.current = b.label + ": " + "OFF"
	}

	return b.current
}

func (b *Bool) increase() {
	*b.ptr = !*b.ptr
	b.markDirty()
}

func (b *Bool) decrease() { b.increase() }

func (b *Bool) input(r rune) bool {
	if r == InputEnter {
		b.increase()
		return true
	}
	return false
}
