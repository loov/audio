package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/loov/audio"
	"github.com/loov/audio/codec/wav"
	"github.com/loov/audio/native"
)

var (
	loop = flag.Bool("loop", false, "")
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
		source = &Looper{reader}
	}

	pipe := audio.Pipe{source, nil, output}
	check(pipe.Run(buf))
}
