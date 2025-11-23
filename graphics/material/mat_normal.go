// Copyright (c) 2020-2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

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

	m, err := newProgramSimple(vertexShader, fragmentShader)
	if err != nil {
		panic("failed to make normal material: " + err.Error())
	}

	return m
}
