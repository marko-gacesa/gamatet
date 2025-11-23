// Copyright (c) 2024, 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package material

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/marko-gacesa/gamatet/graphics/camera"
	"github.com/marko-gacesa/gamatet/graphics/geometry"
)

var _ Material = (*programSimple)(nil)

func newProgramSimple(vertexShader, fragmentShader string) (*programSimple, error) {
	p, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		return nil, err
	}

	return &programSimple{
		program:     *p,
		uniView:     p.uniformLocation("viewMatrix"),
		uniModel:    p.uniformLocation("modelMatrix"),
		uniNorm:     p.uniformLocation("normalMatrix"),
		attribVert:  p.attribLocation("geometryPosition"),
		attribNorm:  p.attribLocation("geometryNormal"),
		attribTexUV: p.attribLocation("geometryTexture"),
	}, nil
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
