// Copyright (c) 2020 by Marko Gaćeša

package geometry

import (
	"github.com/go-gl/mathgl/mgl32"
	"unsafe"
)

type vertexGeometry struct {
	vertexList
	bind
}

func MakeCubeGeometry(makeSide func(side mgl32.Mat3, v *[]vertex)) Geometry {
	model := makeCubeModel(makeSide)
	list := vertexList(model)
	b := bind{}
	b.Load(len(model), vertexSize, unsafe.Pointer(&model[0].v[0]))
	return &vertexGeometry{
		vertexList: list,
		bind:       b,
	}
}

func MakeSphereGeometry(radius float64, wSegs, hSegs int) Geometry {
	model := makeSphereModel(radius, wSegs, hSegs)
	list := vertexList(model)
	b := bind{}
	b.Load(len(model), vertexSize, unsafe.Pointer(&model[0].v[0]))
	return &vertexGeometry{
		vertexList: list,
		bind:       b,
	}
}
