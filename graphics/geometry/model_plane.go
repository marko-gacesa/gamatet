// Copyright (c) 2024 by Marko Gaćeša

package geometry

import (
	"github.com/go-gl/mathgl/mgl32"
)

func makeSquareModel(makePlane func(model mgl32.Mat3, v *[]blockVertex)) []blockVertex {
	v := make([]blockVertex, 0, 8)
	makePlane(mgl32.Ident3(), &v)
	return v
}
