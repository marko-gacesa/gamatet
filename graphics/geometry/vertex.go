// Copyright (c) 2020 by Marko Gaćeša

package geometry

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"unsafe"
)

type vertex struct {
	v  mgl32.Vec3
	n  mgl32.Vec3
	uv mgl32.Vec2
}

const vertexSize = 32 // int(unsafe.Sizeof(vertex{}))

func gen(p, n, uv mgl32.Vec3) (q vertex) {
	q.setV(p)
	q.setN(n)
	q.setUV(uv)
	return
}

func (q *vertex) init(p, n, uv mgl32.Vec3) {
	q.setV(p)
	q.setN(n)
	q.setUV(uv)
}

func (q *vertex) setV(p mgl32.Vec3) {
	q.v = p
}

func (q *vertex) setN(n mgl32.Vec3) {
	q.n = n
}

func (q *vertex) setUV(uv mgl32.Vec3) {
	q.uv = uv.Vec2()
}

type vertexList []vertex

func (v vertexList) VertexSize() int          { return vertexSize }
func (v vertexList) VertexCount() int         { return len(v) }
func (v vertexList) DataPtr() unsafe.Pointer  { return unsafe.Pointer(&v[0].v[0]) }
func (v vertexList) DataOffsetVertex() int    { return 0 }
func (v vertexList) DataOffsetNormal() int    { return 3 * 4 } // 3 float32s are before normal data
func (v vertexList) DataOffsetTextureUV() int { return 6 * 4 } // 6 float32s are before normal data
func (v vertexList) PrimitiveType() uint32    { return gl.TRIANGLES }
