// Copyright (c) 2023,2024 by Marko Gaćeša

package camera

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

func (c *Camera) SetDistance(displayW, displayH, contentW, contentH, contentZ int) {
	const fovY = 45 * math.Pi / 180

	aspectRatio := float32(displayW) / float32(displayH)

	scaleX := float32(displayW) / float32(contentW)
	scaleY := float32(displayH) / float32(contentH)

	var cameraDistance float32
	if scaleX > scaleY {
		cameraDistance = float32(float64(contentH) / (2 * math.Tan(fovY/2)))
	} else {
		cameraDistance = float32(float64(contentW) / (2 * (float64(displayW) / float64(displayH)) * math.Tan(fovY/2)))
	}

	z := float32(contentZ)

	c.LookAt(mgl32.Vec3{0, 0, cameraDistance}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	c.Perspective(fovY, aspectRatio, cameraDistance-z, cameraDistance+z)

	//fmt.Printf("camera distance: %6.2f\n", cameraDistance)
}
