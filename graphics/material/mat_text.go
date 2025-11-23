// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package material

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/marko-gacesa/gamatet/graphics/runeatlas"
)

var _ Material = (*Text)(nil)

func NewText(tex uint32) *Text {
	p, err := newProgramSimple(textVertexShader, textFragmentShader)
	if err != nil {
		panic("failed to make text material: " + err.Error())
	}

	return &Text{
		programSimple:  *p,
		uniTime:        p.uniformLocation("time"),
		uniColor:       p.uniformLocation("objectColor"),
		uniTex:         p.uniformLocation("textureSampler"),
		uniOffsetUV:    p.uniformLocation("uvOffset"),
		uniScaleUV:     p.uniformLocation("uvScale"),
		uniInvertColor: p.uniformLocation("invertColor"),
		tex:            tex,
	}
}

// Text is a material that is used for drawing text.
type Text struct {
	programSimple
	uniTime        int32
	uniColor       int32
	uniTex         int32
	uniOffsetUV    int32
	uniScaleUV     int32
	uniInvertColor int32
	tex            uint32
}

func (t *Text) Use() {
	t.programSimple.Use()

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.Disable(gl.DEPTH_TEST)

	uniform1f(t.uniTime, float32(glfw.GetTime()))
	uniformVec4(t.uniColor, mgl32.Vec4{1, 1, 1, 1})
	uniformTexture(t.uniTex, t.tex)
	uniformVec2(t.uniOffsetUV, mgl32.Vec2{0, 0})
	uniformVec2(t.uniScaleUV, mgl32.Vec2{1, 1})
	uniform1i(t.uniInvertColor, 0)
}

func (t *Text) Reset() {
	t.programSimple.Reset()

	gl.Enable(gl.DEPTH_TEST)

	gl.Disable(gl.BLEND)
}

func (t *Text) Color(color mgl32.Vec4) { uniformVec4(t.uniColor, color) }

func (t *Text) Texture(tex uint32) {
	uniformTexture(t.uniTex, tex)
	t.tex = tex
}

func (t *Text) TexUV(uv runeatlas.RectUV) {
	uniformVec2(t.uniOffsetUV, uv.OffsetUV())
	uniformVec2(t.uniScaleUV, uv.ScaleUV())
}

func (t *Text) InvertOn() {
	uniform1i(t.uniInvertColor, 1)
}

func (t *Text) InvertOff() {
	uniform1i(t.uniInvertColor, 0)
}

const textVertexShader = `
#version 330

uniform mat4 viewMatrix;
uniform mat4 modelMatrix;
uniform vec2 uvOffset;
uniform vec2 uvScale;

in vec3 geometryPosition;
in vec2 geometryTexture;

out vec2 fragmentTexture;
out vec4 fragmentPosition;

void main() {
	fragmentTexture = uvOffset + uvScale * geometryTexture;
	fragmentPosition = modelMatrix * vec4(geometryPosition, 1);
	gl_Position = viewMatrix * fragmentPosition;
}` + z

const textFragmentShader = `
#version 330

uniform sampler2D textureSampler;
uniform vec4 objectColor;
uniform int invertColor;
uniform float time;

in vec2 fragmentTexture;

out vec4 outputColor;

void main() {
	float alpha = texture(textureSampler, fragmentTexture).r;

	if (invertColor != 0) {
		alpha = 1.0 - alpha;
	}

	outputColor = objectColor * alpha;
}` + z
