// Copyright (c) 2020 by Marko Gaćeša

package material

const defaultFragmentShader = `
#version 330

uniform sampler2D textureSampler;
uniform vec4 objectColor;
uniform vec3 lightDirection;
uniform float time;

in vec2 fragmentTexture;
in vec3 fragmentNormal;

out vec4 outputColor;

void main() {
	vec3 halfVector = normalize( vec3(0, 0, 1) + lightDirection );
	float shininess = 6;
	float strength = 1.2;
	vec3 ambientColor = vec3(0.3);
	vec3 lightDirectionColor = vec3(1.0);

	float diffuse = max(0.0, dot(fragmentNormal, lightDirection));
	float specular = max(0.0, dot(fragmentNormal, halfVector));
	if (diffuse == 0.0) specular = 0.0;
	else specular = pow(specular, shininess); // sharpen the highlight

	vec3 scatteredLight = ambientColor + lightDirectionColor * diffuse;
	vec3 reflectedLight = lightDirectionColor * specular * strength;

	vec3 rgb = min(objectColor.rgb * scatteredLight + reflectedLight, vec3(1.0));

	//outputColor = texture(textureSampler, fragmentTexture) * vec4(rgb, objectColor.a);
	float gray = texture(textureSampler, fragmentTexture).r;
	outputColor = gray * vec4(rgb, objectColor.a);
}` + z
