// Copyright (c) 2023-2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package keypress

import (
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/marko-gacesa/gamatet/logic/screen"
)

type Arbiter interface {
	Process(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)
	Events(list []KeyEvent) []KeyEvent
}

type KeyEvent struct {
	Key      glfw.Key
	Scancode int
	Action   glfw.Action
	Mods     glfw.ModifierKey
}

func ConvertAction(action glfw.Action) screen.KeyAction {
	switch action {
	case glfw.Press:
		return screen.KeyActionPress
	case glfw.Release:
		return screen.KeyActionRelease
	case glfw.Repeat:
		return screen.KeyActionRepeat
	default:
		return screen.KeyActionNothing
	}
}

type arbiterSimple struct {
	events []KeyEvent
}

func (x *arbiterSimple) Process(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	x.events = append(x.events, KeyEvent{Key: key, Scancode: scancode, Action: action, Mods: mods})
}

func (x *arbiterSimple) Events(list []KeyEvent) []KeyEvent {
	list = append(list, x.events...)
	x.events = x.events[:0]
	return list
}

type arbiterX11 struct {
	events map[glfw.Key][]entry
}

type entry struct {
	KeyEvent
	wait bool
}

func (x arbiterX11) Process(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	x.events[key] = append(x.events[key], entry{
		KeyEvent: KeyEvent{Key: key, Scancode: scancode, Action: action, Mods: mods},
		wait:     false,
	})

	for key, events := range x.events {
		for i := 0; i < len(events)-1; i++ {
			if events[i].Action == glfw.Release && events[i+1].Action == glfw.Press {
				events[i].Action = glfw.Repeat
				x.events[key] = append(events[:i+1], events[i+2:]...)
			}
		}
	}
}

func (x arbiterX11) Events(list []KeyEvent) []KeyEvent {
	for key, events := range x.events {
		if lastIdx := len(events) - 1; events[lastIdx].Action == glfw.Release && !events[lastIdx].wait {
			for _, e := range events[:lastIdx] {
				list = append(list, e.KeyEvent)
			}
			events[0] = events[lastIdx]
			events[0].wait = true
			x.events[key] = events[:1]
		} else {
			for _, e := range events {
				list = append(list, e.KeyEvent)
			}
			delete(x.events, key)
		}
	}

	return list
}
