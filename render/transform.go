package render

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Transform represents a 3D transform.
type Transform struct {
	Rotation    float32
	Translation mgl32.Vec2
	Scale       mgl32.Vec2
}

// NewTransform instantiates and returns a new identity transform/
func NewTransform() *Transform {
	return &Transform{
		Rotation:    0,
		Translation: mgl32.Vec2{0, 0},
		Scale:       mgl32.Vec2{1, 1},
	}
}

// SetRotation sets the rotation of the transform.
func (t *Transform) SetRotation(angle float32) {
	t.Rotation = angle
}

// SetScale2 sets the scale of the transform.
func (t *Transform) SetScale2(scale mgl32.Vec2) {
	t.Scale = scale
}

// SetScale sets the scale of the transform.
func (t *Transform) SetScale(scale float32) {
	t.Scale = mgl32.Vec2{scale, scale}
}

// Rotate rotates the transform relative to its current rotation.
func (t *Transform) Rotate(angle float32) {
	t.Rotation += angle
}

// Translate translates the transform.
func (t *Transform) Translate(translation mgl32.Vec2) {
	t.Translation = t.Translation.Add(translation)
}

// Matrix returns the transform matrix.
func (t *Transform) Matrix() mgl32.Mat4 {
	return t.TranslationMatrix().Mul4(t.RotationMatrix()).Mul4(t.ScaleMatrix())
}

// ViewMatrix returns the view matrix.
func (t *Transform) ViewMatrix() mgl32.Mat4 {
	eye := mgl32.Vec3{
		t.Translation[0],
		t.Translation[1],
		0.0,
	}
	center := eye.Sub(mgl32.Vec3{0, 0, 1})
	up := mgl32.Vec3{0, 1, 0}
	return mgl32.LookAtV(eye, center, up)
}

// RotationMatrix returns the rotation matrix.
func (t *Transform) RotationMatrix() mgl32.Mat4 {
	return mgl32.HomogRotate3DZ(t.Rotation)
}

// TranslationMatrix returns the translation matrix.
func (t *Transform) TranslationMatrix() mgl32.Mat4 {
	return mgl32.Translate3D(
		t.Translation[0],
		t.Translation[1],
		0.0)
}

// ScaleMatrix returns the scale matrix.
func (t *Transform) ScaleMatrix() mgl32.Mat4 {
	return mgl32.Scale3D(
		t.Scale[0],
		t.Scale[1],
		0.0)
}
