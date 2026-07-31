package main

import (
	"unsafe"

	"github.com/pspdev-go/pspsdk-go/psp/ctrl"
	"github.com/pspdev-go/pspsdk-go/psp/debugscreen"
	"github.com/pspdev-go/pspsdk-go/psp/display"
	"github.com/pspdev-go/pspsdk-go/psp/kernel"
	"github.com/pspdev-go/pspsdk-go/psp/wlan"
)

func onOff(value int32) string {
	if value != 0 {
		return "ON "
	}
	return "OFF"
}

func main() {
	debugscreen.Init()
	ctrl.SetSamplingCycle(0)
	ctrl.SetSamplingMode(ctrl.ModeDigital)

	var pad ctrl.SceCtrlData
	for {
		ctrl.PeekBufferPositive(&pad, 1)
		if pad.Buttons&ctrl.ButtonStart != 0 {
			break
		}

		switchState := wlan.SceWlanGetSwitchState()
		powerState := wlan.SceWlanDevIsPowerOn()

		debugscreen.SetXY(0, 0)
		debugscreen.PutString("Wi-Fi Information Example\n\n")
		debugscreen.Printf("WLAN switch: %s\n", onOff(switchState))
		debugscreen.Printf("WLAN power:  %s\n", onOff(powerState))

		if switchState == 0 {
			debugscreen.PutString("\nTurn on the WLAN switch to read the adapter.")
			debugscreen.PutString("        ")
		} else {
			var macAddress [6]byte
			result := wlan.SceWlanGetEtherAddr(unsafe.Pointer(&macAddress[0]))
			if result < 0 {
				debugscreen.Printf("\nMAC address error: %d                 ", result)
			} else {
				debugscreen.Printf(
					"\nMAC address: %02X:%02X:%02X:%02X:%02X:%02X",
					macAddress[0], macAddress[1], macAddress[2],
					macAddress[3], macAddress[4], macAddress[5],
				)
			}
		}

		debugscreen.PutString("\n\nPress START to exit.")
		display.WaitVblankStart()
	}

	kernel.ExitGame()
}
