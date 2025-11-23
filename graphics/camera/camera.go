// Copyright (c) 2020 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package camera

import "github.com/go-gl/mathgl/mgl32"

type Camera struct {
	lookAt     mgl32.Mat4
	projection mgl32.Mat4
	view       mgl32.Mat4
}

func Default() Camera {
	lookAt := mgl32.LookAtV(mgl32.Vec3{0, 0, 1}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	projection := mgl32.Perspective(60, 16/9, 1, 10)
	view := projection.Mul4(lookAt)
	return Camera{
		lookAt:     lookAt,
		projection: projection,
		view:       view,
	}
}

func (c *Camera) LookAt(eye, center, up mgl32.Vec3) {
	c.lookAt = mgl32.LookAtV(eye, center, up)
	c.view = c.projection.Mul4(c.lookAt)
}

func (c *Camera) Projection(projection mgl32.Mat4) {
	c.projection = projection
	c.view = c.projection.Mul4(c.lookAt)
}

func (c *Camera) Perspective(fovy, aspect, near, far float32) {
	c.projection = mgl32.Perspective(fovy, aspect, near, far)
	c.view = c.projection.Mul4(c.lookAt)
}

func (c *Camera) GetView() *mgl32.Mat4 {
	return &c.view
}

func (c *Camera) GetLookAt() *mgl32.Mat4 {
	return &c.lookAt
}

func (c *Camera) GetProjection() *mgl32.Mat4 {
	return &c.projection
}
