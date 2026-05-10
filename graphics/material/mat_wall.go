// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package material

import "github.com/go-gl/mathgl/mgl32"

var _ Material = (*Wall)(nil)

func NewWall(tex uint32) *Wall {
	p, err := newProgramBlock(defaultVertexShader, wallFragmentShader, tex)
	if err != nil {
		panic("failed to make wall material: " + err.Error())
	}

	const (
		uniPlaneDim    = "planeDim"
		uniPointLights = "pointLights[0].position"
	)

	return &Wall{
		programBlock:   *p,
		uniPointLights: p.uniformLocation(uniPointLights),
		uniPlaneDim:    p.uniformLocation(uniPlaneDim),
	}
}

type Wall struct {
	programBlock
	uniPointLights int32
	uniPlaneDim    int32
}

func (p *Wall) Use() {
	p.programBlock.Use()
}

func (p *Wall) Dim(d mgl32.Vec2) {
	uniformVec2(p.uniPlaneDim, d)
}

func (p *Wall) Lights(lights []PointLight) {
	setLights(p.uniPointLights, lights)
}

const wallFragmentShader = `
#version 330

struct Light {
	vec3 position;
	vec3 color;
	float intensity;
};

uniform sampler2D textureSampler;
uniform vec4 objectColor;
uniform vec3 lightDirection;
uniform Light pointLights[16];
uniform vec2 planeDim;

in vec2 fragmentTexture;
in vec3 fragmentNormal;
in vec4 fragmentPosition;

out vec4 outputColor;

void main() {
	vec3 scatteredLight = max(vec3(0.0), dot(fragmentNormal, lightDirection));

	for (int i = 0; i < 16; i++) {
		vec3 lightPosition = pointLights[i].position;
		vec3 lightColor = pointLights[i].color;
		float lightIntensity = pointLights[i].intensity;
		if (lightIntensity <= 0.0) break;

		float d = length(pointLights[i].position - fragmentPosition.xyz);
		if (d > 3.0) continue;

		vec3 lightDir = normalize(lightPosition - fragmentPosition.xyz);

		float attenuation = 1.0 / (1.0 + 0.75*d*d);
		float diffuse = max(0.0, dot(fragmentNormal, lightDir));

		scatteredLight += lightIntensity * lightColor * attenuation * diffuse;
	}

	vec3 ambientColor = vec3(0.1);

	vec3 rgb = min(objectColor.rgb * (ambientColor + scatteredLight), vec3(1.0));
	float gray = texture(textureSampler, mod(fragmentTexture * planeDim, 1.0)).r;
	gray = mix(0.6, 1.0, gray);
	outputColor = vec4(gray * rgb, 1);
}` + z
