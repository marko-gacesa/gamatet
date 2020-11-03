// Copyright (c) 2020 by Marko Gaćeša

package material

const defaultVertexShader = `
#version 330

uniform mat4 viewMatrix;
uniform mat4 modelMatrix;
uniform mat3 normalMatrix;

in vec3 geometryPosition;
in vec3 geometryNormal;
in vec2 geometryTexture;

out vec2 fragmentTexture;
out vec3 fragmentNormal;
out vec4 fragmentPosition;

void main() {
	fragmentTexture = geometryTexture;
	fragmentNormal = normalize(normalMatrix * geometryNormal);
	fragmentPosition = modelMatrix * vec4(geometryPosition, 1);
	gl_Position = viewMatrix * fragmentPosition;
}` + z
