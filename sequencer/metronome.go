package sequencer

import (
	"math/rand"
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
}

func (metronome *Metronome) UpdateSampleRate(sampleRate int) {
	if metronome.SampleRate == sampleRate {
		return
	}

	metronome.SampleRate = sampleRate
	// calculate measure length
	beatsPerSecond := metronome.BeatsPerMinute / 60
	metronome.Measure.Duration = time.Duration(float32(time.Second) / beatsPerSecond)
	metronome.Measure.Length = audio.DurationToFrameCount(metronome.Measure.Duration, metronome.SampleRate)

	// update calculated fields
	metronome.Measure.Frame = audio.DurationToFrameCount(metronome.Measure.Time, metronome.SampleRate)
	metronome.Measure.Next = metronome.Measure.Length - metronome.Measure.Frame
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

	scale := float32(0.0)

	var saws [8]float32
	baserise := float32(0.002)
	basescale := float32(1.6)

	var lastsaw float32
	//decay := 1 * float32(buf.SampleRate())
	generate.MonoF32(buf, func() float32 {
		//scale -= 0.0008
		scale -= 0.00005
		if k <= 0 {
			scale = 1.0
			baserise = 0.002 + rand.Float32()*0.01
			basescale = 1.0 + rand.Float32()
			k += metronome.Measure.Length
		}
		k--
		if scale < 0 {
			lastsaw = 0.0
			return 0
		}

		saw := float32(0.0)
		rise := baserise
		vol := float32(1.0)
		for i := range saws {
			saws[i] += rise
			if saws[i] > 1.0 {
				saws[i] -= 2.0
			}
			saw += saws[i] * vol
			vol *= 0.7
			rise *= basescale
		}
		saw = lastsaw*0.7 + saw*0.3
		lastsaw = saw

		return (saw + rand.Float32()*0.1 - 0.05) * scale
	})

	metronome.Advance(buf)

	return nil
}
