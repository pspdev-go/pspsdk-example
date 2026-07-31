package main

import (
	"unsafe"

	"github.com/pspdev-go/pspsdk-go/psp/ctrl"
	"github.com/pspdev-go/pspsdk-go/psp/debugscreen"
	"github.com/pspdev-go/pspsdk-go/psp/display"
	"github.com/pspdev-go/pspsdk-go/psp/iofilemgr"
	"github.com/pspdev-go/pspsdk-go/psp/kernel"
)

const (
	openReadOnly  = 0x0001
	openWriteOnly = 0x0002
	openCreate    = 0x0200
	openTruncate  = 0x0400
)

func cString(value string) []byte {
	return append([]byte(value), 0)
}

func writeAndReadFile() (string, int32) {
	path := cString("example.txt")
	message := []byte("Hello from the PSP file system!")

	fd := iofilemgr.SceIoOpen(
		unsafe.Pointer(&path[0]),
		openWriteOnly|openCreate|openTruncate,
		0o777,
	)
	if fd < 0 {
		return "", fd
	}
	written := iofilemgr.SceIoWrite(fd, unsafe.Pointer(&message[0]), uint32(len(message)))
	iofilemgr.SceIoClose(fd)
	if written < 0 {
		return "", written
	}

	fd = iofilemgr.SceIoOpen(unsafe.Pointer(&path[0]), openReadOnly, 0)
	if fd < 0 {
		return "", fd
	}
	var buffer [64]byte
	read := iofilemgr.SceIoRead(fd, unsafe.Pointer(&buffer[0]), uint32(len(buffer)))
	iofilemgr.SceIoClose(fd)
	if read < 0 {
		return "", read
	}
	return string(buffer[:read]), 0
}

func main() {
	debugscreen.Init()
	debugscreen.PutString("File I/O Example\n\n")

	text, result := writeAndReadFile()
	if result < 0 {
		debugscreen.Printf("File operation failed: %d", result)
	} else {
		debugscreen.PutString("Created: example.txt\n")
		debugscreen.PutString("Read back:\n")
		debugscreen.Printf("  %s", text)
	}
	debugscreen.PutString("\n\nPress START to exit.")

	ctrl.SetSamplingCycle(0)
	ctrl.SetSamplingMode(ctrl.ModeDigital)
	var pad ctrl.SceCtrlData
	for pad.Buttons&ctrl.ButtonStart == 0 {
		ctrl.PeekBufferPositive(&pad, 1)
		display.WaitVblankStart()
	}

	kernel.ExitGame()
}
