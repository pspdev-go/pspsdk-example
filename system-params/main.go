package main

import (
	"unsafe"

	"github.com/pspdev-go/pspsdk-go/psp/ctrl"
	"github.com/pspdev-go/pspsdk-go/psp/debugscreen"
	"github.com/pspdev-go/pspsdk-go/psp/display"
	"github.com/pspdev-go/pspsdk-go/psp/kernel"
	"github.com/pspdev-go/pspsdk-go/psp/utility_sysparam"
)

// IDs defined by psputility_sysparam.h.
const (
	paramNickname         = 1
	paramDateFormat       = 4
	paramTimeFormat       = 5
	paramTimezone         = 6
	paramLanguage         = 8
	paramButtonPreference = 9
)

func getInt(id int32) (int32, int32) {
	var value int32
	result := utility_sysparam.SceUtilityGetSystemParamInt(id, unsafe.Pointer(&value))
	return value, result
}

func main() {
	debugscreen.Init()
	debugscreen.PutString("System Parameters Example\n\n")

	var nickname [128]byte
	result := utility_sysparam.SceUtilityGetSystemParamString(
		paramNickname, unsafe.Pointer(&nickname[0]), int32(len(nickname)),
	)
	if result < 0 {
		debugscreen.Printf("Nickname error: %d\n", result)
	} else {
		end := 0
		for end < len(nickname) && nickname[end] != 0 {
			end++
		}
		debugscreen.Printf("Nickname:   %s\n", string(nickname[:end]))
	}

	printParam("Language", paramLanguage)
	printParam("Date format", paramDateFormat)
	printParam("Time format", paramTimeFormat)
	printParam("Timezone", paramTimezone)
	printParam("Enter button", paramButtonPreference)
	debugscreen.PutString("\nValues use the enums from psputility_sysparam.h.")
	debugscreen.PutString("\nPress START to exit.")

	ctrl.SetSamplingCycle(0)
	ctrl.SetSamplingMode(ctrl.ModeDigital)
	var pad ctrl.SceCtrlData
	for pad.Buttons&ctrl.ButtonStart == 0 {
		ctrl.PeekBufferPositive(&pad, 1)
		display.WaitVblankStart()
	}

	kernel.ExitGame()
}

func printParam(label string, id int32) {
	value, result := getInt(id)
	if result < 0 {
		debugscreen.Printf("%s error: %d\n", label, result)
		return
	}
	debugscreen.Printf("%-12s %d\n", label+":", value)
}
