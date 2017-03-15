#version 410

uniform vec4 uColor;

in float vSize;
out vec4 oColor;

float rand(vec2 co) {
	return fract(sin(dot(co.xy ,vec2(12.9898,78.233))) * 43758.5453);
}

void main() {
	float r = rand(uColor.rg * vSize);
	oColor = vec4(uColor.rgb * (vSize + r), uColor.a);
}
