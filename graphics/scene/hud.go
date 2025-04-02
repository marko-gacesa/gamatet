// Copyright (c) 2024,2025 by Marko Gaćeša

package scene

import (
	"context"
	"gamatet/graphics/render"
	"gamatet/graphics/scene/base"
	"gamatet/graphics/texture"
	"github.com/go-gl/mathgl/mgl32"
	"strconv"
	"strings"
)

type Hud struct {
	base.Base
	text *render.Text
	fps  *render.FPS
}

func NewHud(
	renderer *render.Renderer,
	tex *texture.Manager,
) *Hud {
	text := render.MakeText(tex, render.HudFont)
	return &Hud{
		Base: base.NewBase(renderer, tex),
		text: text,
		fps:  render.NewFPS(),
	}
}

func (h *Hud) Release() {
	h.text.Release()
}

func (h *Hud) Render(context.Context) {
	v := h.fps.Get()
	if v == 0 {
		return
	}

	sb := strings.Builder{}
	sb.WriteString("fps=")
	sb.WriteString(strconv.Itoa(v))
	s := sb.String()

	//tw, th := h.text.Dim(s)
	//_, _ = tw, th

	const contentW = 50
	const contentH = contentW * 9 / 16
	h.Renderer().OrthogonalFull(contentW, contentH, contentW, contentH, 1)

	model := mgl32.Translate3D(float32(-contentW)/2, float32(contentH)/2-0.5, 0)
	h.text.String(h.Renderer(), model, mgl32.Vec4{0.5, 0.5, 0, 0.7}, s)
}
