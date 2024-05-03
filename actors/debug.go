package actors

import "github.com/rocco-gossmann/RayTheater/stage/components"

type TheDebut struct {
	trns components.Transform2D
}

func (d TheDebut) GetTransform() *components.Transform2D {
	return &(d.trns)
}
