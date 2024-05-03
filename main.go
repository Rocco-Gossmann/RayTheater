package main

import (
	"fmt"

	"github.com/rocco-gossmann/RayTheater/scenes"
	"github.com/rocco-gossmann/RayTheater/stage"
)

// ==============================================================================
// Setup the Stage
// ==============================================================================
func main() {

	fmt.Println("hello world")

	var ms scenes.SceneMain

	stage.Build(256, 192, 3.0).
		Title("RayTheater - DemoProject").
		IntegerScale(true).
		FPS(60).
		Play(&ms)

}
