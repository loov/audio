package sequencer

import (
	"github.com/loov/audio"
	"github.com/loov/audio/bufferutil"
)

type Sequencer struct {
	Metronome Metronome
}

func NewSequencer(bpm float32, beatsPerMeasure int) *Sequencer {
	sequencer := &Sequencer{}
	sequencer.Metronome.BeatsPerMinute = bpm
	sequencer.Metronome.BeatsPerMeasure = beatsPerMeasure
	return sequencer
}

func (sequencer *Sequencer) Process(buf audio.Buffer) error {
	bufferutil.Zero(buf)
	sequencer.Metronome.Process(buf)
	return nil
}
