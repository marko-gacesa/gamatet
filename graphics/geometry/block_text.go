// Copyright (c) 2024 by Marko Gaćeša

package geometry

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"unsafe"
)

func NewText(uv [4]float32) Text {
	return NewTextWithHeightAndScale(uv, 1, 1)
}

func NewTextWithHeight(height float32, uv [4]float32) Text {
	return NewTextWithHeightAndScale(uv, height, 1)
}

func NewTextWithHeightAndScale(uv [4]float32, height, scale float32) Text {
	uvw := uv[2] - uv[0]
	uvh := uv[3] - uv[1]
	if height <= 0 {
		height = 1
	}
	if scale <= 0 {
		scale = 1
	}
	width := scale * height * uvw / uvh

	norm := mgl32.Vec3{0, 0, 1}
	v0 := blockVertex{
		v:  mgl32.Vec3{0, 0, 0},
		n:  norm,
		uv: mgl32.Vec2{uv[0], uv[3]},
	}
	v1 := blockVertex{
		v:  mgl32.Vec3{width, 0, 0},
		n:  norm,
		uv: mgl32.Vec2{uv[2], uv[3]},
	}
	v2 := blockVertex{
		v:  mgl32.Vec3{width, height, 0},
		n:  norm,
		uv: mgl32.Vec2{uv[2], uv[1]},
	}
	v3 := blockVertex{
		v:  mgl32.Vec3{0, height, 0},
		n:  norm,
		uv: mgl32.Vec2{uv[0], uv[1]},
	}

	g := &text{
		vertices: [4]blockVertex{v3, v2, v1, v0},
		width:    width,
		height:   height,
		bind:     bind{},
	}

	g.bind.Load(6, BlockVertexSize, unsafe.Pointer(&g.vertices[0].v[0]))

	return g
}

type text struct {
	vertices [4]blockVertex
	width    float32
	height   float32
	bind
}

func (v *text) VertexSize() int       { return BlockVertexSize }
func (v *text) VertexCount() int      { return 4 }
func (v *text) PrimitiveType() uint32 { return gl.TRIANGLE_FAN }

func (v *text) Width() float32  { return v.width }
func (v *text) Height() float32 { return v.height }
