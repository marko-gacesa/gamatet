// Copyright (c) 2020-2024 by Marko Gaćeša

package render

import (
	"gamatet/logic/anim"
	"github.com/go-gl/mathgl/mgl32"
)

func animListUpdate(r *anim.Result) (matrix mgl32.Mat4, color mgl32.Vec4) {
	if r.Feature&anim.Translate > 0 {
		matrix = mgl32.Translate3D(r.DX, r.DY, r.DZ)
	} else {
		matrix = mgl32.Ident4()
	}

	if r.Feature&anim.Rotate > 0 {
		if r.RX != 0.0 {
			matrix = matrix.Mul4(mgl32.HomogRotate3DX(r.RX))
		}
		if r.RY != 0.0 {
			matrix = matrix.Mul4(mgl32.HomogRotate3DY(r.RY))
		}
		if r.RZ != 0.0 {
			matrix = matrix.Mul4(mgl32.HomogRotate3DZ(r.RZ))
		}
	}

	if r.Feature&anim.Scale > 0 {
		matrix = matrix.Mul4(mgl32.Scale3D(r.SX, r.SY, r.SZ))
	}

	if r.Feature&anim.Color > 0 {
		color[0] = r.R
		color[1] = r.G
		color[2] = r.B
		color[3] = r.A
	} else {
		color[0] = 1.0
		color[1] = 1.0
		color[2] = 1.0
		color[3] = 1.0
	}

	return
}
