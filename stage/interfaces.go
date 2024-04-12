package stage

type GenericFunction func(ctx Context)
type UnloadFunction func(ctx Context) *Scene[interface{}]
type TickFunction func(ctx Context) bool

type Loadable interface {
	Load(ctx Context)
}

type Unloadable interface {
	Unload(ctx Context) *Scene[interface{}]
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

// Default Scene
// ==============================================================================
type defaultSceneStruct struct{}

var defaultTick TickFunction = func(_ Context) bool { return false }
var defaultGeneric GenericFunction = func(_ Context) {}
var defaultUnload UnloadFunction = func(_ Context) *Scene[interface{}] { return nil }
