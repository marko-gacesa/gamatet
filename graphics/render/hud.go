// Copyright (c) 2024, 2025 by Marko Gaćeša

package render

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

const hudContentW = 80
const hudContentH = hudContentW * 9 / 16

var hudColor = mgl32.Vec4{0.5, 0.5, 0, 0.7}

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

	value string
	model mgl32.Mat4
}

func NewHUD(v fmt.Stringer, pos HUDPos, text *Text) *HUD {
	return &HUD{v: v, pos: pos, text: text}
}

func (hud *HUD) Prepare() {
	if hud.v == nil {
		return
	}

	s := hud.v.String()
	hud.value = s
	if s == "" {
		return
	}

	const sideMargin = 0.1

	switch hud.pos {
	case HUDUpperLeft:
		hud.model = mgl32.Translate3D(float32(-hudContentW)/2+sideMargin, float32(hudContentH)/2-0.5, 0)
	case HUDUpperRight:
		tw, _ := hud.text.Dim(s)
		hud.model = mgl32.Translate3D(float32(hudContentW)/2-sideMargin-tw, float32(hudContentH)/2-0.5, 0)
	case HUDLowerLeft:
		_, th := hud.text.Dim(s)
		hud.model = mgl32.Translate3D(float32(-hudContentW)/2+sideMargin, -float32(hudContentH)/2-0.5+th, 0)
	case HUDLowerRight:
		tw, th := hud.text.Dim(s)
		hud.model = mgl32.Translate3D(float32(hudContentW)/2-sideMargin-tw, -float32(hudContentH)/2-0.5+th, 0)
	}
}

func (hud *HUD) Render(r *Renderer) {
	if hud.value == "" {
		return
	}

	r.OrthogonalFull(hudContentW, hudContentH, hudContentW, hudContentH, 1)
	hud.text.String(r, hud.model, hudColor, hud.value)
}
