// Copyright (c) 2024 by Marko Gaćeša

package material

import (
	"gamatet/graphics/gtypes"
	"github.com/go-gl/mathgl/mgl32"
)

var _ Material = (*Iron)(nil)

func NewIron(tex uint32) *Iron {
	p, err := newProgramBlock(defaultVertexShader, ironFragmentShader, tex)
	if err != nil {
		panic("failed to make rock material: " + err.Error())
	}

	const (
		uniPointLights = "pointLights[0].position"
	)

	return &Iron{
		programBlock:   *p,
		uniPointLights: p.uniformLocation(uniPointLights),
	}
}

type Iron struct {
	programBlock
	uniPointLights int32
}

func (p *Iron) Use() {
	p.programBlock.Use()
	p.programBlock.Color(mgl32.Vec4{1, 1, 1, 1})
}

func (p *Iron) Lights(lights []gtypes.PointLight) {
	n := int32(len(lights))
	if n > MaxLights {
		n = MaxLights
	}

	for i := int32(0); i < n; i++ {
		uniformVec3(p.uniPointLights+i*3, lights[i].Position)
		uniformVec3(p.uniPointLights+i*3+1, lights[i].Color)
		uniform1f(p.uniPointLights+i*3+2, lights[i].Intensity)
	}
	for i := n; i < MaxLights; i++ {
		uniform1f(p.uniPointLights+i*3+2, 0)
	}
}

const ironFragmentShader = `
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
uniform sampler2D textureSamplerChain;
uniform int shouldDrawChain;

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

	if (shouldDrawChain > 0) {
		float chain = texture(textureSamplerChain, fragmentTexture).r;
		if (chain > 0.0) {
			vec3 rgb = min(chain * (ambientColor + scatteredLight), vec3(1.0));
			outputColor = vec4(rgb, 1);
			return;
		}
	}

	vec3 rgb = min(objectColor.rgb * (ambientColor + scatteredLight), vec3(1.0));
	vec2 uv = fragmentTexture;

	if (uv.x > 0.1 && uv.x < 0.9 && uv.y > 0.1 && uv.y < 0.9) {
		uv.x = mod(uv.x*0.2, 1.0);
		uv.y = mod(uv.y*1.2, 1.0);
		float gray = texture(textureSampler, uv).r;
		gray = mix(0.8, 1.0, gray);
		outputColor = vec4(gray * rgb, objectColor.a);
	} else {
		float gray = texture(textureSampler, uv).r;
		gray = mix(0.5, 1.0, gray);
		outputColor = vec4(gray * rgb, objectColor.a);
	}
}` + z
