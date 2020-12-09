package slice

import "github.com/loov/audio"

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

func CopySliceTo32(nchan int, buf []float32, dst audio.Buffer) (frameCount int) {
	if dst.ChannelCount() != nchan {
		panic("channel count does not match")
	}

	maxSamples := (len(buf) / nchan) * nchan
	if dst.SampleCount() < maxSamples {
		maxSamples = dst.SampleCount()
	}

	switch dst := dst.(type) {
	case *audio.BufferF32:
		copy(dst.Interleaved(), buf[:maxSamples])
	case *audio.BufferF64:
		main := dst.Interleaved()
		for i := range main {
			main[i] = float64(buf[i])
		}
	default:
		panic("missing")
	}

	return maxSamples / nchan
}

func CopySliceFrom32(dst audio.Buffer, nchan int, buf []float32) (frameCount int) {
	if dst.ChannelCount() != nchan {
		panic("channel count does not match")
	}

	maxSamples := (len(buf) / nchan) * nchan
	if dst.SampleCount() < maxSamples {
		maxSamples = dst.SampleCount()
	}

	switch dst := dst.(type) {
	case *audio.BufferF32:
		copy(buf[:maxSamples], dst.Interleaved())
	case *audio.BufferF64:
		main := dst.Interleaved()
		for i := range main {
			buf[i] = float32(main[i])
		}
	default:
		panic("missing")
	}

	return maxSamples / nchan
}
