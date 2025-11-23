// Copyright (c) 2020-2024 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

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
