#version 410

uniform vec4 uColor;

in float vSize;
out vec4 oColor;

float rand(vec2 co){
	return fract(sin(dot(co.xy ,vec2(12.9898,78.233))) * 43758.5453);
}

void main() {
	float r = rand(uColor.rg * vSize) * 0.5;
	float factor = min(1.0, 0.2 * vSize);
	float intensity = max(0.4, 1.0 - factor);
	float alpha = max(0, 1.0 - factor);
	oColor = vec4(uColor.rgb * (intensity + r), uColor.a * alpha);
}
