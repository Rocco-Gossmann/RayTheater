package stage

type Scene[T any] struct {
	core       T
	load       *GenericFunction
	tick       *TickFunction
	stageDraw  *GenericFunction
	windowDraw *GenericFunction
	unload     *UnloadFunction
}

func NewScene[T any](template T) *Scene[T] {

	sc := new(Scene[T])
	(*sc).core = template
	(*sc).windowDraw = &defaultGeneric
	(*sc).stageDraw = &defaultGeneric
	(*sc).load = &defaultGeneric
	(*sc).unload = &defaultUnload
	(*sc).tick = &defaultTick

	return sc
}

func (sc *Scene[T]) GetCore() T {
	return (*sc).core
}

func (sc *Scene[T]) OnLoad(fnc *GenericFunction) *Scene[T] {

	if fnc == nil {
		(*sc).load = &defaultGeneric
	} else {
		(*sc).load = fnc
	}

	return sc
}

func (sc *Scene[T]) OnUnload(fnc *UnloadFunction) *Scene[T] {

	if fnc == nil {
		(*sc).unload = &defaultUnload
	} else {
		(*sc).unload = fnc
	}

	return sc
}

func (sc *Scene[T]) OnStageDraw(fnc *GenericFunction) *Scene[T] {

	if fnc == nil {
		(*sc).stageDraw = &defaultGeneric
	} else {
		(*sc).stageDraw = fnc
	}

	return sc
}

func (sc *Scene[T]) OnWindowDraw(fnc *GenericFunction) *Scene[T] {

	if fnc == nil {
		(*sc).windowDraw = &defaultGeneric
	} else {
		(*sc).windowDraw = fnc
	}

	return sc
}

func (sc *Scene[T]) OnTick(fnc *TickFunction) *Scene[T] {
	if fnc == nil {
		(*sc).tick = &defaultTick
	} else {
		(*sc).tick = fnc
	}

	return sc
}
