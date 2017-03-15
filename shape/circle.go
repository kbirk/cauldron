package shape

import (
	"math"
)

// Circle returns an interleave vertex array and indices array.
func Circle(radius float32, segments int, positions bool, uvs bool) ([]float32, []uint16) {
	var p []float32
	if positions {
		p = circlePositions(radius, segments)
	}
	var u []float32
	if uvs {
		u = circleUVs(segments)
	}
	return interleave(p, u), circleIndices(segments)
}

func circlePositions(radius float32, segments int) []float32 {
	theta := (2 * math.Pi) / float32(segments)
	// precalculate sine and cosine
	c := float32(math.Cos(float64(theta)))
	s := float32(math.Sin(float64(theta)))
	// start at angle = 0
	x := radius
	y := float32(0.0)
	length := (segments + 2) * 3
	// create positions
	positions := make([]float32, length)
	// first point
	positions[0] = 0.0
	positions[1] = 0.0
	positions[2] = 0.0
	// last point
	positions[length-3] = radius
	positions[length-2] = 0
	positions[length-1] = 0
	// intermediate
	for i := 0; i < segments; i++ {
		positions[(i+1)*3] = x
		positions[(i+1)*3+1] = y
		positions[(i+1)*3+2] = 0.0
		// apply the rotation
		t := x
		x = c*x - s*y
		y = s*t + c*y
	}
	return positions
}

func circleUVs(segments int) []float32 {
	theta := (2 * math.Pi) / float32(segments)
	// precalculate sine and cosine
	c := float32(math.Cos(float64(theta)))
	s := float32(math.Sin(float64(theta)))
	// start at angle = 0
	x := float32(1.0)
	y := float32(0.0)
	length := (segments + 2) * 2
	// create uvs
	uvs := make([]float32, length)
	// first point
	uvs[0] = 0.0
	uvs[1] = 0.0
	// last point
	uvs[length-2] = 1.0
	uvs[length-1] = 0
	// intermediate
	for i := 0; i < segments; i++ {
		uvs[(i+1)*2] = x
		uvs[(i+1)*2+1] = y
		// apply the rotation
		t := x
		x = c*x - s*y
		y = s*t + c*y
	}
	return uvs
}

func circleIndices(segments int) []uint16 {
	indices := make([]uint16, segments*3)
	for i := 0; i < segments; i++ {
		indices[i*3] = 0
		indices[i*3+1] = uint16(i + 1)
		indices[i*3+2] = uint16(i + 2)
	}
	return indices
}
