package bufferutil

import (
	"github.com/loov/audio"
	"github.com/loov/audio/slice"
)

func Zero(buf audio.Buffer) {
	switch buf := buf.(type) {
	case *audio.BufferF32:
		slice.Zero32(buf.InternalBuffer())
	case *audio.BufferF64:
		slice.Zero64(buf.InternalBuffer())
	default:
		panic("unknown buffer")
	}
}

func Add(dst, src audio.Buffer) {
	// TODO: handle combinations
	switch dst := dst.(type) {
	case *audio.BufferF32:
		src := src.(*audio.BufferF32)
		slice.Add32(dst.InternalBuffer(), src.InternalBuffer())
	case *audio.BufferF64:
		src := src.(*audio.BufferF64)
		slice.Add64(dst.InternalBuffer(), src.InternalBuffer())
	default:
		panic("unknown buffer")
	}
}

func Scale(buf audio.Buffer, amount float64) {
	switch buf := buf.(type) {
	case *audio.BufferF32:
		slice.Scale32(buf.InternalBuffer(), float32(amount))
	case *audio.BufferF64:
		slice.Scale64(buf.InternalBuffer(), amount)
	default:
		panic("unknown buffer")
	}
}
