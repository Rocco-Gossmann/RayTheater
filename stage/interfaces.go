package stage

import "github.com/rocco-gossmann/RayTheater/stage/components"

type Loadable interface {
	Load(ctx Context)
}

type Unloadable interface {
	Unload(ctx Context) any
}

type Tickable interface {
	Tick(ctx Context) bool
}

type StageDrawable interface {
	StageDraw(ctx Context)
}

type WindowDrawable interface {
	WindowDraw(ctx Context)
}

type ITransform2D interface {
	getTransform() *components.Transform2D
}

type iFreeable interface {
	FreeFromStage(f *func())
}
