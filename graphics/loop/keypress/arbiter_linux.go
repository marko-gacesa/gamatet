// Copyright (c) 2023-2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

//go:build (linux && !wayland) || (freebsd && !wayland) || (netbsd && !wayland) || (openbsd && !wayland)

package keypress

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

func NewArbiter() Arbiter {
	return arbiterX11{
		events: make(map[glfw.Key][]entry),
	}
}
