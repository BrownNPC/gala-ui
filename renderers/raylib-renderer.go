package renderers

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RaylibRenderer struct{}

func (r RaylibRenderer) DrawRect(x, y, width, height int32, col color.RGBA) {
	rl.DrawRectangle(x, y, width, height, col)
}
