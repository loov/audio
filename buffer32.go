package audio

import "time"

type BufferF32 struct {
	format Format
	offset uint32
	frames uint32
	data   []float32
}

func NewBufferF32(format Format, duration time.Duration) *BufferF32 {
	return NewBufferF32Frames(format, format.FrameCount(duration))
}

func NewBufferF32Frames(format Format, frames int) *BufferF32 {
	samples := format.ChannelCount * frames
	return &BufferF32{
		format: format,
		offset: 0,
		frames: uint32(frames),
		data:   make([]float32, samples, samples),
	}
}

func (b *BufferF32) InternalBuffer() []float32 { return b.data }

func (b *BufferF32) Interleaved() []float32 {
	start := int(b.offset)
	return b.data[start : start+b.SampleCount()]
}

func (b *BufferF32) SampleRate() int   { return b.format.SampleRate }
func (b *BufferF32) ChannelCount() int { return b.format.ChannelCount }

func (b *BufferF32) Empty() bool      { return b.frames == 0 }
func (b *BufferF32) FrameCount() int  { return int(b.frames) }
func (b *BufferF32) SampleCount() int { return int(b.frames) * b.ChannelCount() }

func (b *BufferF32) Duration() time.Duration {
	return time.Duration(int(time.Second) * b.FrameCount() / b.SampleRate())
}

func (b *BufferF32) ShallowCopy() Buffer {
	x := *b
	return &x
}

func (b *BufferF32) DeepCopy() Buffer {
	x := *b
	x.data = make([]float32, len(b.data), len(b.data))
	copy(x.data, b.data)
	return &x
}

func (b *BufferF32) Slice(low, high int) {
	b.offset += uint32(low * b.ChannelCount())
	b.frames = uint32((high - low) * b.ChannelCount())
}

func (b *BufferF32) CutLeading(low int) {
	b.offset += uint32(low * b.ChannelCount())
	b.frames -= uint32(low * b.ChannelCount())
}
