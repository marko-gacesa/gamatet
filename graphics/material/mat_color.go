// Copyright (c) 2024 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package material

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/marko-gacesa/gamatet/graphics/camera"
	"github.com/marko-gacesa/gamatet/graphics/geometry"
)

var _ Material = (*Color)(nil)

func NewColor() *Color {
	p, err := newProgram(defaultVertexShader, colorFragmentShader)
	if err != nil {
		panic("failed to make color material: " + err.Error())
	}

	return &Color{
		program:     *p,
		uniView:     p.uniformLocation("viewMatrix"),
		uniModel:    p.uniformLocation("modelMatrix"),
		uniNorm:     p.uniformLocation("normalMatrix"),
		uniColor:    p.uniformLocation("objectColor"),
		uniLightDir: p.uniformLocation("lightDirection"),
		attribVert:  p.attribLocation("geometryPosition"),
		attribNorm:  p.attribLocation("geometryNormal"),
		attribTexUV: p.attribLocation("geometryTexture"),
	}
}

type Color struct {
	program
	uniView     int32
	uniModel    int32
	uniNorm     int32
	uniColor    int32
	uniLightDir int32
	attribVert  uint32
	attribNorm  uint32
	attribTexUV uint32
}

func (p *Color) Camera(cam *camera.Camera) { uniformMat4(p.uniView, cam.GetView()) }
func (p *Color) Model(model *mgl32.Mat4)   { uniformModel(p.uniModel, p.uniNorm, model) }

func (p *Color) Use() {
	p.program.Use()
	uniformVec3(p.uniLightDir, lightDir)
}

func (p *Color) Reset() {
	p.program.Reset()
}

func (p *Color) Color(color mgl32.Vec4) { uniformVec4(p.uniColor, color) }

func (p *Color) Geometry(g geometry.Geometry) {
	g.Bind()
	geometry.BindBlockVertex(p.attribVert, p.attribNorm, p.attribTexUV)
	p.meshPrimitiveCount = int32(g.VertexCount())
	p.meshPrimitiveType = g.PrimitiveType()
}

const colorFragmentShader = `
#version 330 core

uniform vec3 lightDirection;
uniform vec4 objectColor;

in vec3 fragmentNormal;

out vec4 outputColor;

void main() {
	vec3 ambientColor = vec3(0.3);
	float diffuse = max(0.0, dot(fragmentNormal, lightDirection));
	vec3 scatteredLight = ambientColor + diffuse;
	vec3 rgb = min(objectColor.rgb * scatteredLight, vec3(1.0));
	outputColor = vec4(rgb, objectColor.a);
}` + z
