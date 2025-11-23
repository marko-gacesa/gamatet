// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package menu

var _ Item = (*Static)(nil)

// Static is a simple, static menu item.
type Static struct {
	base
	inputHandler func(r rune) bool
}

// NewStatic creates a new Static menu item.
func NewStatic(label, description string, inputHandler func(rune) bool, options ...func(Item)) *Static {
	stat := &Static{
		base:         makeBase(label, description),
		inputHandler: inputHandler,
	}
	applyOptions(stat, options...)
	return stat
}

func (s Static) input(r rune) bool {
	if s.inputHandler == nil {
		return false
	}
	return s.inputHandler(r)
}
