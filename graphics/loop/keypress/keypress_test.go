// Copyright (c) 2023-2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package keypress

import (
	"slices"
	"testing"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func TestArbiterX11(t *testing.T) {
	x := &arbiterX11{
		events: make(map[glfw.Key][]entry),
	}

	x.Process(glfw.Key0, 0, glfw.Press, 0)
	if want, got := []KeyEvent{{Key: glfw.Key0, Scancode: 0, Action: glfw.Press, Mods: 0}}, x.Events(nil); !slices.Equal(want, got) {
		t.Errorf("1) want=%v got=%v", want, got)
	}

	x.Process(glfw.Key0, 0, glfw.Release, 0)
	if want, got := []KeyEvent(nil), x.Events(nil); !slices.Equal(want, got) {
		t.Errorf("2) want=%v got=%v", want, got)
	}

	if want, got := []KeyEvent{{Key: glfw.Key0, Scancode: 0, Action: glfw.Release, Mods: 0}}, x.Events(nil); !slices.Equal(want, got) {
		t.Errorf("3) want=%v got=%v", want, got)
	}

	x.Process(glfw.Key1, 0, glfw.Press, 0)
	if want, got := []KeyEvent{{Key: glfw.Key1, Scancode: 0, Action: glfw.Press, Mods: 0}}, x.Events(nil); !slices.Equal(want, got) {
		t.Errorf("4) want=%v got=%v", want, got)
	}

	x.Process(glfw.Key1, 0, glfw.Release, 0)
	x.Process(glfw.Key1, 0, glfw.Press, 0)
	if want, got := []KeyEvent{{Key: glfw.Key1, Scancode: 0, Action: glfw.Repeat, Mods: 0}}, x.Events(nil); !slices.Equal(want, got) {
		t.Errorf("5) want=%v got=%v", want, got)
	}

	x.Process(glfw.Key1, 0, glfw.Release, 0)
	x.Process(glfw.Key1, 0, glfw.Press, 0)
	x.Process(glfw.Key1, 0, glfw.Release, 0)
	if want, got := []KeyEvent{{Key: glfw.Key1, Scancode: 0, Action: glfw.Repeat, Mods: 0}}, x.Events(nil); !slices.Equal(want, got) {
		t.Errorf("6) want=%v got=%v", want, got)
	}

	if want, got := []KeyEvent{{Key: glfw.Key1, Scancode: 0, Action: glfw.Release, Mods: 0}}, x.Events(nil); !slices.Equal(want, got) {
		t.Errorf("7) want=%v got=%v", want, got)
	}
}
