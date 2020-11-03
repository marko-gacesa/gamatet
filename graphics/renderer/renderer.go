// Copyright (c) 2020 by Marko Gaćeša

package renderer

import (
	"gamatet/graphics/camera"
	"gamatet/graphics/geometry"
	"gamatet/graphics/material"
	"github.com/go-gl/mathgl/mgl32"
)

type Renderer struct {
	mat material.Material
	geo geometry.Geometry
	cam *camera.Camera
}

func (r *Renderer) Camera(cam *camera.Camera) {
	r.cam = cam
	if r.mat != nil {
		r.mat.Camera(cam)
	}
}

func (r *Renderer) Material(mat material.Material) {
	if r.mat == mat {
		mat.Refresh()
		return
	}

	mat.Use()
	mat.Refresh()
	mat.Camera(r.cam)
	if r.geo != nil {
		mat.Geometry(r.geo)
	}

	r.mat = mat
}

func (r *Renderer) Geometry(geo geometry.Geometry) {
	if r.geo == geo {
		return
	}

	r.geo = geo

	if r.mat != nil {
		r.mat.Geometry(geo)
	}
}

func (r *Renderer) Render(model *mgl32.Mat4) {
	r.mat.Model(model)
	r.mat.Render()
}
