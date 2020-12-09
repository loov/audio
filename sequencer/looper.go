package sequencer

import (
	"io"

	"github.com/loov/audio"
)

type Looper struct{ audio.ReadSeeker }

func (reader *Looper) Read(buf audio.Buffer) (int, error) {
	totalFrameCount := 0
	dst := buf.ShallowCopy()
	for !dst.Empty() {
		frameCount, err := reader.ReadSeeker.Read(dst)
		dst.CutLeading(frameCount)
		totalFrameCount += frameCount
		if err != nil {
			if err == io.EOF {
				_, err := reader.ReadSeeker.Seek(0, 0)
				if err != nil {
					return totalFrameCount, err
				}
				continue
			}
			return totalFrameCount, err
		}
	}
	return totalFrameCount, nil
}
