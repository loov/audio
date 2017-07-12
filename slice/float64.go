package slice

func Zero64(data []float64) {
	for i := range data {
		data[i] = 0
	}
}

func Add64(dst, src []float64) {
	for i := range dst {
		dst[i] += src[i]
	}
}

func Scale64(data []float64, v float64) {
	for i := range data {
		data[i] *= v
	}
}

func ScaleLinearLerp64(data []float64, from, to float64) {
	inc := (to - from) / float64(len(data))
	for i := range data {
		data[i] *= from
		from += inc
	}
}

func Equal64(a, b []float64) bool {
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
