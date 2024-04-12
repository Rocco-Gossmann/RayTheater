package stage

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	_stageWidth         float32 = 640
	_stageHeight        float32 = 360
	_stage              rl.RenderTexture2D
	_stageRect          rl.Rectangle
	_stageOrigin        rl.Vector2
	_stageBackground    rl.Color = rl.NewColor(24, 96, 128, 128)
	_stageUnload        Unloadable
	_sceneTicklistEntry *LinkedList[*Tickable]

	// Defines if the viewport was created
	_viewportScale      float32 = 1.0
	_viewportRect       rl.Rectangle
	_viewportBackground rl.Color = rl.NewColor(64, 192, 255, 255)

	_setup        bool = false
	_playing      bool = false
	_integerScale bool = false

	_tickableStack = NewLinkedList[*Tickable]()

	_scene *Scene = nil

	ctx = Context{}
)

type Stage struct{ _valid bool }

// Must be called first, to create your stage
func Setup(stageWidth int32, stageHeight int32, scale float32) (s Stage) {
	if _setup {
		panic(" stage has been setup already ")
	}

	_setup = true

	_stageWidth = float32(stageWidth)
	_stageHeight = float32(stageHeight)
	_stageRect.Width = _stageWidth
	_stageRect.Height = -_stageHeight

	_viewportRect.Width = float32(stageWidth) * scale
	_viewportRect.Height = float32(stageHeight) * scale

	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(int32(_viewportRect.Width), int32(_viewportRect.Height), "< RayTheater - Project >")

	if !rl.IsWindowReady() {
		panic(" failed to open window ")
	}

	_stage = rl.LoadRenderTexture(stageWidth, stageHeight)
	if !rl.IsRenderTextureReady(_stage) {
		panic(" failed to create state rendertexture ")
	}

	s._valid = true

	return
}

// Changes the Stages Title
func (s Stage) Title(t string) Stage {
	s._validate()
	rl.SetWindowTitle(t)
	return s
}

// Defines with how many Frames per Second the Stage is sopposed to run
func (s Stage) FPS(f int32) Stage {
	s._validate()
	rl.SetTargetFPS(f)
	return s
}

func (s Stage) IntegerScale(enable bool) Stage {
	s._validate()
	_integerScale = enable
	return s
}

func (s Stage) Debug(level rl.TraceLogLevel) Stage {
	s._validate()
	rl.SetTraceLogLevel(level)
	return s
}

/** Starts to play the given Scene on the Stage
* @param scene - a struct that implements various interfaces  (@see ./stage/scene.go)
 */
func (s Stage) Play(sc *Scene) {
	s._validate()

	rl.BeginDrawing()
	rl.ClearBackground(rl.Red)
	rl.EndDrawing()

	rl.SetExitKey(0)

	switchScene(sc)

	_playing = true

MainLoop:
	for _playing && _scene != nil {

		if rl.WindowShouldClose() {
			_playing = false
			break
		}

		if rl.IsWindowResized() {
			onWindowResize(&_viewportRect, &_viewportScale)
		}

		tickle := _tickableStack.Next
		for tickle != nil {
			if (*tickle.Item).Tick(ctx) {
				tickle = (*tickle).Next
			} else {
				dead := tickle
				tickle = (*tickle).Next
				dead.Drop()

				// if the Dead handler belongs to the current scene
				if dead == _sceneTicklistEntry {
					// Switch the scene
					switchScene(nil)
					// and skip the rendering for this frame
					// That way, we don't need an extra check for the
					// _scene becoming nil (aka. the Loop ending)
					continue MainLoop
				}
			}
		}

		// Drawing on Gamebuffer (Resolution/Window-Size independed)
		//=====================================================================
		rl.BeginTextureMode(_stage)
		rl.ClearBackground(_stageBackground)
		(*_scene.StageDraw).StageDraw(ctx)
		rl.EndTextureMode()

		// Drawing on Screen (Resolution/Window-Size depended)
		//=====================================================================
		rl.BeginDrawing()
		rl.ClearBackground(_viewportBackground)
		rl.DrawTexturePro(_stage.Texture, _stageRect, _viewportRect, _stageOrigin, 0, rl.White)
		(*_scene.WindowDraw).WindowDraw(ctx)
		rl.EndDrawing()

		// Debug outputs
		//=====================================================================
		//		if debug {
		//			rl.DrawFPS(4, 4)
		//		}
	}

	// Unload Scene if needed
	if _scene != nil && (*_scene).Unload != nil {
		(*_scene.Unload).Unload(ctx)
	}

}

func switchScene(scene *Scene) {

	if _scene != nil {
		if (*_scene).Unload != nil {
			if scene == nil {
				scene = (*_scene.Unload).Unload(ctx)
			} else {
				(*_scene.Unload).Unload(ctx)
			}
		}

		if _sceneTicklistEntry != nil {
			_sceneTicklistEntry.Drop()
			_sceneTicklistEntry = nil
		}

		_scene = nil
	}

	if scene != nil {
		_scene = scene

		if (*_scene).Load != nil {
			(*_scene.Load).Load(ctx)
		}

		if scene.Tick != nil {
			_sceneTicklistEntry = _tickableStack.Prepend(scene.Tick)
		}

	}
}

func onWindowResize(vp *rl.Rectangle, scale *float32) {

	wW, wH := float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())

	*scale = min(wW/float32(_stageWidth), wH/float32(_stageHeight))

	// TODO: Reimplement Integer-Scaling
	//
	if _integerScale && *scale > 1 {
		*scale = float32(math.Floor(float64(*scale)))
	}

	vp.Width = float32(_stageWidth) * *scale
	vp.Height = float32(_stageHeight) * *scale
	vp.X = (wW - vp.Width) * 0.5
	vp.Y = (wH - vp.Height) * 0.5
}

func (s Stage) _validate() {
	if !s._valid {
		panic(" can't modify a stage that is not set up ")
	}

	if _playing {
		panic(" can't interact with a running stageplay ")
	}
}
