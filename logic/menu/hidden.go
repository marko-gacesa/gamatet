// Copyright (c) 2025 by Marko Gaćeša

package menu

var _ Item = (*Hidden[string])(nil)

// Hidden is a hidden menu item that activates when the escape key is pressed.
type Hidden[T any] struct {
	base
	ptr   *T
	value T
	key   rune
}

// NewHidden creates new Hidden menu item that activates when the a specific key is pressed.
// The activation means write of the provided value to the provided pointer.
func NewHidden[T any](key rune, ptr *T, value T) *Hidden[T] {
	if ptr == nil {
		panic(strNilPointer)
	}
	hidden := &Hidden[T]{
		base:  makeBase("", ""),
		ptr:   ptr,
		value: value,
		key:   key,
	}

	hidden.base.global = true
	hidden.base.visible = false
	hidden.base.visibleFn = func() bool { return false }

	return hidden
}

func (h *Hidden[T]) input(r rune) bool {
	if r != h.key {
		return false
	}
	*h.ptr = h.value
	return true
}
