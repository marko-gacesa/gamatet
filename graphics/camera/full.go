// Copyright (c) 2023,2024 by Marko Gaćeša

package camera

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// PerspectiveFull sets the camera to the location (0, 0, Z) so that the thin content at XY plane can be fully seen.
func (c *Camera) PerspectiveFull(displayW, displayH, contentW, contentH, contentZ int) {
	const fovY = 45 * math.Pi / 180
	const fovY2Z = 0.8284271247461902 // 2 * math.Tan(fovY/2)

	aspectRatio := float32(displayW) / float32(displayH)

	scaleX := float32(displayW) / float32(contentW)
	scaleY := float32(displayH) / float32(contentH)

	var cameraDistance float32
	if scaleX > scaleY {
		cameraDistance = float32(contentH) / fovY2Z
	} else {
		cameraDistance = float32(contentW) / (fovY2Z * aspectRatio)
	}

	z := float32(contentZ)

	c.LookAt(mgl32.Vec3{0, 0, cameraDistance}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	c.Perspective(fovY, aspectRatio, cameraDistance-z, cameraDistance+z)

	//fmt.Printf("camera distance: %6.2f\n", cameraDistance)
}

// OrthogonalFull sets the orthogonal projection so that the entire content is visible.
func (c *Camera) OrthogonalFull(displayW, displayH, contentW, contentH, contentZ int) {
	scaleX := float32(displayW) / float32(contentW)
	scaleY := float32(displayH) / float32(contentH)

	var w, h float32
	z := float32(contentZ)

	if scaleX > scaleY {
		w = float32(contentW) * (scaleX / scaleY) / 2
		h = float32(contentH) / 2
	} else {
		w = float32(contentW) / 2
		h = float32(contentH) * (scaleY / scaleX) / 2
	}

	cameraDistance := z
	c.LookAt(mgl32.Vec3{0, 0, cameraDistance}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	c.Projection(mgl32.Ortho(-w, w, -h, h, cameraDistance-z, cameraDistance+z))

	//fmt.Printf("orthogonal projection: w=%6.2f h=%6.2f\n", w, h)
}
