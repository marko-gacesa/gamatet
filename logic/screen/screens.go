// Copyright (c) 2025 by Marko Gaćeša

package screen

import (
	"time"
)

type Screens []Screen

var _ Screen = (*Screens)(nil)

func (screens Screens) UpdateViewSize(w, h int) {
	for _, s := range screens {
		s.UpdateViewSize(w, h)
	}
}

func (screens Screens) Release() {
	for _, s := range screens {
		s.Release()
	}
}

func (screens Screens) InputKeyPress(key int, act KeyAction) {
	for _, s := range screens {
		s.InputKeyPress(key, act)
	}
}

func (screens Screens) InputChar(char rune) {
	for _, s := range screens {
		s.InputChar(char)
	}
}

func (screens Screens) Prepare(now time.Time) {
	for _, s := range screens {
		s.Prepare(now)
	}
}

func (screens Screens) Render() {
	for _, s := range screens {
		s.Render()
	}
}
