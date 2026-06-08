// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package loop

import (
	"github.com/marko-gacesa/gamatet/graphics/loop/keypress"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

var keyboardBuffer = make([]keypress.KeyEvent, 0, 16)

func processKeyboard(scr screen.Screen, keyArbiter keypress.Arbiter) {
	keyboardBuffer = keyboardBuffer[:0]
	keyEvents := keyArbiter.Events(keyboardBuffer)
	for _, keyEvent := range keyEvents {
		scr.InputKeyPress(int(keyEvent.Key), keypress.ConvertAction(keyEvent.Action))
	}
}
