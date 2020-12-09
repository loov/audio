package sequencer

import (
	"github.com/loov/audio"
	"github.com/loov/audio/bufferutil"
)

type Sequencer struct {
	Metronome Metronome
	Pool      *Pool

	Temporary audio.Buffer
	Tracks    []*Track
}

func NewSequencer(bpm float32, beatsPerMeasure int) *Sequencer {
	sequencer := &Sequencer{}
	sequencer.Metronome.BeatsPerMinute = bpm
	sequencer.Metronome.BeatsPerMeasure = beatsPerMeasure
	sequencer.Pool = NewPool()
	return sequencer
}

func (sequencer *Sequencer) Process(buf audio.Buffer) error {
	bufferutil.Zero(buf)
	if sequencer.Temporary == nil {
		sequencer.Temporary = buf.DeepCopy()
	}

	sequencer.Metronome.Process(buf)
	for _, track := range sequencer.Tracks {
		track.Process(sequencer.Temporary)
		bufferutil.Add(buf, sequencer.Temporary)
	}

	return nil
}
