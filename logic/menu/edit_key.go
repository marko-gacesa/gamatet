// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package menu

var _ Item = (*Key[byte])(nil)

// Key is menu item that listens to keystrokes.
type Key[T any] struct {
	base
	ptr       *T
	listening bool
	stringFn  func(T) string
	convertFn func(byte) T
}

// NewKey creates new Key menu item.
func NewKey[T any](
	ptr *T,
	label, description string,
	stringFn func(T) string,
	convertFn func(byte) T,
	options ...func(Item),
) *Key[T] {
	if ptr == nil {
		panic(strNilPointer)
	}
	k := &Key[T]{
		base:      makeBase(label, description),
		ptr:       ptr,
		stringFn:  stringFn,
		convertFn: convertFn,
	}
	applyOptions(k, options...)
	return k
}

func (k *Key[T]) Text() string {
	if k.current != "" {
		return k.current
	}

	if k.listening {
		k.current = k.getLabel() + ": " + string(cursor)
	} else {
		k.current = k.getLabel() + ": " + k.stringFn(*k.ptr)
	}

	return k.current
}

func (k *Key[T]) input(r rune) bool {
	if r == InputEnter && !k.listening {
		k.listening = true
		k.markDirty()
		return true
	}
	return false
}

func (k *Key[T]) inputKey(key byte) bool {
	if !k.listening {
		return false
	}

	*k.ptr = k.convertFn(key)
	k.listening = false
	k.markDirty()
	return true
}
