package main

import (
	"unsafe"

	"github.com/pspdev-go/pspsdk-go/psp/ctrl"
	"github.com/pspdev-go/pspsdk-go/psp/debugscreen"
	"github.com/pspdev-go/pspsdk-go/psp/display"
	"github.com/pspdev-go/pspsdk-go/psp/kernel"
	"github.com/pspdev-go/pspsdk-go/psp/rtc"
)

// dateTime mirrors the 16-byte pspTime structure used by the RTC API.
type dateTime struct {
	Year        uint16
	Month       uint16
	Day         uint16
	Hour        uint16
	Minute      uint16
	Second      uint16
	Microsecond uint32
}

var weekdays = [...]string{
	"Sunday", "Monday", "Tuesday", "Wednesday",
	"Thursday", "Friday", "Saturday",
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

		var now dateTime
		result := rtc.SceRtcGetCurrentClockLocalTime(unsafe.Pointer(&now))

		debugscreen.SetXY(0, 0)
		debugscreen.PutString("RTC Clock Example\n\n")
		if result < 0 {
			debugscreen.Printf("RTC error: %d", result)
		} else {
			weekday := rtc.SceRtcGetDayOfWeek(
				int32(now.Year), int32(now.Month), int32(now.Day),
			)
			weekdayName := "Unknown"
			if weekday >= 0 && weekday < int32(len(weekdays)) {
				weekdayName = weekdays[weekday]
			}
			debugscreen.Printf(
				"%04d-%02d-%02d  %s        \n",
				now.Year, now.Month, now.Day, weekdayName,
			)
			debugscreen.Printf(
				"%02d:%02d:%02d.%06d",
				now.Hour, now.Minute, now.Second, now.Microsecond,
			)
		}
		debugscreen.PutString("\n\nPress START to exit.")
		display.WaitVblankStart()
	}

	kernel.ExitGame()
}
