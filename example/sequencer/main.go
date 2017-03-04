package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/loov/audio"
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
		ChannelCount: 2,
		SampleRate:   44100,
		//SamplesPerChannel: 128,
		SamplesPerChannel: 44100,
	})
	check(err)

	info := output.OutputInfo()
	buf := audio.NewBufferF32Frames(
		audio.Format{
			ChannelCount: info.ChannelCount,
			SampleRate:   info.SampleRate,
		}, info.SamplesPerChannel)

	seq := sequencer.NewSequencer(120, 4)
	for {
		if err := seq.Process(buf); err != nil {
			check(err)
		}
		if _, err := output.Write(buf); err != nil {
			check(err)
		}
	}
}
