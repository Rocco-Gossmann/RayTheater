package stage

/** used to pass information around through the engines live-cycle */
type Context struct {
	deltatime float32
}

var stageContext Context
