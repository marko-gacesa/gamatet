// Copyright (c) 2020 by Marko Gaćeša

package geometry

import "unsafe"

type Geometry interface {
	VertexSize() int
	VertexCount() int

	DataPtr() unsafe.Pointer
	DataOffsetVertex() int
	DataOffsetNormal() int
	DataOffsetTextureUV() int

	PrimitiveType() uint32

	GLBinder
}
