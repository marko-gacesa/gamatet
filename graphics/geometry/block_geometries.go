// Copyright (c) 2020-2024 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package geometry

import (
	"unsafe"

	"github.com/go-gl/mathgl/mgl32"
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

func MakeSquareGeometry(makeSquare func(model mgl32.Mat3, v *[]blockVertex)) Geometry {
	model := makeSquareModel(makeSquare)
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
