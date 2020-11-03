// Copyright (c) 2020 by Marko Gaćeša

package geometry

import (
	"github.com/go-gl/mathgl/mgl32"
	"unsafe"
)

var (
	Frame       Geometry
	RoundedCube Geometry
	Gem         Geometry
	Cube        Geometry
	Die         Geometry
	Octahedron  Geometry
	DentCube    Geometry
	Star6       Geometry
	Star8       Geometry
)

func LoadAll() {
	Cube = makeCubeGeometry(makeSideSimple)
	Frame = makeCubeGeometry(makeSideFrame)
	Die = makeCubeGeometry(makeSideDie)
	RoundedCube = makeCubeGeometry(makeSideRounded)
	Gem = makeCubeGeometry(makeSideTruncated)
	Octahedron = makeCubeGeometry(makeSideStellatedOctahedron)
	DentCube = makeCubeGeometry(makeSideDent)
	Star6 = makeCubeGeometry(makeSideStar6)
	Star8 = makeCubeGeometry(makeSideStar8)
}

type cubeGeometry struct {
	vertexList
	bind
}

func makeCubeGeometry(makeSide func(side mgl32.Mat3, v *[]vertex)) Geometry {
	model := makeCubeModel(makeSide)
	list := vertexList(model)
	glb := bind{}
	glb.GLLoad(len(model), vertexSize, unsafe.Pointer(&model[0].v[0]))
	return &cubeGeometry{
		vertexList: list,
		bind:       glb,
	}
}
