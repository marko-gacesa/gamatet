// Copyright (c) 2023-2024 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package render

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/marko-gacesa/gamatet/graphics/camera"
	"github.com/marko-gacesa/gamatet/graphics/geometry"
	"github.com/marko-gacesa/gamatet/graphics/material"
)

// Renderer is an object that holds the current render state:
// The camera position, the selected material and the selected geometry.
type Renderer struct {
	cam camera.Camera
	mat material.Material
	geo geometry.Geometry
}

func NewRenderer() *Renderer {
	r := &Renderer{}
	r.cam = camera.Default()
	return r
}

func (r *Renderer) PerspectiveFull(displayW, displayH, contentW, contentH, contentZ int) {
	r.cam.PerspectiveFull(displayW, displayH, contentW, contentH, contentZ)
	if r.mat != nil {
		r.mat.Camera(&r.cam)
	}
}

func (r *Renderer) CameraLookAt(eye, center, up mgl32.Vec3) {
	r.cam.LookAt(eye, center, up)
	if r.mat != nil {
		r.mat.Camera(&r.cam)
	}
}

func (r *Renderer) CameraPerspective(fovy, aspect, near, far float32) {
	r.cam.Perspective(fovy, aspect, near, far)
	if r.mat != nil {
		r.mat.Camera(&r.cam)
	}
}

func (r *Renderer) OrthogonalFull(displayW, displayH, contentW, contentH, contentZ int) {
	r.cam.OrthogonalFull(displayW, displayH, contentW, contentH, contentZ)
	if r.mat != nil {
		r.mat.Camera(&r.cam)
	}
}

func (r *Renderer) Material(mat material.Material) {
	if r.mat != nil {
		r.mat.Reset()
	}

	mat.Use()
	mat.Camera(&r.cam)
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
