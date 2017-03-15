#version 410

uniform vec4 uColor;
uniform float uTime;

in float vOpacity;

out vec4 oColor;

float cube(float v) {
	return v*v*v;
}

void main() {
	float intensity = max(0, 1.0 - (uTime  / 0.5));
	float opacity = cube(vOpacity);
	oColor = vec4(uColor.rgb * intensity, uColor.a * intensity * opacity);
}
