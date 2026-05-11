// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package gamepad

type Gamepad struct {
	Connected bool
	Handle    int
	Name      string
	Ident     string
	Buttons   Buttons
	Axes      Axes
}

const Count = 4

var Gamepads [Count]Gamepad
var Indices []int

type Button int8

const (
	ButtonBack Button = iota
	ButtonForward

	ButtonLeftBump
	ButtonRightBump

	// Movement buttons.

	ButtonUp
	ButtonDown
	ButtonLeft
	ButtonRight

	// The pad buttons are like 4 universal actions on a gamepad.

	ButtonPadUp    // Y or Triangle
	ButtonPadDown  // A or Cross
	ButtonPadLeft  // X or Square
	ButtonPadRight // B or Circle
)

type Change int8

const (
	Nothing Change = 0
	Press   Change = 1
	Release Change = -1
)

type ButtonChange struct {
	Button Button
	Change Change
}

func init() {
	for idx := range Gamepads {
		clear(idx)
	}

	Indices = make([]int, Count)
	for i := range Indices {
		Indices[i] = i
	}
}

func clear(idx int) {
	Gamepads[idx].Connected = false
	Gamepads[idx].Handle = -1
	Gamepads[idx].Name = ""
	Gamepads[idx].Ident = ""
	Gamepads[idx].Buttons.stateMap = nil
	Gamepads[idx].Axes.stateMap = nil
}

func Connect(handle int, name, ident string) int {
	for idx := range Gamepads {
		if Gamepads[idx].Connected {
			continue
		}

		Gamepads[idx].Connected = true
		Gamepads[idx].Handle = handle
		Gamepads[idx].Name = name
		Gamepads[idx].Ident = ident
		Gamepads[idx].Buttons = makeButtons()
		Gamepads[idx].Axes = makeAxes()

		return idx
	}

	return -1
}

func Disconnect(handle int) (int, string, string) {
	for idx := range Gamepads {
		if !Gamepads[idx].Connected || Gamepads[idx].Handle != handle {
			continue
		}

		name := Gamepads[idx].Name
		ident := Gamepads[idx].Ident

		clear(idx)

		return idx, name, ident
	}

	return -1, "", ""
}
