// Copyright (c) 2020-2024 by Marko Gaćeša

package material

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/marko-gacesa/gamatet/graphics/camera"
	"github.com/marko-gacesa/gamatet/graphics/geometry"
)

// Material is an abstraction of shaders.
type Material interface {
	Use()
	Reset()
	Delete()

	Camera(cam *camera.Camera)
	Model(model *mgl32.Mat4)
	Geometry(geo geometry.Geometry)

	Render()
}
