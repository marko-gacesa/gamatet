// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package gamepad

type Buttons struct {
	stateMap map[Key]struct{}
}

type Key int

func makeButtons() Buttons {
	m := make(map[Key]struct{}, 8)
	return Buttons{stateMap: m}
}

// Change returns Press if gamepad key has just been pressed, Release if gamepad key has just been released
// and Nothing if there is no change in the key's state.
func (b Buttons) Change(key Key, state bool) Change {
	_, pressed := b.stateMap[key]

	if state {
		if !pressed {
			b.stateMap[key] = struct{}{}
			return Press
		}
	} else {
		if pressed {
			delete(b.stateMap, key)
			return Release
		}
	}

	return Nothing
}
