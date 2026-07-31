package main

import (
	"unsafe"

	"github.com/pspdev-go/pspsdk-go/psp/ctrl"
	"github.com/pspdev-go/pspsdk-go/psp/display"
	"github.com/pspdev-go/pspsdk-go/psp/gu"
	"github.com/pspdev-go/pspsdk-go/psp/gum"
	"github.com/pspdev-go/pspsdk-go/psp/kernel"
)

const displayListWords = 262144

type vertex struct {
	Color uint32
	X     float32
	Y     float32
	Z     float32
}

var displayList [displayListWords + 4]uint32

var triangle = [...]vertex{
	{0xff4f4fff, 0, 1, 0},
	{0xff4fff4f, -1, -1, 0},
	{0xffff4f4f, 1, -1, 0},
}

func alignedDisplayList() unsafe.Pointer {
	address := uintptr(unsafe.Pointer(&displayList[0]))
	return unsafe.Pointer((address + 15) &^ uintptr(15))
}

func initGraphics() {
	frameBuffer0 := gu.GuGetStaticVramBuffer(
		uint32(gu.VRAMWidth), uint32(gu.ScreenHeight), uint32(gu.PSM8888),
	)
	frameBuffer1 := gu.GuGetStaticVramBuffer(
		uint32(gu.VRAMWidth), uint32(gu.ScreenHeight), uint32(gu.PSM8888),
	)

	gu.SceGuInit()
	gu.SceGuStart(gu.Direct, alignedDisplayList())
	gu.SceGuDrawBuffer(gu.PSM8888, frameBuffer0, gu.VRAMWidth)
	gu.SceGuDispBuffer(gu.ScreenWidth, gu.ScreenHeight, frameBuffer1, gu.VRAMWidth)
	gu.SceGuOffset(uint32(2048-gu.ScreenWidth/2), uint32(2048-gu.ScreenHeight/2))
	gu.SceGuViewport(2048, 2048, gu.ScreenWidth, gu.ScreenHeight)
	gu.SceGuScissor(0, 0, gu.ScreenWidth, gu.ScreenHeight)
	gu.SceGuEnable(gu.ScissorTest)
	gu.SceGuEnable(gu.ClipPlanes)
	gu.SceGuFinish()
	gu.SceGuSync(gu.SyncFinish, gu.SyncWhatDone)

	display.WaitVblankStart()
	gu.SceGuDisplay(gu.DisplayOn)
}

func drawTriangle(angle float32) {
	gu.SceGuStart(gu.Direct, alignedDisplayList())
	gu.SceGuClearColor(0xff201810)
	gu.SceGuClear(gu.ColorBufferBit)

	draw := gum.Draw3D{
		Prim:        gu.Triangles,
		VType:       gu.Color8888 | gu.Vertex32BitF | gu.Transform3D,
		Count:       int32(len(triangle)),
		Vertices:    unsafe.Pointer(&triangle[0]),
		FOVY:        60,
		Aspect:      float32(gu.ScreenWidth) / float32(gu.ScreenHeight),
		Near:        0.5,
		Far:         100,
		Translation: gum.Vector3{Z: -3},
		Rotation:    gum.Vector3{Z: angle},
	}
	gum.DrawArray3D(&draw)

	gu.SceGuFinish()
	gu.SceGuSync(gu.SyncFinish, gu.SyncWhatDone)
	display.WaitVblankStart()
	gu.SceGuSwapBuffers()
}

func main() {
	ctrl.SetSamplingCycle(0)
	ctrl.SetSamplingMode(ctrl.ModeDigital)
	initGraphics()

	var pad ctrl.SceCtrlData
	var angle float32
	for {
		ctrl.PeekBufferPositive(&pad, 1)
		if pad.Buttons&ctrl.ButtonStart != 0 {
			break
		}
		drawTriangle(angle)
		angle += 0.02
	}

	gu.SceGuTerm()
	kernel.ExitGame()
}
