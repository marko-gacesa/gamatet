// Copyright (c) 2025 by Marko Gaćeša

package menu

var _ Item = (*Cancel[string])(nil)

// Cancel is a hidden menu item that activates when the escape key is pressed.
type Cancel[T any] struct {
	base
	ptr   *T
	value T
}

// NewCancel creates new Cancel hidden menu item that activates when the escape key is pressed.
// The activation means write of the provided value to the provided pointer.
func NewCancel[T any](ptr *T, value T) *Cancel[T] {
	if ptr == nil {
		panic("need non-nil pointer")
	}
	cancel := &Cancel[T]{
		base:  makeBase("", ""),
		ptr:   ptr,
		value: value,
	}

	cancel.canceler = true
	cancel.base.visible = false
	cancel.base.visibleFn = func() bool { return false }

	return cancel
}

func (c *Cancel[T]) input(r rune) bool {
	if r != InputEscape {
		return false
	}
	*c.ptr = c.value
	return true
}
