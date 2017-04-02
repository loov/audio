package main

import (
	"flag"
	"fmt"
	"os"
	"sandbox/raintree/qpc"

	"github.com/loov/audio"
	"github.com/loov/audio/bufferutil"
	"github.com/loov/audio/codec/wav"
	"github.com/loov/audio/native"
	"github.com/loov/audio/sequencer"
)

var (
	loop = flag.Bool("loop", false, "")
)

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	output, err := native.NewOutputDevice(audio.DeviceInfo{
		ChannelCount:      2,
		SampleRate:        44100,
		SamplesPerChannel: 128,
		//SamplesPerChannel: 44100,
	})
	check(err)

	info := output.OutputInfo()
	buf := audio.NewBufferF32Frames(
		audio.Format{
			ChannelCount: info.ChannelCount,
			SampleRate:   info.SampleRate,
		}, info.SamplesPerChannel)

	seq := sequencer.NewSequencer(105, 4)
	seq.Metronome.Sound.Enabled = true
	seq.Metronome.Sound.MasterVolume = 0.001
	seq.Pool.LoadDirectory("~samples\\")

	NewTrack := func(samples []*wav.Reader) *sequencer.Track {
		track := &sequencer.Track{}
		track.Volume = 1
		track.Sequencer = seq
		track.Synced = true
		seq.Tracks = append(seq.Tracks, track)
		track.Samples = samples
		return track
	}

	Bass := NewTrack(seq.Pool.Subset("Bass"))
	Bass.Volume = 0.4
	BassPluck := NewTrack(seq.Pool.Subset("Bass Pluck"))
	BassPluck.Volume = 0
	Cello := NewTrack(seq.Pool.Subset("Cello"))
	Cello.Volume = 0
	LeadFlute := NewTrack(seq.Pool.Subset("Lead Flute"))
	LeadFlute.Volume = 0.4
	ViolinA0 := NewTrack(seq.Pool.Subset("Violin A0"))
	ViolinA0.Volume = 1
	ViolinA1 := NewTrack(seq.Pool.Subset("Violin A1"))
	ViolinA1.Volume = 0.5
	ViolinB0 := NewTrack(seq.Pool.Subset("Violin B0"))
	ViolinB0.Volume = 0
	ViolinB1 := NewTrack(seq.Pool.Subset("Violin B1"))
	ViolinB1.Volume = 0

	for {
		start := qpc.Now()
		if err := seq.Process(buf); err != nil {
			check(err)
		}
		finish := qpc.Now()
		duration := finish.Sub(start).Duration()
		if duration > buf.Duration() {
			fmt.Println("SKIP ", duration, " expected ", buf.Duration())
		}

		bufferutil.Scale(buf, 10)

		if _, err := output.Write(buf); err != nil {
			check(err)
		}
	}
}
