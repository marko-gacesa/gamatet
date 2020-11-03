// Copyright (c) 2020 by Marko Gaćeša

package material

// TexUV is a material that maps the texture UV vectors to RGB colors.
func TexUV() Material {
	const vertexShader = `
		#version 330
		
		uniform mat4 viewMatrix;
		uniform mat4 modelMatrix;
		
		in vec3 geometryPosition;
		in vec3 geometryNormal;
		in vec2 geometryTexture;
		
		out vec2 fragmentTexture;
		
		void main() {
			fragmentTexture = geometryTexture;
			gl_Position = viewMatrix * modelMatrix * vec4(geometryPosition, 1);
		}` + z

	const fragmentShader = `
		#version 330
		
		in vec2 fragmentTexture;	
		out vec4 outputColor;
		
		void main() {
			outputColor = vec4(fragmentTexture, 1, 1);
		}` + z

	return newSimple(vertexShader, fragmentShader)
}
