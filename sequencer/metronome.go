package sequencer

import (
	"time"

	"github.com/loov/audio"
	"github.com/loov/audio/generate"
)

type Metronome struct {
	BeatsPerMinute  float32
	BeatsPerMeasure int
	SampleRate      int

	Global struct {
		Time  time.Duration
		Frame int64
	}

	Measure struct {
		Length   int
		Duration time.Duration

		Time  time.Duration
		Frame int
		Next  int // frames until next frame
	}

	Sound struct {
		Enabled   bool
		Frequency float32

		Volume      float32
		VolumeDecay float32
		Phase       float32
		LastPhase   float32
		PhaseSpeed  float32
	}
}

func (metronome *Metronome) UpdateSampleRate(sampleRate int) {
	if metronome.SampleRate == sampleRate {
		return
	}

	metronome.SampleRate = sampleRate
	{
		// calculate measure length
		beatsPerSecond := metronome.BeatsPerMinute / 60
		metronome.Measure.Duration = time.Duration(float32(time.Second) / beatsPerSecond)
		metronome.Measure.Length = audio.DurationToFrameCount(metronome.Measure.Duration, metronome.SampleRate)

		// update calculated fields
		metronome.Measure.Frame = audio.DurationToFrameCount(metronome.Measure.Time, metronome.SampleRate)
		metronome.Measure.Next = metronome.Measure.Length - metronome.Measure.Frame
	}

	{
		metronome.Sound.Enabled = true
		metronome.Sound.Frequency = 880
		metronome.Sound.Volume = 0.0
		metronome.Sound.VolumeDecay = 1 / (0.1 * float32(metronome.SampleRate))
		metronome.Sound.Phase = 0.0
		metronome.Sound.PhaseSpeed = 2 * metronome.Sound.Frequency / float32(metronome.SampleRate)
	}
}

func (metronome *Metronome) Advance(buf audio.Buffer) {
	metronome.Global.Time += buf.Duration()
	metronome.Global.Frame += int64(buf.FrameCount())

	metronome.Measure.Frame += buf.FrameCount()
	for metronome.Measure.Frame > metronome.Measure.Length {
		metronome.Measure.Frame -= metronome.Measure.Length
	}

	metronome.Measure.Next = metronome.Measure.Length - metronome.Measure.Frame
	metronome.Measure.Time = audio.FrameCountToDuration(metronome.Measure.Frame, metronome.SampleRate)
}

func (metronome *Metronome) Process(buf audio.Buffer) error {
	metronome.UpdateSampleRate(buf.SampleRate())

	k := metronome.Measure.Next
	if metronome.Measure.Frame == 0 {
		k = 0
	}

	snd := &metronome.Sound
	generate.MonoF32(buf, func() float32 {
		snd.Volume -= snd.VolumeDecay
		if k <= 0 {
			snd.Volume = 1.0
			k += metronome.Measure.Length
		}

		if snd.Volume < 0 {
			return 0
		}

		snd.Phase += snd.PhaseSpeed
		if snd.Phase > 1 {
			snd.Phase -= 2.0
		}

		p := snd.LastPhase*0.5 + snd.Phase*0.5
		snd.LastPhase = p
		return p * snd.Volume
	})

	metronome.Advance(buf)

	return nil
}
