// Copyright (c) 2020-2025 by Marko Gaćeša

package material

var _ Material = (*Rock)(nil)

const (
	MaxLights = 16
)

func NewRock(tex uint32) *Rock {
	p, err := newProgramBlock(defaultVertexShader, rockFragmentShader, tex)
	if err != nil {
		panic("failed to make rock material: " + err.Error())
	}

	const (
		uniTextureChain    = "textureSamplerChain"
		uniShouldDrawChain = "shouldDrawChain"
		uniPointLights     = "pointLights[0].position"
	)

	return &Rock{
		programBlock:       *p,
		uniPointLights:     p.uniformLocation(uniPointLights),
		uniTexChain:        p.uniformLocation(uniTextureChain),
		uniShouldDrawChain: p.uniformLocation(uniShouldDrawChain),
	}
}

type Rock struct {
	programBlock
	uniPointLights     int32
	uniTexChain        int32
	uniShouldDrawChain int32
}

func (p *Rock) Use() {
	p.programBlock.Use()
	uniform1i(p.uniShouldDrawChain, 0)
}

func (p *Rock) ChainTexture(tex uint32) {
	uniformTexture(p.uniTexChain, tex)
	uniform1i(p.uniShouldDrawChain, 1)
}

func (p *Rock) ClearChain() {
	uniform1i(p.uniShouldDrawChain, 0)
}

func (p *Rock) Lights(lights []PointLight) {
	n := min(int32(len(lights)), MaxLights)

	for i := range n {
		uniformVec3(p.uniPointLights+i*3, lights[i].Position)
		uniformVec3(p.uniPointLights+i*3+1, lights[i].Color)
		uniform1f(p.uniPointLights+i*3+2, lights[i].Intensity)
	}
	for i := n; i < MaxLights; i++ {
		uniform1f(p.uniPointLights+i*3+2, 0)
	}
}

const rockFragmentShader = `
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
	float gray = texture(textureSampler, fragmentTexture).r;
	gray = mix(0.6, 1.0, gray);
	outputColor = vec4(gray * rgb, 1);
}` + z
