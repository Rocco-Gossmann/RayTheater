package stage

import (
	"log"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/rocco-gossmann/RayTheater/stage/components"
)

// Required to make sure Engine can quit by itself
type exitSceneStruct struct{}

func (s exitSceneStruct) Tick(ctx Context) bool  { return false }
func (s exitSceneStruct) StageDraw(ctx Context)  {}
func (s exitSceneStruct) WindowDraw(ctx Context) {}
func (s exitSceneStruct) Unload(ctx Context) any { return nil }

var exitScene exitSceneStruct

type ActorID uint

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

	_tickable = &defaultTick
	_stagedrawable = &defaultStageDraw
	_windowdrawable = &defaultWindowDraw
	_unloadable = &defaultUnload
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

	_scene          any = nil
	_tickable       *Tickable
	_stagedrawable  *StageDrawable
	_windowdrawable *WindowDrawable
	_unloadable     *Unloadable

	_actorCnt          uint
	_actorsTickable    map[ActorID]*Tickable
	_actorsTransform2D map[ActorID]*components.Transform2D
	_actorsUpdates     map[ActorID][]*func()
	_actorsFreeable    map[ActorID]map[*iFreeable]*func()

	ctx = Context{}
)

/** Starts to play the given Scene on the Stage
* @param scene - a struct that implements various interfaces  (@see ./stage/scene.go)
 */
func stagePlay(sc any) {

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

		for _, at := range _actorsUpdates {
			for _, fnc := range at {
				(*fnc)()
			}
		}

		if !(*_tickable).Tick(ctx) {
			log.Println("tick == false => end")
			switchScene(nil)
			continue MainLoop
		}

		for _, tickle := range _actorsTickable {
			(*tickle).Tick(ctx)
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

func AddActor(a any) (id ActorID) {
	log.SetPrefix("[Stage AddActor] ")
	_actorCnt++
	id = ActorID(_actorCnt)

	_actorsUpdates[id] = make([]*func(), 0, 3)
	_actorsFreeable[id] = make(map[*iFreeable]*func())

	if t, ok := interface{}(a).(components.Transform2D); ok {
		fnc := t.RegisterAtStage()
		_actorsUpdates[id] = append(_actorsUpdates[id], fnc)

		if f, ok := interface{}(t).(iFreeable); ok {
			_actorsFreeable[id][&f] = fnc
		}
		_actorsTransform2D[id] = &t
	}

	if t, ok := interface{}(a).(Tickable); ok {
		_actorsTickable[id] = &t
	}

	log.Printf(`
	Add Actor:      %v
		ID:			%d
		Tickable:   %v
		Transform:  %v
		Freeable:   %v
	\n`, a, id, _actorsTickable[id], _actorsTransform2D[id], _actorsFreeable[id])

	return id
}

func RemoveActor(a ActorID) {

	// Free all components bound to this stage
	for i, f := range _actorsFreeable[a] {
		(*i).FreeFromStage(f)
	}

	// Remove Actor from all Lists
	delete(_actorsFreeable, a)
	delete(_actorsTransform2D, a)
	delete(_actorsTickable, a)
}

func switchScene(scene any) {
	log.SetPrefix("[Stage.switchScene]")
	log.Printf("Called %v %p", scene, scene)
	if scene == nil {
		log.Println("scene is nil => attempt to load scene from Unload")
		scene = (*_unloadable).Unload(ctx)
		log.Println("Got Scene:", scene)

	} else {
		log.Printf("scene is not nil => attempting regular unload: %v %p %v\n", scene, _unloadable, *_unloadable)
		(*_unloadable).Unload(ctx)
		log.Println("Unload done")
	}

	_scene = nil

	if scene == nil || scene == 0x0 {
		log.Printf("new scene is nil\n")

	} else {
		log.Printf("new scene is not nil: %v (%p)\n", scene, scene)
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

	if _integerScale && *scale > 1 {
		*scale = float32(math.Floor(float64(*scale)))
	}

	vp.Width = float32(_stageWidth) * *scale
	vp.Height = float32(_stageHeight) * *scale
	vp.X = (wW - vp.Width) * 0.5
	vp.Y = (wH - vp.Height) * 0.5
}
