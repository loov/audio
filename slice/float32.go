package slice

func Zero32(data []float32) {
	for i := range data {
		data[i] = 0
	}
}

func Add32(dst, src []float32) {
	for i := range dst {
		dst[i] += src[i]
	}
}

func Scale32(data []float32, v float32) {
	for i := range data {
		data[i] *= v
	}
}

func ScaleLinearLerp32(data []float32, from, to float32) {
	inc := (to - from) / float32(len(data))
	for i := range data {
		data[i] *= from
		from += inc
	}
}

func Equal32(a, b []float32) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
