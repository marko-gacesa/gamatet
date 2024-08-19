// Copyright (c) 2024 by Marko Gaćeša

package render

import (
	"gamatet/graphics/geometry"
	"gamatet/graphics/material"
	"gamatet/graphics/runeatlas"
	"gamatet/graphics/texture"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang/freetype/truetype"
	"time"
)

type Text struct {
	tex        uint32
	texManager *texture.Manager
	atlas      runeatlas.RuneAtlas
	geom       geometry.Geometry
	mat        material.Text
}

func MakeText(manager *texture.Manager, font *truetype.Font) *Text {
	fontFace := runeatlas.NewFace(font, 32, 72)
	runeAtlas := runeatlas.NewRuneAtlas(fontFace, 512)

	runeAtlas.Store('|')
	for key := '0'; key <= '9'; key++ {
		runeAtlas.Store(key)
	}

	tex := manager.Bind(runeAtlas.Image())
	runeAtlas.ClearDirty()

	return &Text{
		tex:        tex,
		texManager: manager,
		atlas:      *runeAtlas,
		geom:       geometry.NewSquare(),
		mat:        *material.NewText(tex),
	}
}

func (t *Text) Release() {
	t.atlas.Clear()

	t.texManager.Delete(t.tex)
	t.tex = 0

	t.geom.Delete()
	t.mat.Delete()
}

func (t *Text) Material(r *Renderer, color mgl32.Vec4, ch rune) (runeatlas.RectUV, bool) {
	runeRect, ok := t.atlas.TextUV(ch)
	if !ok {
		return runeatlas.RectUV{}, false
	}

	mat := &t.mat

	r.Material(mat)
	mat.Texture(t.tex)
	mat.Color(color)
	mat.TexUV(runeRect)

	if t.atlas.IsDirty() {
		t.texManager.ReBind(t.tex, t.atlas.Image())
		t.atlas.ClearDirty()
	}

	return runeRect, true
}

func (t *Text) Rune(r *Renderer, model mgl32.Mat4, color mgl32.Vec4, ch rune) {
	runeRect, ok := t.Material(r, color, ch)
	if !ok {
		return
	}

	r.Geometry(t.geom)

	w2h := runeRect.WidthToHeight()
	modelChar := model.Mul4(mgl32.Scale3D(w2h, 1, 1))

	gl.DepthMask(false)
	r.Render(&modelChar)
	gl.DepthMask(true)
}

func (t *Text) String(r *Renderer, model mgl32.Mat4, color mgl32.Vec4, s string) {
	for _, ch := range s {
		if ch > 32 {
			t.atlas.Store(ch)
		}
	}

	gl.DepthMask(false)
	defer gl.DepthMask(true)

	mat := &t.mat

	r.Geometry(t.geom)
	r.Material(mat)
	mat.Texture(t.tex)
	mat.Color(color)

	if t.atlas.IsDirty() {
		t.texManager.ReBind(t.tex, t.atlas.Image())
		t.atlas.ClearDirty()
	}

	modelText := model
	var chPrev rune
	for _, ch := range s {
		var control bool

		switch ch {
		case ' ':
			modelText = modelText.Mul4(mgl32.Translate3D(0.2, 0, 0))
			continue
		case '\n':
			model = model.Mul4(mgl32.Translate3D(0, -1, 0))
			modelText = model
			continue
		case '\x01': // cursor
			if time.Now().UnixMilli()%500 < 250 {
				continue
			}
			control = true
			ch = '|'
		}

		runeRect, ok := t.atlas.TextUV(ch)
		if !ok {
			continue
		}

		mat.TexUV(runeRect)
		w2h := runeRect.WidthToHeight()
		k2h := t.atlas.KernToHeight(chPrev, ch)

		if control {
			modelChar := modelText.Mul4(mgl32.Scale3D(w2h, 1, 1))
			mat.Color(color.Mul(0.8))
			r.Render(&modelChar)
			mat.Color(color)
			continue
		}

		w2h2 := w2h / 2
		modelText = modelText.Mul4(mgl32.Translate3D(w2h2+k2h, 0, 0))
		modelChar := modelText.Mul4(mgl32.Scale3D(w2h, 1, 1))
		r.Render(&modelChar)
		modelText = modelText.Mul4(mgl32.Translate3D(w2h2, 0, 0))

		chPrev = ch
	}
}

func (t *Text) Prepare(strs ...string) {
	for _, s := range strs {
		for _, ch := range s {
			if ch > 32 {
				t.atlas.Store(ch)
			}
		}
	}
}
