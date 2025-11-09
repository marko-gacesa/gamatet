// Copyright (c) 2024, 2025 by Marko Gaćeša

package render

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang/freetype/truetype"
	"github.com/marko-gacesa/gamatet/graphics/geometry"
	"github.com/marko-gacesa/gamatet/graphics/material"
	"github.com/marko-gacesa/gamatet/graphics/runeatlas"
	"github.com/marko-gacesa/gamatet/graphics/texture"
)

type TextBlock struct {
	tex        uint32
	texManager *texture.Manager
	atlas      runeatlas.RuneAtlas
	geom       geometry.Geometry
	mat        material.TextBlock
}

func MakeTextBlock(manager *texture.Manager, font *truetype.Font) *TextBlock {
	fontFace := runeatlas.NewFace(font, 48, 72)
	runeAtlas := runeatlas.NewRuneAtlas(fontFace, 128)

	tex := manager.Bind(runeAtlas.Image())
	runeAtlas.ClearDirty()

	return &TextBlock{
		tex:        tex,
		texManager: manager,
		atlas:      *runeAtlas,
		geom:       geometry.NewSquare(),
		mat:        *material.NewTextBlock(tex),
	}
}

func (t *TextBlock) Release() {
	t.atlas.Clear()

	t.texManager.Delete(t.tex)
	t.tex = 0

	t.geom.Delete()
	t.mat.Delete()
}

func (t *TextBlock) Material(r *Renderer, color mgl32.Vec4, ch rune) (runeatlas.RectUV, bool) {
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

func (t *TextBlock) Rune(r *Renderer, model mgl32.Mat4, color mgl32.Vec4, ch rune) {
	runeRect, ok := t.Material(r, color, ch)
	if !ok {
		return
	}

	r.Geometry(t.geom)

	w2h := runeRect.WidthToHeight()
	modelChar := model.Mul4(mgl32.Scale3D(w2h, 1, 1))

	r.Render(&modelChar)
}
