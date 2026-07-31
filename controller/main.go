package main

import (
	"github.com/pspdev-go/pspsdk-go/psp/ctrl"
	"github.com/pspdev-go/pspsdk-go/psp/debugscreen"
	"github.com/pspdev-go/pspsdk-go/psp/display"
	"github.com/pspdev-go/pspsdk-go/psp/kernel"
)

var buttons = [...]struct {
	mask ctrl.Buttons
	name string
}{
	{ctrl.ButtonUp, "UP"},
	{ctrl.ButtonDown, "DOWN"},
	{ctrl.ButtonLeft, "LEFT"},
	{ctrl.ButtonRight, "RIGHT"},
	{ctrl.ButtonTriangle, "TRIANGLE"},
	{ctrl.ButtonCircle, "CIRCLE"},
	{ctrl.ButtonCross, "CROSS"},
	{ctrl.ButtonSquare, "SQUARE"},
	{ctrl.ButtonLTrigger, "L"},
	{ctrl.ButtonRTrigger, "R"},
	{ctrl.ButtonSelect, "SELECT"},
}

func main() {
	debugscreen.Init()
	ctrl.SetSamplingCycle(0)
	ctrl.SetSamplingMode(ctrl.ModeAnalog)

	var pad ctrl.SceCtrlData
	for {
		ctrl.PeekBufferPositive(&pad, 1)
		if pad.Buttons&ctrl.ButtonStart != 0 {
			break
		}

		debugscreen.SetXY(0, 0)
		debugscreen.PutString("Controller Example\n\n")
		debugscreen.Printf("Analog stick: X=%3d Y=%3d\n\n", pad.Lx, pad.Ly)
		debugscreen.PutString("Pressed: ")

		pressed := false
		for _, button := range buttons {
			if pad.Buttons&button.mask != 0 {
				if pressed {
					debugscreen.PutString(", ")
				}
				debugscreen.PutString(button.name)
				pressed = true
			}
		}
		if !pressed {
			debugscreen.PutString("none")
		}
		debugscreen.PutString("                              \n\nPress START to exit.")
		display.WaitVblankStart()
	}

	kernel.ExitGame()
}
