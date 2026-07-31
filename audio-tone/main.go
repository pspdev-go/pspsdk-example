package main

import (
	"unsafe"

	"github.com/pspdev-go/pspsdk-go/psp/audio"
	"github.com/pspdev-go/pspsdk-go/psp/ctrl"
	"github.com/pspdev-go/pspsdk-go/psp/debugscreen"
	"github.com/pspdev-go/pspsdk-go/psp/kernel"
)

const (
	sampleRate  = 44100
	sampleCount = 1024
	maxVolume   = 0x8000
)

var samples [sampleCount * 2]int16 // Two interleaved channels: left, right.

func fillSquareWave(frequency int, phase *int) {
	for i := 0; i < sampleCount; i++ {
		sample := int16(7000)
		if *phase >= sampleRate/2 {
			sample = -sample
		}
		samples[i*2] = sample
		samples[i*2+1] = sample

		// Advancing by frequency once per sample produces that many complete
		// cycles per second. The explicit wrap avoids integer division and
		// modulo instructions, which are problematic on Allegrex.
		*phase += frequency
		if *phase >= sampleRate {
			*phase -= sampleRate
		}
	}
}

func main() {
	debugscreen.Init()
	debugscreen.PutString("Audio Tone Example\n\n")
	debugscreen.PutString("UP / DOWN: change pitch\n")
	debugscreen.PutString("START: exit\n")

	ctrl.SetSamplingCycle(0)
	ctrl.SetSamplingMode(ctrl.ModeDigital)

	// -1 asks the system to select a free hardware audio channel.
	channel := audio.SceAudioChReserve(-1, sampleCount, audio.PSP_AUDIO_FORMAT_STEREO)
	if channel < 0 {
		debugscreen.Printf("\nCould not reserve audio channel: %d", channel)
		waitForStart()
		kernel.ExitGame()
		return
	}

	frequency := 440
	phase := 0
	var pad ctrl.SceCtrlData
	var previous ctrl.Buttons
	for {
		ctrl.PeekBufferPositive(&pad, 1)
		if pad.Buttons&ctrl.ButtonStart != 0 {
			break
		}
		if pad.Buttons&ctrl.ButtonUp != 0 && previous&ctrl.ButtonUp == 0 {
			frequency += 55
		}
		if pad.Buttons&ctrl.ButtonDown != 0 && previous&ctrl.ButtonDown == 0 && frequency > 110 {
			frequency -= 55
		}
		previous = pad.Buttons

		debugscreen.SetXY(0, 5)
		debugscreen.Printf("Frequency: %4d Hz", frequency)
		fillSquareWave(frequency, &phase)
		audio.SceAudioOutputBlocking(channel, maxVolume, unsafe.Pointer(&samples[0]))
	}

	audio.SceAudioChRelease(channel)
	kernel.ExitGame()
}

func waitForStart() {
	var pad ctrl.SceCtrlData
	for pad.Buttons&ctrl.ButtonStart == 0 {
		ctrl.PeekBufferPositive(&pad, 1)
	}
}
