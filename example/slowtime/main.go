package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"

	"github.com/loov/audio"
	"github.com/loov/audio/example/internal/effect"
	"github.com/loov/audio/example/internal/osc"
	"github.com/loov/audio/native"
)

func main() {
	sine := &osc.Sine{}
	sine.Frequency.Set(440.0)

	gain := &effect.Gain{}
	gain.Value.Set(0.5)

	go func() {
		runtime.LockOSThread()

		device, err := native.NewOutputDevice(audio.DeviceInfo{
			ChannelCount:      2,
			SampleRate:        44100,
			SamplesPerChannel: 64,
		})
		if err != nil {
			log.Fatal("failed to create device: ", err)
		}
		defer device.Close()
		info := device.OutputInfo()

		work := audio.NewBufferF32Frames(audio.Format{
			ChannelCount: info.ChannelCount,
			SampleRate:   info.SampleRate,
		}, info.SamplesPerChannel)

		for {
			sine.Process(work)
			gain.Process(work)
			if err := device.Process(work); err != nil {
				log.Fatal("failed to write audio frame: ", err)
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if len(scanner.Text()) > 0 {
			k := scanner.Text()[0]
			switch k {
			case 'q':
				return
			case '+':
				gain.Value.Set(gain.Value.Get() + 0.10)
			case '-':
				gain.Value.Set(gain.Value.Get() - 0.10)
			default:
				v := float64(math.Abs(float64(int(k - 100))))
				note := 440.0 * math.Pow(2, (v)/12.0)
				if note > 22000 {
					note = 440.0
				}

				fmt.Printf("switching oscillator to %.2f Hz\n", note)
				sine.Frequency.Set(note)
			}
		}
	}
}
