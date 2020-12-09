package slice

import "github.com/loov/audio"

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

func CopyInterleaved64(nchan int, buf []float64, dst audio.Buffer) (frameCount int) {
	if dst.ChannelCount() != nchan {
		panic("channel count does not match")
	}

	maxSamples := len(buf)
	if dst.SampleCount() < maxSamples {
		maxSamples = dst.SampleCount()
	}

	switch dst := dst.(type) {
	case *audio.BufferF32:
		main := dst.Interleaved()
		for i := range main {
			main[i] = float32(buf[i])
		}
	case *audio.BufferF64:
		copy(dst.Interleaved(), buf[:maxSamples])
	default:
		panic("missing")
	}

	return maxSamples / nchan
}
