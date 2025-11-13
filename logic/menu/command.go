// Copyright (c) 2024, 2025 by Marko Gaćeša

package menu

var _ Item = (*Command[string])(nil)

// Command is a menu item that activates when the enter key is pressed while the item is active.
type Command[T any] struct {
	base
	ptr   *T
	value T
}

// NewCommand creates a new Command menu item that activates when the enter key is pressed while the item is active.
// The activation means write of the provided value to the provided pointer.
func NewCommand[T any](ptr *T, value T, label, description string, options ...func(Item)) *Command[T] {
	if ptr == nil {
		panic(strNilPointer)
	}
	cmd := &Command[T]{
		base:  makeBase(label, description),
		ptr:   ptr,
		value: value,
	}
	applyOptions(cmd, options...)
	return cmd
}

func (c *Command[T]) input(r rune) bool {
	if r != InputEnter {
		return false
	}
	*c.ptr = c.value
	return true
}
