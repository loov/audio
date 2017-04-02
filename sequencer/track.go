package sequencer

import (
	"math/rand"

	"github.com/loov/audio"
	"github.com/loov/audio/bufferutil"
	"github.com/loov/audio/codec/wav"
)

type Track struct {
	Name string

	Sequencer *Sequencer
	Enabled   bool
	Volume    float32

	Synced  bool
	Samples []*wav.Reader

	Temporary audio.Buffer
	Playing   *Delay
	Pending   *Delay
}

type Delay struct {
	Frames int
	Reader *wav.Reader
	Done   bool
}

func (delay *Delay) Process(buf audio.Buffer) error {
	if delay.Done {
		return nil
	}

	if delay.Frames > 0 {
		if delay.Frames > buf.FrameCount() {
			delay.Frames -= buf.FrameCount()
			return nil
		}

		bufferutil.Zero(buf)
		buf = buf.ShallowCopy()
		buf.CutLeading(delay.Frames)
		delay.Frames = 0
	}

	_, err := delay.Reader.Read(buf)
	delay.Done = err != nil
	return nil
}

func (track *Track) Process(buf audio.Buffer) error {
	if track.Enabled {
		return nil
	}
	if track.Temporary == nil {
		track.Temporary = buf.DeepCopy()
	}

	if track.Pending == nil && track.Playing == nil {
		track.ScheduleNext()
		return nil
	}
	if track.Pending == nil && track.Playing == nil {
		bufferutil.Zero(buf)
		return nil
	}

	if track.Playing != nil {
		track.Playing.Process(buf)
	}
	if track.Pending != nil {
		if track.Pending.Frames > buf.FrameCount() {
			track.Pending.Process(buf)
		} else if track.Pending.Frames == 0 {
			track.Pending.Process(track.Temporary)
			bufferutil.Add(buf, track.Temporary)
		}
	}
	bufferutil.Scale(buf, float64(track.Volume))

	if track.Playing == nil || track.Playing.Done {
		track.Playing, track.Pending = track.Pending, nil
	}
	if track.Pending == nil || track.Pending.Done {
		track.Pending = nil
	}

	if track.Playing == nil || track.Pending == nil {
		track.ScheduleNext()
	}

	return nil
}

func (track *Track) ScheduleNext() {
	if track.Playing == nil {
		framesToNext := track.Sequencer.Metronome.Measure.Next
		n := rand.Intn(len(track.Samples) - 1)
		sample := track.Samples[n].Clone()
		track.Playing = &Delay{framesToNext, sample, false}
	}
}
