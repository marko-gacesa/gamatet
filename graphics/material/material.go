// Copyright (c) 2020 by Marko Gaćeša

package material

import (
	"gamatet/graphics/camera"
	"gamatet/graphics/geometry"
	"github.com/go-gl/mathgl/mgl32"
)

type Material interface {
	Use()
	Delete()

	Refresh()

	Camera(cam *camera.Camera)
	Model(model *mgl32.Mat4)
	Geometry(geo geometry.Geometry)

	Render()
}

func newSimple(vertexShader, fragmentShader string) Material {
	prog, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic("failed to make material: " + err.Error())
	}

	s := simple{program: prog}
	return s
}

type simple struct {
	*program
}

var _ Material = simple{program: nil}
