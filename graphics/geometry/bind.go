// Copyright (c) 2020 by Marko Gaćeša

package geometry

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"unsafe"
)

type GLBinder interface {
	GLLoad(vertexCount, vertexSize int, ptr unsafe.Pointer)
	GLBind()
}

type bind struct {
	vao uint32
	vbo uint32
}

func (b *bind) GLLoad(vertexCount, vertexSize int, ptr unsafe.Pointer) {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, vertexCount*vertexSize, ptr, gl.STATIC_DRAW)

	b.vao = vao
	b.vbo = vbo
}

func (b *bind) GLBind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, b.vbo)
	gl.BindVertexArray(b.vao)
}
