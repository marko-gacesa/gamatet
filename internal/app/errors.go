// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import "errors"

var (
	errorInputMissing          = errors.New("internal error: Input missing")
	errorPlayerIndexOutOfRange = errors.New("internal error: Player index out of range")
	errorPresetIndexOutOfRange = errors.New("internal error: Preset index out of range")
)
