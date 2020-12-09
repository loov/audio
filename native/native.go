package native

import (
	"github.com/loov/audio"
	"github.com/loov/audio/native/portaudio"
)

func NewOutputDevice(pref audio.DeviceInfo) (audio.OutputDevice, error) {
	return portaudio.NewOutputDevice(pref)
}
