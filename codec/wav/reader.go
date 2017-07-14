package wav

import (
	"fmt"
	"io"
	"time"

	"github.com/loov/audio"
	"github.com/loov/audio/codec/wav/wavdata"
	"github.com/loov/audio/slice"
)

const splitBufferSize = 1 << 10

type Reader struct {
	header header
	format format
	data   data
	head   int

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

func (reader *Reader) Clone() *Reader {
	clone := &Reader{}
	*clone = *reader
	return clone
}

func (reader *Reader) FrameCount() int { return len(reader.data.Data) / int(reader.format.BlockAlign) }
func (reader *Reader) SampleRate() int { return int(reader.format.SampleRate) }

func (reader *Reader) Duration() time.Duration {
	return audio.FrameCountToDuration(reader.FrameCount(), reader.SampleRate())
}

func (reader *Reader) Seek(offset time.Duration, whence int) (time.Duration, error) {
	offsetFrames := audio.DurationToFrameCount(offset, reader.SampleRate())
	reader.head = offsetFrames * int(reader.format.BlockAlign)
	return reader.Position(), nil
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

func (reader *Reader) ReadBlock(block []float32) (frameCount int) {
	if reader.framesLeft() == 0 {
		return 0
	}

	maxFrames := len(block) / reader.ChannelCount()
	if framesLeft := reader.framesLeft(); maxFrames > framesLeft {
		maxFrames = framesLeft
	}
	maxSamples := maxFrames * reader.ChannelCount()

	codec, ok := wavdata.Codecs[wavdata.Format{
		reader.format.Encoding,
		reader.format.BitsPerSample,
	}]

	if !ok {
		panic(fmt.Sprintf("unsupported codec %v", codec))
	}
	reader.head += codec.ReadF32(reader.data.Data[reader.head:], block, maxSamples)
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
			framesRead := reader.ReadBlock(reader._buffered)
			reader.buffered = reader._buffered[:framesRead*channelCount]
		}

		frameCount := slice.CopySliceTo32(reader.ChannelCount(), reader.buffered, dst)
		reader.buffered = reader.buffered[frameCount*channelCount:]
		dst.CutLeading(frameCount)

		totalFrameCount += frameCount
	}

	if reader.framesLeft() == 0 {
		return totalFrameCount, io.EOF
	}
	return totalFrameCount, nil
}
