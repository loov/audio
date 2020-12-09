package audio

import (
	"errors"
	"time"
)

var (
	ErrUnknownBuffer = errors.New("unknown buffer")
)

type Buffer interface {
	SampleRate() int
	// Interleaved() []T
	ChannelCount() int
	FrameCount() int
	SampleCount() int
	Duration() time.Duration
	Empty() bool

	ShallowCopy() Buffer
	DeepCopy() Buffer

	Slice(low, high int)
	CutLeading(low int)
}

type Writer interface {
	Write(Buffer) (int, error)
}

type Reader interface {
	Read(Buffer) (int, error)
}

type ReadSeeker interface {
	Reader
	Seeker
}

type Processor interface {
	Process(Buffer) error
}

type Seeker interface {
	Position() time.Duration
	// Seek seeks to the nearest possible position and returns the final position
	Seek(offset time.Duration, whence int) (time.Duration, error)
}

// Anything implemeting Durationer has a known finite length
type Durationer interface {
	Duration() time.Duration
}
