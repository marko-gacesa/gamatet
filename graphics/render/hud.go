// Copyright (c) 2024, 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package render

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

var hudColor = mgl32.Vec4{0.5, 0.5, 0, 0.5}

type HUDPos byte

const (
	HUDLowerLeft HUDPos = iota
	HUDLowerRight
	HUDUpperLeft
	HUDUpperRight
)

type HUD struct {
	v    fmt.Stringer
	pos  HUDPos
	text *Text

	viewW float32
	viewH float32

	show  bool
	value string
	model mgl32.Mat4
}

func NewHUD(v fmt.Stringer, pos HUDPos, text *Text) *HUD {
	return &HUD{v: v, pos: pos, text: text}
}

func (hud *HUD) ShowToggle() {
	hud.show = !hud.show
}

func (hud *HUD) Prepare(viewW, viewH int) {
	if hud.v == nil || !hud.show {
		hud.value = ""
		return
	}

	hud.viewW = float32(viewW)
	hud.viewH = float32(viewH)

	s := hud.v.String()
	hud.value = s
	if s == "" {
		return
	}

	const sideMargin = 0.1
	const runeDim = 48

	switch hud.pos {
	case HUDUpperLeft:
		hud.model = mgl32.Translate3D(-1, 1, 0).
			Mul4(mgl32.Scale3D(runeDim/hud.viewW, runeDim/hud.viewH, 1)).
			Mul4(mgl32.Translate3D(sideMargin, -0.5, 1))

	case HUDUpperRight:
		tw, _ := hud.text.Dim(s)
		hud.model = mgl32.Translate3D(1, 1, 0).
			Mul4(mgl32.Scale3D(runeDim/hud.viewW, runeDim/hud.viewH, 1)).
			Mul4(mgl32.Translate3D(-sideMargin-tw, -0.5, 1))
	case HUDLowerLeft:
		_, th := hud.text.Dim(s)
		hud.model = mgl32.Translate3D(-1, -1, 0).
			Mul4(mgl32.Scale3D(runeDim/hud.viewW, runeDim/hud.viewH, 1)).
			Mul4(mgl32.Translate3D(sideMargin, -0.5+th, 1))
	case HUDLowerRight:
		tw, th := hud.text.Dim(s)
		hud.model = mgl32.Translate3D(1, -1, 0).
			Mul4(mgl32.Scale3D(runeDim/hud.viewW, runeDim/hud.viewH, 1)).
			Mul4(mgl32.Translate3D(-sideMargin-tw, -0.5+th, 1))
	}
}

func (hud *HUD) Render(r *Renderer) {
	if hud.value == "" {
		return
	}

	r.Orthogonal2D(2, 2)
	hud.text.String(r, hud.model, hudColor, hud.value)
}
