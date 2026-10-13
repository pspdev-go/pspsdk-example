package main

import (
	"image/color"
	"reflect"
	"unsafe"

	"github.com/pspdev-go/pspsdk-go/psp/ctrl"
	"github.com/pspdev-go/pspsdk-go/psp/display"
	"github.com/pspdev-go/pspsdk-go/psp/gu"
	"github.com/pspdev-go/pspsdk-go/psp/loadexec"
	"github.com/pspdev-go/pspsdk-go/psp/threadman"
	"github.com/sago35/koebiten"
)

const (
	screenWidth  = 480
	screenHeight = 272
	bufferWidth  = 512
)

type pspHardware struct {
	display          pspDisplay
	keyState         [koebiten.KeyMax + 1]bool
	previousKeyState [koebiten.KeyMax + 1]bool
	exitRequested    bool
}

var exitRequestedByPSP bool

func (h *pspHardware) Init() error {
	if display.SetMode(display.ModeLCD, screenWidth, screenHeight) != display.ErrorOK {
		return errPSPDisplay
	}
	registerExitCallback()
	ctrl.SetSamplingCycle(0)
	ctrl.SetSamplingMode(ctrl.ModeAnalog)
	return h.display.Init()
}

func (h *pspHardware) GetDisplay() koebiten.Displayer {
	return &h.display
}

func (h *pspHardware) KeyUpdate() error {
	threadman.SceKernelCheckCallback()
	var input ctrl.SceCtrlData
	ctrl.PeekBufferPositive(&input, 1)
	h.exitRequested = exitRequestedByPSP
	for _, key := range pspKeys {
		h.previousKeyState[key.key] = h.keyState[key.key]
		h.keyState[key.key] = input.Buttons&key.button != 0
	}
	return nil
}

func registerExitCallback() bool {
	callbackID := threadman.SceKernelCreateCallback(
		unsafe.Pointer(&exitCallbackName[0]),
		unsafe.Pointer(reflect.ValueOf(pspExitCallback).Pointer()),
		nil,
	)
	if callbackID < 0 {
		return false
	}
	return loadexec.SceKernelRegisterExitCallback(callbackID) >= 0
}

func pspExitCallback(_, _ int32, _ unsafe.Pointer) int32 {
	exitRequestedByPSP = true
	return 0
}

var exitCallbackName = [...]byte{'E', 'x', 'i', 't', ' ', 'C', 'a', 'l', 'l', 'b', 'a', 'c', 'k', 0}

func (h *pspHardware) ExitRequested() bool {
	return h.exitRequested
}

func (h *pspHardware) IsKeyPressed(key koebiten.Key) bool {
	return key >= 0 && int(key) < len(h.keyState) && h.keyState[key]
}

func (h *pspHardware) IsKeyJustPressed(key koebiten.Key) bool {
	return key >= 0 && int(key) < len(h.keyState) && h.keyState[key] && !h.previousKeyState[key]
}

var errPSPDisplay = pspDisplayError{}

type pspDisplayError struct{}

func (pspDisplayError) Error() string { return "failed to initialize PSP display" }

type pspDisplay struct {
	framebuffers [2]*[bufferWidth * screenHeight]uint32
	drawIndex    int
}

type pspKey struct {
	button ctrl.Buttons
	key    koebiten.Key
}

var pspKeys = [...]pspKey{
	{ctrl.ButtonLeft, koebiten.KeyLeft},
	{ctrl.ButtonRight, koebiten.KeyRight},
	{ctrl.ButtonUp, koebiten.KeyUp},
	{ctrl.ButtonDown, koebiten.KeyDown},
	{ctrl.ButtonCross, koebiten.Key0},
	{ctrl.ButtonCircle, koebiten.Key1},
	{ctrl.ButtonSquare, koebiten.Key2},
	{ctrl.ButtonTriangle, koebiten.Key3},
}

func (d *pspDisplay) Init() error {
	d.framebuffers[0] = pspFrameBuffer()
	d.framebuffers[1] = pspFrameBuffer()
	if d.framebuffers[0] == nil || d.framebuffers[1] == nil {
		return errPSPDisplay
	}
	white := rgba8888(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	for _, framebuffer := range d.framebuffers {
		for i := range framebuffer {
			framebuffer[i] = white
		}
	}
	if display.SetFrameBuf(unsafe.Pointer(d.framebuffers[0]), bufferWidth, display.PixelFormat8888, display.SetBufNextVSync) != display.ErrorOK {
		return errPSPDisplay
	}
	d.drawIndex = 1
	display.WaitVblankStart()
	return nil
}

func pspFrameBuffer() *[bufferWidth * screenHeight]uint32 {
	address := uintptr(gu.GuGetStaticVramTexture(bufferWidth, screenHeight, uint32(gu.PSM8888)))
	return (*[bufferWidth * screenHeight]uint32)(unsafe.Pointer(address | 0x40000000))
}

func (d *pspDisplay) Size() (int16, int16) {
	return screenWidth, screenHeight
}

func (d *pspDisplay) SetPixel(x, y int16, c color.RGBA) {
	if x < 0 || x >= screenWidth || y < 0 || y >= screenHeight {
		return
	}
	d.framebuffers[d.drawIndex][int(y)*bufferWidth+int(x)] = rgba8888(c)
}

func (d *pspDisplay) Display() error {
	if display.SetFrameBuf(unsafe.Pointer(d.framebuffers[d.drawIndex]), bufferWidth, display.PixelFormat8888, display.SetBufNextVSync) != display.ErrorOK {
		return errPSPDisplay
	}
	display.WaitVblankStart()
	d.drawIndex = 1 - d.drawIndex
	return nil
}

func (d *pspDisplay) ClearDisplay() {
	d.ClearBuffer()
	_ = d.Display()
}

func (d *pspDisplay) ClearBuffer() {
	framebuffer := d.framebuffers[d.drawIndex]
	if framebuffer == nil {
		return
	}
	white := rgba8888(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	for i := range framebuffer {
		framebuffer[i] = white
	}
}

func rgba8888(c color.RGBA) uint32 {
	return uint32(c.R) | uint32(c.G)<<8 | uint32(c.B)<<16 | uint32(c.A)<<24
}
