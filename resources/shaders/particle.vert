#version 410

layout(location=0) in vec2 aPosition;
layout(location=1) in vec2 aOffset;
layout(location=2) in vec2 aVelocity;
layout(location=3) in float aSize;

uniform mat4 uModel;
uniform mat4 uView;
uniform mat4 uProjection;
uniform float uTime;
uniform vec2 uGravity;

out float vSize;

void main() {
	vec2 displacement = (aVelocity * uTime) + (0.5 * uGravity * aSize * (uTime*uTime));
	float size = max(0, aSize - (aSize * uTime));
	vec2 wPosition = (aPosition * size) + aOffset + displacement;
	gl_Position = uProjection * uView * uModel * vec4(wPosition, 0, 1);
	vSize = size / 4;
}
