// Copyright (c) 2020 by Marko Gaćeša

package gtypes

import "github.com/go-gl/mathgl/mgl32"

type ModelColor struct {
	Model mgl32.Mat4
	Color mgl32.Vec4
}

type ModelColorValue struct {
	Model mgl32.Mat4
	Color mgl32.Vec4
	Value int
}

type PointLight struct {
	Position  mgl32.Vec3
	Color     mgl32.Vec3
	Intensity float32
}
