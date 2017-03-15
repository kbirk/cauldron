#version 410

layout(location=0) in vec2 aPosition;

uniform mat4 uModel;
uniform mat4 uView;
uniform mat4 uProjection;
uniform float uForce;
uniform float uTime;

float easeOut(float t) {
	t -= 1.0;
	return 1.0 + t*t*t*t*t;
}

void main() {
	vec2 wPosition = vec2(
		aPosition.x * uForce * easeOut(uTime),
		aPosition.y * uForce * easeOut(uTime) / 3);
	gl_Position = uProjection * uView * uModel * vec4(wPosition, 0, 1);
}
