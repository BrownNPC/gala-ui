package main

import (
	"gala/gala"
	"gala/renderers"
	"image/color"
	"math/rand"
	"net/http"

	// _ "net/http/pprof" // Import the pprof package to register pprof handlers

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	renderer := renderers.RaylibRenderer{}
	// Start the pprof HTTP server in a separate goroutine
	go func() {
		http.ListenAndServe(":8080", nil) // Start the pprof server on port 6060
	}()

	var layout = gala.NewLayout(1280, 720, 199)
	rl.InitWindow(1280, 720, "yo")
	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.White)
		layout.Box().
			Id("Card").
			Size(gala.Percent(100), gala.Percent(100)).
			Display_Flex().
			AlignItems_Center().
			BackgroundColor(color.RGBA{26, 26, 29, 255}).
			FlexDirection_Row().
			Padding(10).
			Contains(
				layout.Box().
					Size(100, 100).
					Id("sigma").
					BackgroundColor(rl.Blue).
					Flex(1).
					Display_Flex(),

				layout.Box().
					Size(100, 100).
					Id("sigma").
					BackgroundColor(rl.Beige).
					Flex(1).
					Display_Flex(),
				layout.Box().
					Height(200).
					Left(20).
					Right(10).
					Bottom(100).
					// Position_Relative().
					Flex(1).
					BackgroundColor(rl.Pink).
					Hovered(func(box *gala.Box) {
						box.BackgroundColor(color.RGBA{0, uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255})
					}),
			)
		layout.End(renderer)

		rl.EndDrawing()
	}
}
