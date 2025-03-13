package coresystems

import (
	"github.com/TheBitDrifter/blueprint"
	blueprintmotion "github.com/TheBitDrifter/blueprint/motion"
	"github.com/TheBitDrifter/tteokbokki/motion"
)

const (
	DEFAULT_GRAVITY  = 9.8
	PIXELS_PER_METER = 50.0
)

type GravitySystem struct{}

func (GravitySystem) Run(scene blueprint.Scene, dt float64) error {
	cursor := scene.NewCursor(blueprint.Queries.Dynamics)
	for range cursor.Next() {
		dyn := blueprintmotion.Components.Dynamics.GetFromCursor(cursor)
		mass := 1 / dyn.InverseMass
		gravity := motion.Forces.Generator.NewGravityForce(mass, DEFAULT_GRAVITY, PIXELS_PER_METER)
		motion.Forces.AddForce(dyn, gravity)
	}
	return nil
}
