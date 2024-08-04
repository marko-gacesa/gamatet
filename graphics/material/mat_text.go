// Copyright (c) 2024 by Marko Gaćeša

package material

import (
	"gamatet/graphics/runeatlas"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var _ Material = (*Text)(nil)

func NewText(tex uint32) *Text {
	p, err := newProgramBlock(textVertexShader, textFragmentShader, tex)
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
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	uniformVec2(t.uniOffsetUV, mgl32.Vec2{0, 0})
	uniformVec2(t.uniScaleUV, mgl32.Vec2{1, 1})
}

func (t *Text) Reset() {
	t.programBlock.Reset()
	gl.Enable(gl.DEPTH_TEST)
	gl.Disable(gl.BLEND)
}

func (t *Text) Texture(tex uint32) {
	uniformTexture(t.uniTex, tex)
	t.tex = tex
}

func (t *Text) TexUV(uv runeatlas.RectUV) {
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

const textFragmentShader = `
#version 330

uniform sampler2D textureSampler;
uniform vec4 objectColor;
uniform vec3 lightDirection;
uniform float time;

in vec2 fragmentTexture;
in vec3 fragmentNormal;

out vec4 outputColor;

void main() {
	float alpha = texture(textureSampler, fragmentTexture).r;
	if (alpha < 0.001) {
		outputColor = vec4(0);
		return;
	}

	vec3 ambientColor = vec3(0.3);
	vec3 lightDirectionColor = vec3(1.0);

	float diffuse = max(0.0, dot(fragmentNormal, lightDirection));
	vec3 scatteredLight = ambientColor + lightDirectionColor * diffuse;
	vec3 rgb = min(objectColor.rgb * scatteredLight, vec3(1.0));

	outputColor = vec4(rgb, objectColor.a * alpha);
}` + z
