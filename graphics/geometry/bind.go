// Copyright (c) 2020-2024 by Marko Gaćeša

package geometry

import (
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type VertexArray interface {
	Load(vertexCount, vertexSize int, ptr unsafe.Pointer)
	Delete()
	Bind()
}

type bind struct {
	vao uint32
	vbo uint32
}

func (b *bind) Load(vertexCount, vertexSize int, ptr unsafe.Pointer) {
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

func (b *bind) Delete() {
	gl.DeleteVertexArrays(1, &b.vao)
	gl.DeleteBuffers(1, &b.vbo)
}

func (b *bind) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, b.vbo)
	gl.BindVertexArray(b.vao)
}
