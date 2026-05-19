// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package material

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var _ Material = (*Clouds)(nil)

func NewClouds() *Clouds {
	p, err := newProgramSimple(defaultVertexShader, skyFragmentShader)
	if err != nil {
		panic("failed to make clouds material: " + err.Error())
	}

	return &Clouds{
		programSimple: *p,
		uniTime:       p.uniformLocation("time"),
		uniResolution: p.uniformLocation("resolution"),
		uniBackColor0: p.uniformLocation("backColorGrad0"),
		uniBackColor1: p.uniformLocation("backColorGrad1"),
		uniCloudColor: p.uniformLocation("cloudColor"),
		backColor0:    mgl32.Vec3{0.3, 0.5, 0.8},
		backColor1:    mgl32.Vec3{0.7, 0.85, 1.0},
		cloudColor:    mgl32.Vec3{1, 1, 1},
	}
}

type Clouds struct {
	programSimple
	uniTime       int32
	uniResolution int32
	uniBackColor0 int32
	uniBackColor1 int32
	uniCloudColor int32

	backColor0 mgl32.Vec3
	backColor1 mgl32.Vec3
	cloudColor mgl32.Vec3
}

func (p *Clouds) BackColorTop(c mgl32.Vec3)    { p.backColor1 = c }
func (p *Clouds) BackColorBottom(c mgl32.Vec3) { p.backColor0 = c }
func (p *Clouds) CloudColor(c mgl32.Vec3)      { p.cloudColor = c }

func (p *Clouds) Use() {
	p.program.Use()

	var viewport [4]int32
	gl.GetIntegerv(gl.VIEWPORT, &viewport[0])

	width := viewport[2]
	height := viewport[3]

	uniform1f(p.uniTime, float32(glfw.GetTime()))
	uniformVec2(p.uniResolution, mgl32.Vec2{float32(width), float32(height)})

	uniformVec3(p.uniBackColor0, p.backColor0)
	uniformVec3(p.uniBackColor1, p.backColor1)
	uniformVec3(p.uniCloudColor, p.cloudColor)
}

const skyFragmentShader = `
#version 330

out vec4 outputColor;

uniform vec2 resolution;
uniform float time;

uniform vec3 backColorGrad0;
uniform vec3 backColorGrad1;
uniform vec3 cloudColor;

float hash(vec2 p) {
	p = fract(p * vec2(573.84, 731.36));
	p += dot(p, p + 27.61);
	return fract(p.x * p.y);
}

float noise(vec2 p) {
	vec2 i = floor(p);
	vec2 f = fract(p);

	float a = hash(i);
	float b = hash(i + vec2(1.0, 0.0));
	float c = hash(i + vec2(0.0, 1.0));
	float d = hash(i + vec2(1.0, 1.0));

	vec2 u = f * f * (3.0 - 2.0 * f);
	return mix(a, b, u.x) + (c - a) * u.y * (1.0 - u.x) + (d - b) * u.x * u.y;
}

// fbm stands for Fractal Brownian Motion.
float fbm(vec2 frequency, int octaves) {
	float value = 0.0;
	float amplitude = 0.5;
	for (int i = 0; i < octaves; ++i) {
		value += amplitude * noise(frequency);
		frequency *= 2.0;
		amplitude *= 0.5;
	}
	return value;
}

void main() {
	vec2 uv = gl_FragCoord.xy / resolution.xy;
	uv.x *= resolution.x / resolution.y; // keep aspect ratio correct

	// slow cloud movement
	vec2 motion = vec2(time * 0.3, time * 0.04); // hardcoded cloud movement vector (0.3, 0.04)

	// generate cloud field, several layers
	float n = 0.0;
	n += 1.0 * fbm(uv * 8.0 + motion, 4); // small blurry patches
	n += 0.5 * fbm(uv * 2.0 + motion, 6); // large sharper chunks
	n += 0.5 * fbm(uv * 5.0 + motion, 8); // medium detail clouds
	n /= 1.7; // normalize
	float clouds = smoothstep(0.45, 0.75, n);

	vec3 backColor = mix(backColorGrad0, backColorGrad1, uv.y);
	vec3 color = mix(backColor, cloudColor, clouds);

	outputColor = vec4(color, 1.0);
}` + z
