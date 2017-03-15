package shape

// Quad returns an interleave vertex array and indices array.
func Quad(size float32, positions bool, uvs bool) ([]float32, []uint16) {
	var p []float32
	if positions {
		p = quadPositions(size)
	}
	var u []float32
	if uvs {
		u = quadUVs()
	}
	return interleave(p, u), quadIndices()
}

func quadPositions(size float32) []float32 {
	half := size / 2.0
	return []float32{
		-half, -half, 0.0,
		half, -half, 0.0,
		half, half, 0.0,
		-half, half, 0.0,
	}
}

func quadUVs() []float32 {
	return []float32{
		0.0, 0.0,
		1.0, 0.0,
		1.0, 1.0,
		0.0, 1.0,
	}
}

func quadIndices() []uint16 {
	return []uint16{
		0, 1, 2, 0, 2, 3,
	}
}
