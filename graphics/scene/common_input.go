// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package scene

import (
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/marko-gacesa/gamatet/game/action"
	"github.com/marko-gacesa/gamatet/graphics/scene/base"
	"github.com/marko-gacesa/gamatet/internal/config/key"
	"github.com/marko-gacesa/gamatet/logic/gamepad"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func inputKeyboardGlobal(
	key int,
	act screen.KeyAction,
	actionCh chan<- action.Action,
) {
	if act == screen.KeyActionPress {
		switch glfw.Key(key) {
		case glfw.KeyEscape:
			actionCh <- action.Abort
		case glfw.KeyPause:
			actionCh <- action.Pause
		}
	}
}

func inputKeyboardPlayer(
	key int,
	act screen.KeyAction,
	keys key.Input,
	waitDoneCh <-chan struct{},
	playerInCh chan<- []byte,
) {
	if act == screen.KeyActionPress {
		switch glfw.Key(key) {
		case KeyMap[keys.Left]:
			base.SendAction(action.MoveLeft, waitDoneCh, playerInCh)
		case KeyMap[keys.Right]:
			base.SendAction(action.MoveRight, waitDoneCh, playerInCh)
		case KeyMap[keys.Activate]:
			base.SendAction(action.Activate, waitDoneCh, playerInCh)
		case KeyMap[keys.Boost]:
			base.SendAction(action.SpeedUp, waitDoneCh, playerInCh)
		case KeyMap[keys.Drop]:
			base.SendAction(action.Drop, waitDoneCh, playerInCh)
		}
	} else if act == screen.KeyActionRelease {
		switch glfw.Key(key) {
		case KeyMap[keys.Boost]:
			base.SendAction(action.SpeedDown, waitDoneCh, playerInCh)
		}
	}
}

func inputGamepad(
	b gamepad.ButtonChange,
	actionCh chan<- action.Action,
	waitDoneCh <-chan struct{},
	playerInCh chan<- []byte,
) {
	if b.Change == gamepad.Press {
		switch b.Button {
		case gamepad.ButtonBack:
			actionCh <- action.Abort
		case gamepad.ButtonLeftBump, gamepad.ButtonRightBump:
			actionCh <- action.Pause
		case gamepad.ButtonLeft:
			base.SendAction(action.MoveLeft, waitDoneCh, playerInCh)
		case gamepad.ButtonRight:
			base.SendAction(action.MoveRight, waitDoneCh, playerInCh)
		case gamepad.ButtonUp, gamepad.ButtonPadLeft, gamepad.ButtonPadRight, gamepad.ButtonPadUp:
			base.SendAction(action.Activate, waitDoneCh, playerInCh)
		case gamepad.ButtonDown:
			base.SendAction(action.SpeedUp, waitDoneCh, playerInCh)
		case gamepad.ButtonPadDown:
			base.SendAction(action.Drop, waitDoneCh, playerInCh)
		}
	} else if b.Change == gamepad.Release {
		switch b.Button {
		case gamepad.ButtonDown:
			base.SendAction(action.SpeedDown, waitDoneCh, playerInCh)
		}
	}
}
