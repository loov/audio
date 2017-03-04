package wav

import (
	"time"

	"github.com/loov/audio"
)

type Reader struct {
	header header
	format format
	data   data
}

func NewBytesReader(data []byte) (*Reader, error) {
	reader := &Reader{}
	return reader, reader.decode(data)
}

func (reader *Reader) decode(data []byte) error {
	// TODO: handle short buffer
	data = reader.header.Read(data)
	data = reader.format.Read(data)
	data = reader.data.Read(data)
	return nil
}

func (reader *Reader) Duration() time.Duration {
	return 0
}

func (reader *Reader) Position() time.Duration {
	return 0
}

func (reader *Reader) Read(buf audio.Buffer) (int, error) {
	return buf.FrameCount(), nil
}
