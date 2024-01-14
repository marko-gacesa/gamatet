// Copyright (c) 2020-2024 by Marko Gaćeša

package material

import (
	"gamatet/graphics/gtypes"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

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

func NewRock(tex uint32) Rock {
	prog, err := newProgram(defaultVertexShader, rockFragmentShader)
	if err != nil {
		panic("failed to make rock material: " + err.Error())
	}

	prog.registerUniform(uniLightDirection)
	prog.registerUniform(uniPointLights)
	prog.registerUniform(uniTextureChain)
	prog.registerUniform(uniShouldDrawChain)

	return Rock{
		program: prog,
		tex:     tex,
	}
}

type Rock struct {
	*program
	tex uint32
}

func (p Rock) Refresh() {
	p.program.Refresh()
	p.program.uniformVec3(uniLightDirection, mgl32.Vec3{1, 4, 3}.Normalize())
	p.program.Texture(p.tex)
	p.program.uniformTexture(uniTextureChain, p.tex)
	p.program.uniform1i(uniShouldDrawChain, 0)
}

func (p Rock) ChainTexture(tex uint32) {
	p.program.uniformTexture(uniTextureChain, tex)
	p.program.uniform1i(uniShouldDrawChain, 1)
}

func (p Rock) ClearChain() {
	p.program.uniform1i(uniShouldDrawChain, 0)
}

func (p Rock) Lights(lights []gtypes.PointLight) {
	lightsUniform := p.uniforms[uniPointLights]

	n := int32(len(lights))
	if n > MaxLights {
		n = MaxLights
	}

	for i := int32(0); i < n; i++ {
		gl.Uniform3fv(lightsUniform+i*3, 1, &lights[i].Position[0])
		gl.Uniform3fv(lightsUniform+i*3+1, 1, &lights[i].Color[0])
		gl.Uniform1f(lightsUniform+i*3+2, lights[i].Intensity)
	}
	for i := n; i < MaxLights; i++ {
		gl.Uniform1f(lightsUniform+i*3+2, 0)
	}
}

var _ Material = Rock{program: nil}
