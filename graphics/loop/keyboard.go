// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package loop

import (
	"github.com/marko-gacesa/gamatet/graphics/loop/keypress"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func processKeyboard(scr screen.Screen, keyArbiter keypress.Arbiter) {
	var buffer [8]keypress.KeyEvent

	keyEvents := keyArbiter.Events(buffer[:0])
	for _, keyEvent := range keyEvents {
		scr.InputKeyPress(int(keyEvent.Key), keypress.ConvertAction(keyEvent.Action))
	}
}
