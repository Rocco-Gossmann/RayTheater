package scenes

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/rocco-gossmann/RayTheater/stage"
)

type sceneMainStruct struct{}

var sceneMain = sceneMainStruct{}

var mainSceneLoad stage.GenericFunction = func(ctx stage.Context) {
	log.Println("scene loaded")
}

var mainSceneUnload stage.UnloadFunction = func(ctx stage.Context) *stage.Scene {
	log.Println("scene unloaded")
	return nil
}

var SceneMain = stage.NewScene(&sceneMain).
	OnLoad(&mainSceneLoad).
	OnUnload(&mainSceneUnload)

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
