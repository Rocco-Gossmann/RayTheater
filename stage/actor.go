package stage

type Actor[T any] struct {
	Core   *T
	Load   *Loadable
	Tick   *Tickable
	Unload *Unloadable
	dead   bool
}

func NewActor[T any](template *T) (act Actor[T]) {
	act.Core = template

	if l, ok := interface{}(template).(*Loadable); ok {
		act.Load = l
	}

	if l, ok := interface{}(template).(*Tickable); ok {
		act.Tick = l
	}

	if l, ok := interface{}(template).(*Unloadable); ok {
		act.Unload = l
	}

	return
}
