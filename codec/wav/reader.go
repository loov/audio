package wav

import (
	"io"
	"time"

	"github.com/loov/audio"
	"github.com/loov/audio/slice"
)

const splitBufferSize = 1 << 10

type Reader struct {
	header header
	format format
	data   data

	head int // in data

	buffered  []float32
	_buffered []float32
}

func NewBytesReader(data []byte) (*Reader, error) {
	reader := &Reader{}
	reader._buffered = make([]float32, splitBufferSize)
	return reader, reader.decode(data)
}

func (reader *Reader) decode(data []byte) error {
	// TODO: handle short buffer
	data = reader.header.Read(data)
	data = reader.format.Read(data)
	data = reader.data.Read(data)
	return nil
}

func (reader *Reader) FrameCount() int { return len(reader.data.Data) / int(reader.format.BlockAlign) }
func (reader *Reader) SampleRate() int { return int(reader.format.SampleRate) }

func (reader *Reader) Duration() time.Duration {
	return audio.FrameCountToDuration(reader.FrameCount(), reader.SampleRate())
}

func (reader *Reader) Position() time.Duration {
	frameHead := reader.head / int(reader.format.BlockAlign)
	return audio.FrameCountToDuration(frameHead, reader.SampleRate())
}

func (reader *Reader) ChannelCount() int {
	return int(reader.format.NumChannels)
}

func (reader *Reader) framesLeft() int {
	return (len(reader.data.Data) - reader.head) / int(reader.format.BlockAlign)
}

func (reader *Reader) ReadInterleavedBlock(block []float32) (frameCount int) {
	if reader.framesLeft() == 0 {
		return 0
	}

	maxFrames := len(block) / reader.ChannelCount()
	if framesLeft := reader.framesLeft(); maxFrames > framesLeft {
		maxFrames = framesLeft
	}
	maxSamples := maxFrames * reader.ChannelCount()

	h, src := reader.head, reader.data.Data

	switch reader.format.BitsPerSample {
	case 8: // unsigned 8bit
		for k := 0; k < maxSamples; k++ {
			v := uint8(src[h])
			block[k] = float32(v)/128.0 - 1.0
			h += 1
		}
	case 16: // signed 16bit; -32,768 (0x7FFF) to 32,767 (0x8000)
		for k := 0; k < maxSamples; k++ {
			v := int16(src[h]) | int16(src[h+1])<<8
			block[k] = float32(v) / float32(0x8000)
			h += 2
		}
	case 32: // signed 32bit
		h := reader.head
		for k := 0; k < maxSamples; k++ {
			v := int32(src[h]) | int32(src[h+1])<<8 | int32(src[h+2])<<8 | int32(src[h+3])<<8
			block[k] = float32(v) / float32(0x80000000)
			h += 4
		}
	default:
		panic("unimplemented bits per sample")
	}
	reader.head = h

	return maxFrames
}

func (reader *Reader) Read(buf audio.Buffer) (int, error) {
	if reader.framesLeft() == 0 {
		return 0, io.EOF
	}

	channelCount := reader.ChannelCount()
	totalFrameCount := 0
	dst := buf.ShallowCopy()
	for !dst.Empty() {
		if len(reader.buffered) <= 0 {
			n := reader.ReadInterleavedBlock(reader._buffered)
			reader.buffered = reader._buffered[:n]
		}

		frameCount := slice.Split32(reader.ChannelCount(), reader.buffered, dst)
		reader.buffered = reader.buffered[frameCount*channelCount:]
		dst.CutLeading(frameCount)

		totalFrameCount += frameCount
	}

	return totalFrameCount, nil
}
