// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package loop

import (
	"log/slog"

	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/marko-gacesa/gamatet/logic/gamepad"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func connectGamepad(joy glfw.Joystick, log *slog.Logger) {
	number := gamepad.Connect(int(joy), joy.GetGamepadName(), joy.GetGUID())
	if number >= 0 {
		log.Info("Joystick",
			"name", joy.GetGamepadName(),
			"state", "connected",
			"gamepad_number", number)
	}
}

func disconnectGamepad(joy glfw.Joystick, log *slog.Logger) {
	number, name, _ := gamepad.Disconnect(int(joy))
	if number >= 0 {
		log.Info("Joystick",
			"name", name,
			"state", "disconnected",
			"gamepad_number", number)
	}
}

var _gamepadButtonMap = map[glfw.GamepadButton]gamepad.Button{
	glfw.ButtonY: gamepad.ButtonPadUp,
	glfw.ButtonA: gamepad.ButtonPadDown,
	glfw.ButtonX: gamepad.ButtonPadLeft,
	glfw.ButtonB: gamepad.ButtonPadRight,

	glfw.ButtonBack:        gamepad.ButtonBack,
	glfw.ButtonStart:       gamepad.ButtonForward,
	glfw.ButtonLeftBumper:  gamepad.ButtonLeftBump,
	glfw.ButtonRightBumper: gamepad.ButtonRightBump,

	glfw.ButtonDpadUp:    gamepad.ButtonUp,
	glfw.ButtonDpadDown:  gamepad.ButtonDown,
	glfw.ButtonDpadLeft:  gamepad.ButtonLeft,
	glfw.ButtonDpadRight: gamepad.ButtonRight,
}

func processGamepad(gamepadIdx int, scr screen.Screen) {
	g := &gamepad.Gamepads[gamepadIdx]

	if !g.Connected {
		return
	}

	joy := glfw.Joystick(g.Handle)
	if joy < glfw.Joystick1 || joy > glfw.JoystickLast {
		return
	}

	state := joy.GetGamepadState()
	if state == nil {
		return
	}

	var changesBuffer [2]gamepad.ButtonChange
	for _, act := range g.Axes.Changes(gamepad.Axis(glfw.AxisLeftX),
		gamepad.ButtonLeft, gamepad.ButtonRight, changesBuffer[:0], state.Axes[glfw.AxisLeftX]) {
		scr.InputGamepadPress(gamepadIdx, act)
	}
	for _, act := range g.Axes.Changes(gamepad.Axis(glfw.AxisLeftY),
		gamepad.ButtonUp, gamepad.ButtonDown, changesBuffer[:0], state.Axes[glfw.AxisLeftY]) {
		scr.InputGamepadPress(gamepadIdx, act)
	}

	for glButton, button := range _gamepadButtonMap {
		change := g.Buttons.Change(gamepad.Key(glButton), state.Buttons[glButton] == glfw.Press)
		if change != gamepad.Nothing {
			scr.InputGamepadPress(gamepadIdx, gamepad.ButtonChange{Button: button, Change: change})
		}
	}
}
