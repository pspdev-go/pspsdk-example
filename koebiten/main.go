package main

import (
	"github.com/pspdev-go/pspsdk-go/psp/loadexec"
	"github.com/sago35/koebiten"
	"tinygo.org/x/drivers/pixel"
)

type Game struct {
	x          int
	y          int
	colorIndex int
	hardware   *pspHardware
}

func (g *Game) Update() error {
	if g.hardware.IsKeyPressed(koebiten.KeyLeft) {
		g.x -= 4
	}
	if g.hardware.IsKeyPressed(koebiten.KeyRight) {
		g.x += 4
	}
	if g.hardware.IsKeyPressed(koebiten.KeyUp) {
		g.y -= 4
	}
	if g.hardware.IsKeyPressed(koebiten.KeyDown) {
		g.y += 4
	}
	if g.hardware.IsKeyJustPressed(koebiten.Key0) {
		g.colorIndex = 0
	}
	if g.hardware.IsKeyJustPressed(koebiten.Key1) {
		g.colorIndex = 1
	}
	if g.hardware.IsKeyJustPressed(koebiten.Key2) {
		g.colorIndex = 2
	}
	if g.hardware.IsKeyJustPressed(koebiten.Key3) {
		g.colorIndex = 3
	}
	if g.x < 16 {
		g.x = 16
	}
	if g.x > 464 {
		g.x = 464
	}
	if g.y < 16 {
		g.y = 16
	}
	if g.y > 216 {
		g.y = 216
	}
	return nil
}

func (g *Game) Draw(screen *koebiten.Image) {
	koebiten.DrawFilledRect(nil, g.x-12, g.y-12, 24, 24, squareColors[g.colorIndex])
	koebiten.DrawFilledRect(nil, 8, 238, screenWidth-16, 1, pixel.NewRGB888(0, 0, 0))
	koebiten.DrawText(nil, "D-PAD: MOVE", nil, 8, 248, pixel.NewRGB888(0, 0, 0))
	koebiten.DrawText(nil, "X=BLUE  O=RED  SQUARE=GREEN  TRIANGLE=YELLOW", nil, 8, 257, pixel.NewRGB888(0, 0, 0))
	koebiten.DrawText(nil, "PS: QUIT MENU", nil, 8, 266, pixel.NewRGB888(0, 0, 0))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	hardware := &pspHardware{}
	if err := koebiten.SetHardware(hardware); err != nil {
		return
	}
	koebiten.SetWindowSize(screenWidth, screenHeight)

	game := &Game{x: screenWidth / 2, y: screenHeight / 2, hardware: hardware}
	display := hardware.GetDisplay()
	for {
		if err := hardware.KeyUpdate(); err != nil {
			return
		}
		if hardware.ExitRequested() {
			loadexec.SceKernelExitGame()
			return
		}
		display.ClearBuffer()
		if err := game.Update(); err != nil {
			return
		}
		game.Draw(nil)
		if err := display.Display(); err != nil {
			return
		}
	}
}

var squareColors = [...]pixel.RGB888{
	pixel.NewRGB888(0, 160, 255),
	pixel.NewRGB888(255, 64, 64),
	pixel.NewRGB888(0, 200, 64),
	pixel.NewRGB888(255, 196, 0),
}
