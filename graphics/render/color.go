// Copyright (c) 2020-2024 by Marko Gaćeša

package render

import (
	"github.com/go-gl/mathgl/mgl32"
)

var colorWhite = mgl32.Vec4{1, 1, 1, 1}

func colorVector(c uint32) (color mgl32.Vec4) {
	color[0] = float32(c>>24) / 255
	color[1] = float32((c&0xFF0000)>>16) / 255
	color[2] = float32((c&0xFF00)>>8) / 255
	color[3] = float32(c&0xFF) / 255
	return
}

func mulColor(c1, c2 mgl32.Vec4) (color mgl32.Vec4) {
	color[0] = c1[0] * c2[0]
	color[1] = c1[1] * c2[1]
	color[2] = c1[2] * c2[2]
	color[3] = c1[3] * c2[3]
	return
}
