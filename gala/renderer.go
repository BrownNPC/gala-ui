package gala

import "image/color"

type Renderer interface {
	DrawRect(poxX, posY, width, height int32, color color.RGBA)
	MousePos() (x int32, y int32)
}
