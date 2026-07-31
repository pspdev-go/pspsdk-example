package main

import (
	_ "embed"

	"github.com/pspdev-go/pspsdk-go/psp/ctrl"
	"github.com/pspdev-go/pspsdk-go/psp/debugscreen"
)

//go:embed message.txt
var message string

func main() {
	debugscreen.Init()
	ctrl.SetSamplingCycle(0)
	ctrl.SetSamplingMode(ctrl.ModeDigital)

	for {
		debugscreen.SetXY(0, 0)
		debugscreen.Printf("embed: %s", message)
	}
}
