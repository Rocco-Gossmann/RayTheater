package stage

import rl "github.com/gen2brain/raylib-go/raylib"

type sBuilder struct{ _valid bool }

func Build(stageWidth int32, stageHeight int32, scale float32) (s sBuilder) {

	if _setup {
		panic(" stage has been Build already ")
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
func (s sBuilder) Title(t string) sBuilder {
	s._validate()
	rl.SetWindowTitle(t)
	return s
}

// Defines with how many Frames per Second the sBuilder is sopposed to run
func (s sBuilder) FPS(f int32) sBuilder {
	s._validate()
	rl.SetTargetFPS(f)
	return s
}

func (s sBuilder) IntegerScale(enable bool) sBuilder {
	s._validate()
	_integerScale = enable
	return s
}

func (s sBuilder) Debug(level rl.TraceLogLevel) sBuilder {
	s._validate()
	rl.SetTraceLogLevel(level)
	return s
}

func (s sBuilder) Play(scene any) {
	stagePlay(scene)
}

func (s sBuilder) _validate() {
	if !s._valid {
		panic(" can't modify a stage that is not set up ")
	}

	if _playing {
		panic(" can't interact with a running stageplay ")
	}
}
