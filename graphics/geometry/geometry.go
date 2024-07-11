// Copyright (c) 2020-2024 by Marko Gaćeša

package geometry

type Geometry interface {
	VertexSize() int
	VertexCount() int
	PrimitiveType() uint32
	VertexArray
}

type Text interface {
	Geometry
	Width() float32
	Height() float32
}
