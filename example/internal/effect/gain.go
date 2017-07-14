package effect

import (
	"github.com/loov/audio"
	"github.com/loov/audio/example/internal/atomic2"
	"github.com/loov/audio/slice"
)

type Gain struct {
	Value   atomic2.Float64
	current float64
}

func NewGain(value float64) *Gain {
	gain := &Gain{}
	gain.Value.Set(value)
	return gain
}

func (gain *Gain) Process(buf audio.Buffer) error {
	target := gain.Value.Get()
	current := gain.current

	channelCount := buf.ChannelCount()
	if target == current {
		switch buf := buf.(type) {
		case *audio.BufferF32:
			slice.Scale32(buf.Interleaved(), float32(current))
		case *audio.BufferF64:
			slice.Scale64(buf.Interleaved(), float64(current))
		default:
			return audio.ErrUnknownBuffer
		}
		return nil
	}

	var active float64

	switch buf := buf.(type) {
	case *audio.BufferF32:
		data := buf.Interleaved()
		for i := 0; i < len(data); i += channelCount {
			scale := float32(active)
			for k := 0; k < channelCount; k++ {
				data[i+k] *= scale
			}
			active = (active + target) * 0.5
		}
	case *audio.BufferF64:
		data := buf.Interleaved()
		for i := 0; i < len(data); i += channelCount {
			scale := active
			for k := 0; k < channelCount; k++ {
				data[i+k] *= scale
			}
			active = (active + target) * 0.5
		}
	default:
		return audio.ErrUnknownBuffer
	}

	current = active

	if atomic2.AlmostEqual64(current, target) {
		gain.current = target
	} else {
		gain.current = current
	}

	return nil
}
