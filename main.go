package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/rocco-gossmann/RayTheater/scenes"
	"github.com/rocco-gossmann/RayTheater/stage"
)

// ==============================================================================
// Setup the Stage
// ==============================================================================
func main() {

	var ms scenes.SceneMain

	stage.Build(256, 192, 3.0).
		Title("RayTheater - DemoProject").
		IntegerScale(true).
		FPS(60).
		Debug(rl.LogDebug).
		Play(&ms)

}
