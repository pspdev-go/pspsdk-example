package main

import (
	"github.com/pspdev-go/pspsdk-go/psp/ctrl"
	"github.com/pspdev-go/pspsdk-go/psp/debugscreen"
	"github.com/pspdev-go/pspsdk-go/psp/display"
	"github.com/pspdev-go/pspsdk-go/psp/kernel"
)

func main() {
	debugscreen.Init()
	debugscreen.SetTextColor(0xffffffff)
	debugscreen.PutString("Hello from Go on PSP!\n\n")
	debugscreen.PutString("Press START to exit.")

	ctrl.SetSamplingCycle(0)
	ctrl.SetSamplingMode(ctrl.ModeDigital)

	var pad ctrl.SceCtrlData
	for {
		ctrl.PeekBufferPositive(&pad, 1)
		if pad.Buttons&ctrl.ButtonStart != 0 {
			break
		}
		display.WaitVblankStart()
	}

	kernel.ExitGame()
}
