package scenes

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/rocco-gossmann/RayTheater/stage"
)

type SceneMain struct{}

func (s SceneMain) Tick(ctx stage.Context) bool {

	if rl.IsKeyDown(rl.KeyLeftAlt) && rl.IsKeyPressed(rl.KeyF4) {
		log.Println("Hit ALT+F4")
		return false
	}

	if rl.IsKeyReleased(rl.KeyEscape) {
		log.Println("Hit ESC")
		return false
	}

	return true
}

func (s SceneMain) StageDraw(ctx stage.Context) {
	rl.DrawText("Press and relase the ESC key\nor hit ALT+F4 to quit! ", 8, 8, 8, rl.LightGray)
}

// Debug overlay
func (s SceneMain) WindowDraw(ctx stage.Context) {
	rl.DrawRectangle(0, 0, 90, 30, rl.NewColor(0, 0, 0, 128))
	rl.DrawFPS(4, 4)
}
