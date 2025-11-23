// Copyright (c) 2020-2024 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

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
