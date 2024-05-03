package components

import rl "github.com/gen2brain/raylib-go/raylib"

type Transform2D struct {

	// Where Actor is this frame
	loc rl.Vector2

	// Where it will be next frame
	_loc rl.Vector2

	_scRegistered *func()
}

func (t Transform2D) RegisterAtStage() *func() {

	if t._scRegistered != nil {
		panic("can't register Transform2D-Component at 2 different stages at once")
	}

	f := func() {
		t.loc.X = t._loc.X
		t.loc.Y = t._loc.Y
	}

	t._scRegistered = &f

	return t._scRegistered

}

func (t Transform2D) FreeFromStage(f *func()) {
	if f != t._scRegistered {
		panic("can't free a Transform2D-Component you did not register")
	}

	t._scRegistered = nil
}

func (t Transform2D) GetX() float32 {
	return t.loc.X
}

func (t Transform2D) GetY() float32 {
	return t.loc.Y
}

func (t Transform2D) GetLoc() rl.Vector2 {
	return t.loc
}

func (t Transform2D) SetLoc(x, y float32) {
	t._loc.X = x
	t._loc.Y = y
}
