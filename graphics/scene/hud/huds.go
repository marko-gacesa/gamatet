// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package hud

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/marko-gacesa/gamatet/graphics/render"
	"github.com/marko-gacesa/gamatet/graphics/texture"
	"github.com/marko-gacesa/gamatet/internal/values"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

const (
	PosFPS       = render.HUDUpperRight
	PosLatencies = render.HUDUpperLeft
	PosAuthor    = render.HUDLowerRight
)

type HUDs struct {
	textHUD *render.Text
	huds    []*render.HUD
}

func NewHUDs(tex *texture.Manager) *HUDs {
	textHUD := render.MakeText(tex, render.HudFont)

	h := &HUDs{
		textHUD: textHUD,
		huds:    make([]*render.HUD, 0, 2),
	}

	h.Add(String(fmt.Sprintf("%v by Marko Gaćeša", values.ProgramName)), PosAuthor)

	return h
}

func (h *HUDs) Release() {
	h.textHUD.Release()
}

func (h *HUDs) Add(v fmt.Stringer, pos render.HUDPos) {
	h.huds = append(h.huds, render.NewHUD(v, pos, h.textHUD))
}

func (h *HUDs) InputKeyPress(key int, act screen.KeyAction) {
	k := glfw.Key(key)
	if k == glfw.KeyF2 && act == screen.KeyActionPress {
		for _, hud := range h.huds {
			hud.ShowToggle()
		}
	}
}

func (h *HUDs) Prepare(viewW, viewH int) {
	for _, hud := range h.huds {
		hud.Prepare(viewW, viewH)
	}
}

func (h *HUDs) Render(r *render.Renderer) {
	for _, hud := range h.huds {
		hud.Render(r)
	}
}

type String string

func (s String) String() string { return string(s) }
