// Copyright (c) 2020-2024 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package material

var _ Material = (*Lava)(nil)

func NewLava(tex uint32) *Lava {
	p, err := newProgramBlock(defaultVertexShader, lavaFragmentShader, tex)
	if err != nil {
		panic("failed to make lava material: " + err.Error())
	}

	return &Lava{programBlock: *p}
}

type Lava struct {
	programBlock
}

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
