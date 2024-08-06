// Copyright (c) 2020-2024 by Marko Gaćeša

package material

import (
	"gamatet/graphics/camera"
	"gamatet/graphics/geometry"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var _ Material = (*programBlock)(nil)

var lightDir = mgl32.Vec3{1, 4, 3}.Normalize()

type programBlock struct {
	program
	uniView     int32
	uniModel    int32
	uniNorm     int32
	uniTime     int32
	uniTex      int32
	uniColor    int32
	uniLightDir int32
	attribVert  uint32
	attribNorm  uint32
	attribTexUV uint32
	tex         uint32
}

func newProgramBlock(vertexShaderSource, fragmentShaderSource string, tex uint32) (*programBlock, error) {
	p, err := newProgram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		return nil, err
	}

	return &programBlock{
		program:     *p,
		uniView:     p.uniformLocation("viewMatrix"),
		uniModel:    p.uniformLocation("modelMatrix"),
		uniNorm:     p.uniformLocation("normalMatrix"),
		uniTime:     p.uniformLocation("time"),
		uniTex:      p.uniformLocation("textureSampler"),
		uniColor:    p.uniformLocation("objectColor"),
		uniLightDir: p.uniformLocation("lightDirection"),
		attribVert:  p.attribLocation("geometryPosition"),
		attribNorm:  p.attribLocation("geometryNormal"),
		attribTexUV: p.attribLocation("geometryTexture"),
		tex:         tex,
	}, nil
}

func (p *programBlock) Camera(cam *camera.Camera) { uniformMat4(p.uniView, cam.GetView()) }
func (p *programBlock) Model(model *mgl32.Mat4)   { uniformModel(p.uniModel, p.uniNorm, model) }

func (p *programBlock) Use() {
	p.program.Use()
	uniform1f(p.uniTime, float32(glfw.GetTime()))
	uniformVec3(p.uniLightDir, lightDir)
	uniformTexture(p.uniTex, p.tex)
	uniformVec4(p.uniColor, mgl32.Vec4{1, 1, 1, 1})
}

func (p *programBlock) Color(color mgl32.Vec4) { uniformVec4(p.uniColor, color) }

func (p *programBlock) Geometry(g geometry.Geometry) {
	g.Bind()
	geometry.BindBlockVertex(p.attribVert, p.attribNorm, p.attribTexUV)
	p.meshPrimitiveCount = int32(g.VertexCount())
	p.meshPrimitiveType = g.PrimitiveType()
}
