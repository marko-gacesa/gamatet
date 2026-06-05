// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package scene

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/marko-gacesa/gamatet/graphics/geometry"
	"github.com/marko-gacesa/gamatet/graphics/material"
	"github.com/marko-gacesa/gamatet/graphics/render"
	"github.com/marko-gacesa/gamatet/graphics/scene/base"
	"github.com/marko-gacesa/gamatet/graphics/texture"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

type Background struct {
	base.Base
	matBack         *material.Clouds
	geomSquarePlain geometry.Geometry
}

var _ screen.Screen = (*Background)(nil)

func NewBackground(renderer *render.Renderer, tex *texture.Manager) *Background {
	matBack := material.NewClouds()
	matBack.BackColorTop(mgl32.Vec3{0.1, 0.14, 0.2})
	matBack.BackColorBottom(mgl32.Vec3{0.05, 0.07, 0.1})
	matBack.CloudColor(mgl32.Vec3{0.0, 0.0, 0.0})
	geomSquarePlain := geometry.NewSquare() // -0.5..0.5 square
	return &Background{
		Base:            base.NewBase(renderer, tex),
		matBack:         matBack,
		geomSquarePlain: geomSquarePlain,
	}
}

func (m *Background) Release() {
	m.matBack.Delete()
	m.geomSquarePlain.Delete()
}

func (m *Background) Render() {
	gl.DepthMask(false)
	gl.Disable(gl.DEPTH_TEST)

	r := m.Renderer()
	r.Orthogonal2D(1, 1) // x and y: -0.5..0.5
	r.Material(m.matBack)
	r.Geometry(m.geomSquarePlain)
	r.Render(mgl32.Ident4())

	gl.DepthMask(true)
	gl.Enable(gl.DEPTH_TEST)
}
