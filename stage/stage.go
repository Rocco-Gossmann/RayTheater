package stage

import (
	"log"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Required to make sure Engine can quit by itself
type exitSceneStruct struct{}

func (s exitSceneStruct) Tick(ctx Context) bool           { return false }
func (s exitSceneStruct) StageDraw(ctx Context)           {}
func (s exitSceneStruct) WindowDraw(ctx Context)          {}
func (s exitSceneStruct) Unload(ctx Context) *interface{} { return nil }

var exitScene exitSceneStruct

var (
	defaultStageDraw  StageDrawable
	defaultWindowDraw WindowDrawable
	defaultUnload     Unloadable
	defaultTick       Tickable
)

func init() {
	log.SetPrefix("[Stage.init] ")
	var ok bool = false

	log.Println("exit scene: ", exitScene)

	defaultStageDraw, ok = interface{}(exitScene).(StageDrawable)
	log.Println("init StageDraw:", ok, defaultStageDraw)

	defaultWindowDraw, ok = interface{}(exitScene).(WindowDrawable)
	log.Println("init WindowDraw:", ok, defaultWindowDraw)

	defaultTick, ok = interface{}(exitScene).(Tickable)
	log.Println("init Tick:", ok, defaultTick)

	defaultUnload, ok = interface{}(exitScene).(Unloadable)
	log.Println("init Unload:", ok, defaultUnload)
}

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

	_scene          any             = nil
	_tickable       *Tickable       = &defaultTick
	_stagedrawable  *StageDrawable  = &defaultStageDraw
	_windowdrawable *WindowDrawable = &defaultWindowDraw
	_unloadable     *Unloadable     = &defaultUnload

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
func (s Stage) Play(sc any) {
	s._validate()

	// Just so a window opens, at all, before we do anything else.
	rl.BeginDrawing()
	rl.ClearBackground(rl.Red)
	rl.EndDrawing()

	// disable ESC to Quit
	// TODO: enable
	//	rl.SetExitKey(0)

	switchScene(sc)

	_playing = true

MainLoop:
	for _scene != nil {
		log.SetPrefix("[Stage.Play] ")
		if rl.WindowShouldClose() {
			log.Println("Should close")
			if _scene != &exitScene {
				switchScene(&exitScene)
			}
			break MainLoop
		}

		if rl.IsWindowResized() {
			onWindowResize(&_viewportRect, &_viewportScale)
		}

		if !(*_tickable).Tick(ctx) {
			log.Println("tick == false => end")
			switchScene(nil)
			continue MainLoop
		}

		tickle := _tickableStack.Next
		for tickle != nil {
			if (*tickle.Item).Tick(ctx) {
				tickle = (*tickle).Next
			} else {
				dead := tickle
				tickle = (*tickle).Next
				dead.Drop()
			}
		}

		// Drawing on Gamebuffer (Resolution/Window-Size independed)
		//=====================================================================
		rl.BeginTextureMode(_stage)
		rl.ClearBackground(_stageBackground)
		(*_stagedrawable).StageDraw(ctx)
		rl.EndTextureMode()

		// Drawing on Screen (Resolution/Window-Size depended)
		//=====================================================================
		rl.BeginDrawing()
		rl.ClearBackground(_viewportBackground)
		rl.DrawTexturePro(_stage.Texture, _stageRect, _viewportRect, _stageOrigin, 0, rl.White)
		(*_windowdrawable).WindowDraw(ctx)
		rl.EndDrawing()

		// Debug outputs
		//=====================================================================
		//		if debug {
		//			rl.DrawFPS(4, 4)
		//		}
	}
}

func switchScene(scene any) {
	log.SetPrefix("[Stage.switchScene]")
	log.Printf("Called %v %p", scene, scene)
	if scene == nil {
		log.Println("scene is nil => attempt to load scene from Unlaod")
		scene = (*_unloadable).Unload(ctx)
		log.Println("Got Scene:", scene)

	} else {
		log.Println("scene is not nil => attempting regular unload")
		(*_unloadable).Unload(ctx)
		log.Println("Unload done")
	}

	_scene = nil

	if scene != nil {
		log.Println("new scene is not nil: ", scene)
		_scene = scene

		if i, ok := interface{}(_scene).(Tickable); ok {
			log.Println("Scene is Tickable:", i)
			_tickable = &i
		} else {
			log.Println("Scene not Tickable:")
			_tickable = &defaultTick
		}

		if i, ok := interface{}(_scene).(Unloadable); ok {
			log.Println("Scene is Unloadable:", i)
			_unloadable = &i
		} else {
			log.Println("Scene not Unloadable:")
			_unloadable = &defaultUnload
		}

		if i, ok := interface{}(_scene).(WindowDrawable); ok {
			log.Println("Scene is WindowDrawable:", i)
			_windowdrawable = &i
		} else {
			log.Println("Scene not WindowDrawable:")
			_windowdrawable = &defaultWindowDraw
		}

		if i, ok := interface{}(scene).(StageDrawable); ok {
			log.Println("Scene is StageDrawable:", i)
			_stagedrawable = &i
		} else {
			log.Println("Scene not StageDrawable:")
			_stagedrawable = &defaultStageDraw
		}

		if l, ok := interface{}(_scene).(Loadable); ok {
			log.Println("Scene is Loadable:", l)
			l.Load(ctx)
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
