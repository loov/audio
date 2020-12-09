package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/loov/audio"
	"github.com/loov/audio/bufferutil"
	"github.com/loov/audio/native"
	"github.com/loov/audio/sequencer"
	"github.com/loov/hrtime"
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
		SamplesPerChannel: 256,
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

	NewTrack := func(name string) *sequencer.Track {
		track := &sequencer.Track{}
		track.Name = name
		track.Volume = 1
		track.Sequencer = seq
		track.Synced = true
		seq.Tracks = append(seq.Tracks, track)
		track.Samples = seq.Pool.Subset(name)
		return track
	}

	Bass := NewTrack("Bass")
	Bass.Volume = 0
	BassPluck := NewTrack("Bass Pluck")
	BassPluck.Volume = 1
	Cello := NewTrack("Cello")
	Cello.Volume = 0
	LeadFlute := NewTrack("Lead Flute")
	LeadFlute.Volume = 0
	ViolinA0 := NewTrack("Violin A0")
	ViolinA0.Volume = 1
	ViolinA1 := NewTrack("Violin A1")
	ViolinA1.Volume = 0.3
	ViolinB0 := NewTrack("Violin B0")
	ViolinB0.Volume = 0
	ViolinB1 := NewTrack("Violin B1")
	ViolinB1.Volume = 0

	for {
		start := hrtime.Now()
		if err := seq.Process(buf); err != nil {
			check(err)
		}
		finish := hrtime.Now()
		duration := finish - start
		if duration > buf.Duration() {
			fmt.Println("SKIP ", duration, " expected ", buf.Duration())
		}

		bufferutil.Scale(buf, 10)

		if _, err := output.Write(buf); err != nil {
			check(err)
		}
	}
}
