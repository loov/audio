package generate

import (
	"github.com/loov/audio"
)

func MonoF32(out audio.Buffer, sample func() float32) error {
	channelCount := out.ChannelCount()
	switch out := out.(type) {
	case *audio.BufferF32:
		main := out.Interleaved()
		for i := 0; i < len(main); i += channelCount {
			s := float32(sample())
			for k := 0; k < channelCount; k++ {
				main[i+k] = s
			}
		}
	case *audio.BufferF64:
		main := out.Interleaved()
		for i := 0; i < len(main); i += channelCount {
			s := float64(sample())
			for k := 0; k < channelCount; k++ {
				main[i+k] = s
			}
		}
	default:
		//TODO: implement the slowest path
		return audio.ErrUnknownBuffer
	}
	return nil
}

func StereoF32(out audio.Buffer, sample func() (float32, float32)) error {
	channelCount := out.ChannelCount()
	switch out := out.(type) {
	case *audio.BufferF32:
		if channelCount >= 2 {
			main := out.Interleaved()
			for i := 0; i < len(main); i += channelCount {
				leftsample, rightsample := sample()
				main[i], main[i+1] = float32(leftsample), float32(rightsample)
				for k := 2; k < channelCount; k++ {
					main[i+k] = 0
				}
			}
		} else {
			main := out.Interleaved()
			for i := range main {
				leftsample, rightsample := sample()
				main[i] = float32(leftsample+rightsample) * 0.5
			}
		}
	case *audio.BufferF64:
		if channelCount >= 2 {
			main := out.Interleaved()
			for i := 0; i < len(main); i += channelCount {
				leftsample, rightsample := sample()
				main[i], main[i+1] = float64(leftsample), float64(rightsample)
				for k := 2; k < channelCount; k++ {
					main[i+k] = 0
				}
			}
		} else {
			main := out.Interleaved()
			for i := range main {
				leftsample, rightsample := sample()
				main[i] = float64(leftsample+rightsample) * 0.5
			}
		}
	default:
		//TODO: implement the slowest path
		return audio.ErrUnknownBuffer
	}
	return nil
}
