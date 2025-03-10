package clientsystems

import (
	"math"

	"github.com/TheBitDrifter/blueprint"
	blueprintclient "github.com/TheBitDrifter/blueprint/client"
	blueprintmotion "github.com/TheBitDrifter/blueprint/motion"
	"github.com/TheBitDrifter/coldbrew"
	"github.com/TheBitDrifter/coldbrew/examples/platformer/animations"
	"github.com/TheBitDrifter/coldbrew/examples/platformer/components"
)

type PlayerAnimationSystem struct{}

func (PlayerAnimationSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	cursor := scene.NewCursor(blueprint.Queries.InputBuffer)
	for cursor.Next() {
		bundle := blueprintclient.Components.SpriteBundle.GetFromCursor(cursor)
		spriteBlueprint := &bundle.Blueprints[0]
		dyn := blueprintmotion.Components.Dynamics.GetFromCursor(cursor)
		grounded, onGround := components.OnGroundComponent.GetFromCursorSafe(cursor)
		if grounded {
			grounded = scene.CurrentTick() == onGround.LastTouch
		}

		// Normal animation state transitions
		if math.Abs(dyn.Vel.X) > 0 && grounded {
			spriteBlueprint.TryAnimation(animations.RunAnimation)
		} else if dyn.Vel.Y > 0 && !grounded {
			spriteBlueprint.TryAnimation(animations.FallAnimation)
		} else if dyn.Vel.Y <= 0 && !grounded {
			spriteBlueprint.TryAnimation(animations.JumpAnimation)
		} else {
			spriteBlueprint.TryAnimation(animations.IdleAnimation)
		}
	}
	return nil
}
