// Copyright (c) 2024, 2025 by Marko Gaćeša

package menu

var _ Item = (*Bool)(nil)

// Bool is a menu item that manages a boolean variable.
type Bool struct {
	base
	ptr       *bool
	strValues [2]string
}

// NewBool creates new Bool menu item.
func NewBool(ptr *bool, label, description string, options ...func(Item)) *Bool {
	if ptr == nil {
		panic(strNilPointer)
	}
	b := &Bool{
		base:      makeBase(label, description),
		ptr:       ptr,
		strValues: [2]string{"FALSE", "TRUE"},
	}
	applyOptions(b, options...)
	return b
}

func (b *Bool) Text() string {
	if b.current != "" {
		return b.current
	}

	if *b.ptr {
		b.current = b.getLabel() + ": " + b.strValues[1]
	} else {
		b.current = b.getLabel() + ": " + b.strValues[0]
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
