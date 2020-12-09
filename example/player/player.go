package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/loov/audio"
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

	files := flag.Args()
	if len(files) == 0 {
		flag.PrintDefaults()
		os.Exit(2)
	}

	data, err := ioutil.ReadFile(files[0])
	check(err)

	reader, err := wav.NewBytesReader(data)
	check(err)

	output, err := native.NewOutputDevice(audio.DeviceInfo{
		ChannelCount:      reader.ChannelCount(),
		SampleRate:        reader.SampleRate(),
		SamplesPerChannel: 128,
	})
	check(err)

	info := output.OutputInfo()
	buf := audio.NewBufferF32Frames(
		audio.Format{
			ChannelCount: info.ChannelCount,
			SampleRate:   info.SampleRate,
		}, info.SamplesPerChannel)

	fmt.Printf("\n\n\n")
	fmt.Println(reader.Duration())

	var source audio.Reader = reader
	if *loop {
		source = &sequencer.Looper{reader}
	}

	pipe := audio.Pipe{source, nil, output}
	check(pipe.Run(buf))
}
