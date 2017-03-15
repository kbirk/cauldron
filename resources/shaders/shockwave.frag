#version 410

uniform vec4 uColor;
uniform float uTime;

out vec4 oColor;

void main() {
	float intensity = max(0, 1.0 - (uTime  / 0.5));
	oColor = vec4(uColor.rgb * intensity, uColor.a * intensity);
}
