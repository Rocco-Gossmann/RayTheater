package stage

type Loadable interface {
	Load(ctx Context)
}

type Unloadable interface {
	Unload(ctx Context) *interface{}
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
