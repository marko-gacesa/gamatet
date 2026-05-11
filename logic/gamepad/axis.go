// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package gamepad

type Axes struct {
	stateMap map[Axis]float32
}

type Axis int

func makeAxes() Axes {
	m := make(map[Axis]float32, 8)
	return Axes{stateMap: m}
}

// Changes converts gamepad axis state to a slice of game actions.
// If the axis moves to/from the negative side, a Change for the "neg" Button would be produced.
// if the axis moves to/from the positive side, a Change for the "pos" Button would be produced.
func (b Axes) Changes(axis Axis, neg, pos Button, buffer []ButtonChange, state float32) []ButtonChange {
	oldState := b.stateMap[axis]

	const threshold = 0.75

	if oldState <= -threshold && state > -threshold {
		buffer = append(buffer, ButtonChange{Button: neg, Change: Release})
	}
	if oldState >= threshold && state < threshold {
		buffer = append(buffer, ButtonChange{Button: pos, Change: Release})
	}

	if oldState > -threshold && state <= -threshold {
		buffer = append(buffer, ButtonChange{Button: neg, Change: Press})
	}
	if oldState < threshold && state >= threshold {
		buffer = append(buffer, ButtonChange{Button: pos, Change: Press})
	}

	b.stateMap[axis] = state

	return buffer
}
