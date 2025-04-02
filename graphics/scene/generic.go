// Copyright (c) 2025 by Marko Gaćeša

package scene

import (
	"gamatet/logic/screen"
	"time"
)

type ScreenMap map[string]screen.Screen

func (m ScreenMap) Add(key string, s screen.Screen) {
	m[key] = s
}

func (m ScreenMap) Remove(key string) {
	s, ok := m[key]
	if !ok {
		return
	}
	s.Release()
}

func (m ScreenMap) UpdateViewSize(w, h int) {
	for _, s := range m {
		s.UpdateViewSize(w, h)
	}
}

func (m ScreenMap) Release() {
	for _, s := range m {
		s.Release()
	}
}

func (m ScreenMap) InputKeyPress(key, scancode int) {
	for _, s := range m {
		s.InputKeyPress(key, scancode)
	}
}

func (m ScreenMap) InputChar(char rune) {
	for _, s := range m {
		s.InputChar(char)
	}
}

func (m ScreenMap) Prepare(now time.Time) {
	for _, s := range m {
		s.Prepare(now)
	}
}

func (m ScreenMap) Render() {
	for _, s := range m {
		s.Render()
	}
}
