package shape

func calcStride(positions bool, uvs bool) int {
	stride := 0
	if positions {
		stride += 3
	}
	if uvs {
		stride += 2
	}
	return stride
}

func interleave(positions []float32, uvs []float32) []float32 {
	stride := calcStride(positions != nil, uvs != nil)
	length := (len(positions) / 3)
	buff := make([]float32, length*stride)
	for i := 0; i < length; i++ {
		if positions != nil {
			buff[i*stride] = positions[i*3]
			buff[i*stride+1] = positions[i*3+1]
			buff[i*stride+2] = positions[i*3+2]
		}
		if uvs != nil {
			buff[i*stride+3] = uvs[i*2]
			buff[i*stride+4] = uvs[i*2+1]
		}
	}
	return buff
}
