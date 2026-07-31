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

var cube = [...]vertex{
	// Front.
	{0xffe85d75, -1, -1, 1}, {0xffe85d75, -1, 1, 1}, {0xffe85d75, 1, 1, 1},
	{0xffe85d75, -1, -1, 1}, {0xffe85d75, 1, 1, 1}, {0xffe85d75, 1, -1, 1},
	// Back.
	{0xff5d8ee8, -1, -1, -1}, {0xff5d8ee8, 1, -1, -1}, {0xff5d8ee8, 1, 1, -1},
	{0xff5d8ee8, -1, -1, -1}, {0xff5d8ee8, 1, 1, -1}, {0xff5d8ee8, -1, 1, -1},
	// Right.
	{0xff62d49b, 1, -1, -1}, {0xff62d49b, 1, -1, 1}, {0xff62d49b, 1, 1, 1},
	{0xff62d49b, 1, -1, -1}, {0xff62d49b, 1, 1, 1}, {0xff62d49b, 1, 1, -1},
	// Left.
	{0xffd4b362, -1, -1, -1}, {0xffd4b362, -1, 1, -1}, {0xffd4b362, -1, 1, 1},
	{0xffd4b362, -1, -1, -1}, {0xffd4b362, -1, 1, 1}, {0xffd4b362, -1, -1, 1},
	// Top.
	{0xffb06be8, -1, 1, -1}, {0xffb06be8, 1, 1, -1}, {0xffb06be8, 1, 1, 1},
	{0xffb06be8, -1, 1, -1}, {0xffb06be8, 1, 1, 1}, {0xffb06be8, -1, 1, 1},
	// Bottom.
	{0xff56cfe1, -1, -1, -1}, {0xff56cfe1, -1, -1, 1}, {0xff56cfe1, 1, -1, 1},
	{0xff56cfe1, -1, -1, -1}, {0xff56cfe1, 1, -1, 1}, {0xff56cfe1, 1, -1, -1},
}

func alignedDisplayList() unsafe.Pointer {
	base := unsafe.Pointer(&displayList[0])
	offset := -uintptr(base) & 15
	return unsafe.Add(base, offset)
}

func initGraphics() {
	frameBuffer0 := gu.GuGetStaticVramBuffer(
		uint32(gu.VRAMWidth), uint32(gu.ScreenHeight), uint32(gu.PSM8888),
	)
	frameBuffer1 := gu.GuGetStaticVramBuffer(
		uint32(gu.VRAMWidth), uint32(gu.ScreenHeight), uint32(gu.PSM8888),
	)
	depthBuffer := gu.GuGetStaticVramBuffer(
		uint32(gu.VRAMWidth), uint32(gu.ScreenHeight), uint32(gu.PSM4444),
	)

	gu.SceGuInit()
	gu.SceGuStart(gu.Direct, alignedDisplayList())
	gu.SceGuDrawBuffer(gu.PSM8888, frameBuffer0, gu.VRAMWidth)
	gu.SceGuDispBuffer(gu.ScreenWidth, gu.ScreenHeight, frameBuffer1, gu.VRAMWidth)
	gu.SceGuDepthBuffer(depthBuffer, gu.VRAMWidth)
	gu.SceGuOffset(
		uint32(2048-gu.ScreenWidth/2),
		uint32(2048-gu.ScreenHeight/2),
	)
	gu.SceGuViewport(2048, 2048, gu.ScreenWidth, gu.ScreenHeight)
	gu.SceGuDepthRange(65535, 0)
	gu.SceGuScissor(0, 0, gu.ScreenWidth, gu.ScreenHeight)
	gu.SceGuEnable(gu.ScissorTest)
	gu.SceGuDepthFunc(gu.GEqual)
	gu.SceGuEnable(gu.DepthTest)
	gu.SceGuFrontFace(gu.CW)
	gu.SceGuShadeModel(gu.Smooth)
	gu.SceGuEnable(gu.CullFace)
	gu.SceGuEnable(gu.ClipPlanes)
	gu.SceGuFinish()
	gu.SceGuSync(gu.SyncFinish, gu.SyncWhatDone)

	display.WaitVblankStart()
	gu.SceGuDisplay(gu.DisplayOn)
}

func drawCube(angle float32) {
	gu.SceGuStart(gu.Direct, alignedDisplayList())
	gu.SceGuClearColor(0xff181221)
	gu.SceGuClearDepth(0)
	gu.SceGuClear(gu.ColorBufferBit | gu.DepthBufferBit)

	gu.SceGuAmbientColor(0xffffffff)
	draw := gum.Draw3D{
		Prim:        gu.Triangles,
		VType:       gu.Color8888 | gu.Vertex32BitF | gu.Transform3D,
		Count:       int32(len(cube)),
		Vertices:    unsafe.Pointer(&cube[0]),
		FOVY:        60,
		Aspect:      float32(gu.ScreenWidth) / float32(gu.ScreenHeight),
		Near:        0.5,
		Far:         100,
		Translation: gum.Vector3{Z: -3.5},
		Rotation:    gum.Vector3{X: angle * 0.7, Y: angle, Z: angle * 0.35},
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
		ctrl.ReadBufferPositive(&pad, 1)
		if pad.Buttons&ctrl.ButtonStart != 0 {
			break
		}
		drawCube(angle)
		angle += 0.025
	}

	gu.SceGuTerm()
	kernel.ExitGame()
}
