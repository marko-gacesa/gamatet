// Copyright (c) 2020 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package gutil

import (
	"github.com/go-gl/mathgl/mgl32"
)

type CatmullRom struct {
	seg mgl32.Vec4
}

var catmullRomSpline mgl32.Mat4

func init() {
	const s = 0.5
	catmullRomSpline.SetRow(0, mgl32.Vec4{-s, 2 - s, s - 2, s})
	catmullRomSpline.SetRow(1, mgl32.Vec4{2 * s, s - 3, 3 - 2*s, -s})
	catmullRomSpline.SetRow(2, mgl32.Vec4{-s, 0, s, 0})
	catmullRomSpline.SetRow(3, mgl32.Vec4{0, 1, 0, 0})
}

func (i *CatmullRom) Update(cPrev, cZero, cOne, cNext float32) {
	i.seg = catmullRomSpline.Mul4x1(mgl32.Vec4{cPrev, cZero, cOne, cNext})
}

func (i *CatmullRom) Value(x float32) float32 {
	return mgl32.Vec4{x * x * x, x * x, x, 1}.Dot(i.seg)
}
