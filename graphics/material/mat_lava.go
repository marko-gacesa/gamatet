// Copyright (c) 2020 by Marko Gaćeša

package material

import (
	"github.com/go-gl/mathgl/mgl32"
)

const lavaFragmentShader = `
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
    float turbulence = 0.5 * (sin(uv.x * 10.0 + time) + sin(uv.y * 10.0 + time));

    // Color gradient for the lava effect
    vec3 color1 = vec3(0.4, 0.1, 0.0);
    vec3 color2 = vec3(1.0, 0.5, 0.0);

	float intensity = mod(texture(textureSampler, uv).r + turbulence, 1.0);
	intensity = 1.0 - intensity * intensity;
	vec3 finalColor = objectColor.rgb * mix(color1, color2, intensity) + diffuse * 0.08;

    outputColor = vec4(finalColor, 1.0);
}` + z

func NewLava(tex uint32) Lava {
	prog, err := newProgram(defaultVertexShader, lavaFragmentShader)
	if err != nil {
		panic("failed to make lava material: " + err.Error())
	}

	prog.registerUniform(uniLightDirection)

	return Lava{
		program: prog,
		tex:     tex,
	}
}

type Lava struct {
	*program
	tex uint32
}

func (p Lava) Refresh() {
	p.program.Refresh()
	p.program.uniformVec3(uniLightDirection, mgl32.Vec3{1, 4, 3}.Normalize())
	p.program.Texture(p.tex)
}

var _ Material = Lava{program: nil}
