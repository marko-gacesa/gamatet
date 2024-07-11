// Copyright (c) 2024 by Marko Gaćeša

package material

import (
	"gamatet/graphics/camera"
	"gamatet/graphics/geometry"
	"github.com/go-gl/mathgl/mgl32"
)

var _ Material = (*programSimple)(nil)

func newSimple(vertexShader, fragmentShader string) *programSimple {
	p, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic("failed to make simple material: " + err.Error())
	}

	return &programSimple{
		program:     *p,
		uniView:     p.uniformLocation("viewMatrix"),
		uniModel:    p.uniformLocation("modelMatrix"),
		uniNorm:     p.uniformLocation("normalMatrix"),
		attribVert:  p.attribLocation("geometryPosition"),
		attribNorm:  p.attribLocation("geometryNormal"),
		attribTexUV: p.attribLocation("geometryTexture"),
	}
}

type programSimple struct {
	program
	uniView     int32
	uniModel    int32
	uniNorm     int32
	attribVert  uint32
	attribNorm  uint32
	attribTexUV uint32
}

func (p *programSimple) Camera(cam *camera.Camera) { uniformMat4(p.uniView, cam.GetView()) }
func (p *programSimple) Model(model *mgl32.Mat4)   { uniformModel(p.uniModel, p.uniNorm, model) }

func (p *programSimple) Geometry(g geometry.Geometry) {
	g.Bind()
	geometry.BindBlockVertex(p.attribVert, p.attribNorm, p.attribTexUV)
	p.meshPrimitiveCount = int32(g.VertexCount())
	p.meshPrimitiveType = g.PrimitiveType()
}
