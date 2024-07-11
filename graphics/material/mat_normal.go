// Copyright (c) 2020-2024 by Marko Gaćeša

package material

// Normal is a material that maps the normal vectors to RGB colors.
func Normal() Material {
	const vertexShader = `
		#version 330

		uniform mat4 viewMatrix;
		uniform mat4 modelMatrix;
		uniform mat3 normalMatrix;

		in vec3 geometryPosition;
		in vec3 geometryNormal;

		out vec3 fragmentNormal;

		void main() {
			fragmentNormal = normalize(normalMatrix * geometryNormal);
			gl_Position = viewMatrix * modelMatrix * vec4(geometryPosition, 1);
		}` + z

	const fragmentShader = `
		#version 330

		in vec3 fragmentNormal;
		out vec4 outputColor;

		void main() {
			outputColor = vec4(fragmentNormal, 1);
		}` + z

	return newSimple(vertexShader, fragmentShader)
}
