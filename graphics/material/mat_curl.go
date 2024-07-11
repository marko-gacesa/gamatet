// Copyright (c) 2024 by Marko Gaćeša

package material

var _ Material = (*Curl)(nil)

func NewCurl(tex uint32) *Curl {
	p, err := newProgramBlock(defaultVertexShader, curlFragmentShader, tex)
	if err != nil {
		panic("failed to make curl material: " + err.Error())
	}

	return &Curl{programBlock: *p}
}

type Curl struct {
	programBlock
}

const curlFragmentShader = `
#version 330 core

uniform sampler2D textureSampler;
uniform vec3 lightDirection;
uniform float time;
uniform vec4 objectColor;

in vec2 fragmentTexture;
in vec3 fragmentNormal;

out vec4 outputColor;

void main() {
    vec2 tc = fragmentTexture - vec2(0.5, 0.5);

    float dist = length(tc);
    float radius = 0.7071;
    float percent = (radius - dist) / radius;

    float angle = sin(time * 0.628); 
    float theta = percent * percent * angle * 8.0 + 3.0 * time;

    float s = sin(theta);
    float c = cos(theta);
    tc = vec2(dot(tc, vec2(c, -s)), dot(tc, vec2(s, c)));

    tc = tc + vec2(0.5, 0.5);

	float intensity = texture(textureSampler, tc * 0.5).r;
	vec3 color = vec3(intensity, 0, 0.75 * intensity);
    outputColor = vec4(objectColor.rgb * color, 1.0);
}` + z
