// Copyright (c) 2020-2024 by Marko Gaćeša

package geometry

import (
	"github.com/go-gl/mathgl/mgl32"
	"unsafe"
)

type vertexGeometry struct {
	blockVertexList
	bind
}

func MakeCubeGeometry(makeSide func(side mgl32.Mat3, v *[]blockVertex)) Geometry {
	model := makeCubeModel(makeSide)
	b := bind{}
	b.Load(len(model), BlockVertexSize, unsafe.Pointer(&model[0].v[0]))
	return &vertexGeometry{
		blockVertexList: blockVertexList(model),
		bind:            b,
	}
}

func MakeSphereGeometry(radius float64, wSegs, hSegs int) Geometry {
	model := makeSphereModel(radius, wSegs, hSegs)
	b := bind{}
	b.Load(len(model), BlockVertexSize, unsafe.Pointer(&model[0].v[0]))
	return &vertexGeometry{
		blockVertexList: blockVertexList(model),
		bind:            b,
	}
}
