package main

import (
	"github.com/pspdev-go/pspsdk-go/psp/ctrl"
	"github.com/pspdev-go/pspsdk-go/psp/debugscreen"
	"github.com/pspdev-go/pspsdk-go/psp/display"
	"github.com/pspdev-go/pspsdk-go/psp/kernel"
	"github.com/pspdev-go/pspsdk-go/psp/power"
)

func yesNo(value int32) string {
	if value != 0 {
		return "yes"
	}
	return "no"
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

		debugscreen.SetXY(0, 0)
		debugscreen.PutString("Power Status Example\n\n")
		debugscreen.Printf("External power:  %-3s\n", yesNo(power.ScePowerIsPowerOnline()))
		debugscreen.Printf("Battery present: %-3s\n", yesNo(power.ScePowerIsBatteryExist()))
		debugscreen.Printf("Charging:        %-3s\n", yesNo(power.ScePowerIsBatteryCharging()))
		debugscreen.Printf("Charge:          %3d %%\n", power.ScePowerGetBatteryLifePercent())
		debugscreen.Printf("Time remaining:  %3d minutes\n", power.ScePowerGetBatteryLifeTime())
		debugscreen.Printf("Temperature:     %3d C\n", power.ScePowerGetBatteryTemp())
		debugscreen.Printf("Voltage:         %4d mV\n", power.ScePowerGetBatteryVolt())
		debugscreen.Printf("CPU / bus clock: %3d / %3d MHz\n",
			power.ScePowerGetCpuClockFrequencyInt(),
			power.ScePowerGetBusClockFrequencyInt(),
		)
		debugscreen.PutString("\nPress START to exit.")
		display.WaitVblankStart()
	}

	kernel.ExitGame()
}
