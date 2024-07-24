// Copyright (c) 2024 by Marko Gaćeša

package material

import (
	"gamatet/graphics/textcanvas"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var _ Material = (*Text)(nil)

func NewText(tex uint32) *Text {
	p, err := newProgramBlock(textVertexShader, defaultFragmentShader, tex)
	if err != nil {
		panic("failed to make text material: " + err.Error())
	}

	return &Text{
		programBlock: *p,
		uniOffsetUV:  p.uniformLocation("uvOffset"),
		uniScaleUV:   p.uniformLocation("uvScale"),
	}
}

// Text is a material that is used for drawing text.
type Text struct {
	programBlock
	uniOffsetUV int32
	uniScaleUV  int32
}

func (t *Text) Use() {
	t.programBlock.Use()
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	uniformVec2(t.uniOffsetUV, mgl32.Vec2{0, 0})
	uniformVec2(t.uniScaleUV, mgl32.Vec2{1, 1})
}

func (t *Text) Reset() {
	t.programBlock.Reset()
	gl.Disable(gl.BLEND)
}

func (t *Text) TexUV(uv textcanvas.RectUV) {
	uniformVec2(t.uniOffsetUV, uv.OffsetUV())
	uniformVec2(t.uniScaleUV, uv.ScaleUV())
}

const textVertexShader = `
#version 330

uniform mat4 viewMatrix;
uniform mat4 modelMatrix;
uniform mat3 normalMatrix;
uniform vec2 uvOffset;
uniform vec2 uvScale;

in vec3 geometryPosition;
in vec3 geometryNormal;
in vec2 geometryTexture;

out vec2 fragmentTexture;
out vec3 fragmentNormal;
out vec4 fragmentPosition;

void main() {
	fragmentTexture = uvOffset + uvScale * geometryTexture;
	fragmentNormal = normalize(normalMatrix * geometryNormal);
	fragmentPosition = modelMatrix * vec4(geometryPosition, 1);
	gl_Position = viewMatrix * fragmentPosition;
}` + z
