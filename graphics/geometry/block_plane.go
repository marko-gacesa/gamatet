// Copyright (c) 2024 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package geometry

import (
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// NewSquare creates a square geometry from (-0.5, -0.5) to (0.5, 0.5), z=0.
func NewSquare() Geometry {
	norm := mgl32.Vec3{0, 0, 1}
	v0 := blockVertex{v: mgl32.Vec3{-0.5, -0.5, 0}, n: norm, uv: mgl32.Vec2{0, 1}}
	v1 := blockVertex{v: mgl32.Vec3{0.5, -0.5, 0}, n: norm, uv: mgl32.Vec2{1, 1}}
	v2 := blockVertex{v: mgl32.Vec3{0.5, 0.5, 0}, n: norm, uv: mgl32.Vec2{1, 0}}
	v3 := blockVertex{v: mgl32.Vec3{-0.5, 0.5, 0}, n: norm, uv: mgl32.Vec2{0, 0}}

	g := &square{
		vertices: [4]blockVertex{v3, v2, v1, v0},
		bind:     bind{},
	}

	g.bind.Load(4, BlockVertexSize, unsafe.Pointer(&g.vertices[0].v[0]))

	return g
}

// NewSquare0 creates a square geometry from (0, 0) to (1, 1), z=0.
func NewSquare0() Geometry {
	norm := mgl32.Vec3{0, 0, 1}
	v0 := blockVertex{v: mgl32.Vec3{0, 0, 0}, n: norm, uv: mgl32.Vec2{0, 1}}
	v1 := blockVertex{v: mgl32.Vec3{1, 0, 0}, n: norm, uv: mgl32.Vec2{1, 1}}
	v2 := blockVertex{v: mgl32.Vec3{1, 1, 0}, n: norm, uv: mgl32.Vec2{1, 0}}
	v3 := blockVertex{v: mgl32.Vec3{0, 1, 0}, n: norm, uv: mgl32.Vec2{0, 0}}

	g := &square{
		vertices: [4]blockVertex{v3, v2, v1, v0},
		bind:     bind{},
	}

	g.bind.Load(4, BlockVertexSize, unsafe.Pointer(&g.vertices[0].v[0]))

	return g
}

type square struct {
	vertices [4]blockVertex
	bind
}

func (v *square) VertexSize() int       { return BlockVertexSize }
func (v *square) VertexCount() int      { return 4 }
func (v *square) PrimitiveType() uint32 { return gl.TRIANGLE_FAN }
