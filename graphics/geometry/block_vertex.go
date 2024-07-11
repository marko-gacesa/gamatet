// Copyright (c) 2020-2024 by Marko Gaćeša

package geometry

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type blockVertex struct {
	v  mgl32.Vec3
	n  mgl32.Vec3
	uv mgl32.Vec2
}

const (
	BlockVertexSize = 32 // int(unsafe.Sizeof(blockVertex{}))
)

func gen(p, n, uv mgl32.Vec3) (q blockVertex) {
	q.setV(p)
	q.setN(n)
	q.setUV(uv)
	return
}

func (q *blockVertex) init(p, n, uv mgl32.Vec3) {
	q.setV(p)
	q.setN(n)
	q.setUV(uv)
}

func (q *blockVertex) setV(p mgl32.Vec3)   { q.v = p }
func (q *blockVertex) setN(n mgl32.Vec3)   { q.n = n }
func (q *blockVertex) setUV(uv mgl32.Vec3) { q.uv = uv.Vec2() }

type blockVertexList []blockVertex

func (v blockVertexList) VertexSize() int       { return BlockVertexSize }
func (v blockVertexList) VertexCount() int      { return len(v) }
func (v blockVertexList) PrimitiveType() uint32 { return gl.TRIANGLES }

func BindBlockVertex(attribVert, attribNorm, attribTexUV uint32) {
	const (
		offsetVertex = 0  // 0 float32s are before vertex data - there are three floats here: x, y, z
		offsetNormal = 12 // 3 float32s (3*4=12bytes) are before normal data - there are three floats here
		offsetTexUV  = 24 // 6 float32s (6*4=24bytes) are before tex UV data - there are two floats here: u, v
	)

	gl.EnableVertexAttribArray(attribVert)
	gl.VertexAttribPointerWithOffset(attribVert, 3, gl.FLOAT, false, BlockVertexSize, offsetVertex)

	gl.EnableVertexAttribArray(attribNorm)
	gl.VertexAttribPointerWithOffset(attribNorm, 3, gl.FLOAT, false, BlockVertexSize, offsetNormal)

	gl.EnableVertexAttribArray(attribTexUV)
	gl.VertexAttribPointerWithOffset(attribTexUV, 2, gl.FLOAT, false, BlockVertexSize, offsetTexUV)
}
