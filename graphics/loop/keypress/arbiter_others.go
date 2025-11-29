// Copyright (c) 2023-2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

//go:build !((linux && !wayland) || (freebsd && !wayland) || (netbsd && !wayland) || (openbsd && !wayland))

package keypress

func NewArbiter() Arbiter {
	return &arbiterSimple{
		events: make([]KeyEvent, 0, 16),
	}
}
