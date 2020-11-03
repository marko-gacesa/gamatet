// Copyright (c) 2020 by Marko Gaćeša

package material

import (
	"github.com/go-gl/mathgl/mgl32"
)

const acidFragmentShader = `
#version 330 core

uniform sampler2D textureSampler;
uniform vec3 lightDirection;
uniform float time;
uniform vec4 objectColor;

in vec2 fragmentTexture;
in vec3 fragmentNormal;

out vec4 outputColor;

void main() {
	float diffuse = max(0.0, dot(fragmentNormal, lightDirection));

	vec2 uv = fragmentTexture;

    // Color gradient for the acid effect
    vec3 color1 = vec3(0.0, 1.0, 0.0);
    vec3 color2 = vec3(0.0, 0.0, 0.0);

	float intensity = mod(texture(textureSampler, uv).r + time, 1.0);
	intensity *= intensity;
	vec3 finalColor = objectColor.rbg * mix(color1, color2, intensity) + diffuse * 0.2;

    outputColor = vec4(finalColor, 1.0);
}` + z

func NewAcid(tex uint32) Acid {
	prog, err := newProgram(defaultVertexShader, acidFragmentShader)
	if err != nil {
		panic("failed to make acid material: " + err.Error())
	}

	prog.registerUniform(uniLightDirection)

	return Acid{
		program: prog,
		tex:     tex,
	}
}

type Acid struct {
	*program
	tex uint32
}

func (p Acid) Refresh() {
	p.program.Refresh()
	p.program.uniformVec3(uniLightDirection, mgl32.Vec3{1, 4, 3}.Normalize())
	p.program.Texture(p.tex)
}

var _ Material = Acid{program: nil}
