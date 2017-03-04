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
