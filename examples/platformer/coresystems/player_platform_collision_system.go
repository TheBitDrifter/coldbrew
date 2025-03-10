package coresystems

import (
	"github.com/TheBitDrifter/blueprint"
	blueprintmotion "github.com/TheBitDrifter/blueprint/motion"
	blueprintspatial "github.com/TheBitDrifter/blueprint/spatial"
	"github.com/TheBitDrifter/coldbrew/examples/platformer/components"
	"github.com/TheBitDrifter/tteokbokki/motion"
	"github.com/TheBitDrifter/tteokbokki/spatial"
	"github.com/TheBitDrifter/warehouse"
)

type PlayerPlatformCollisionSystem struct{}

func (s PlayerPlatformCollisionSystem) Run(scene blueprint.Scene, dt float64) error {
	// Create cursors
	platformTerrainQuery := warehouse.Factory.NewQuery().And(components.PlatformTerrainTag)
	platformCursor := scene.NewCursor(platformTerrainQuery)
	playerCursor := scene.NewCursor(blueprint.Queries.InputBuffer)

	// Outer loop is blocks
	for platformCursor.Next() {
		// Inner is players
		for playerCursor.Next() {
			// Delegate to helper
			err := s.resolve(scene, platformCursor, playerCursor)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (PlayerPlatformCollisionSystem) resolve(scene blueprint.Scene, platformCursor, playerCursor *warehouse.Cursor) error {
	playerShape := blueprintspatial.Components.Shape.GetFromCursor(playerCursor)
	playerPosition := blueprintspatial.Components.Position.GetFromCursor(playerCursor)
	playerDynamics := blueprintmotion.Components.Dynamics.GetFromCursor(playerCursor)

	platformShape := blueprintspatial.Components.Shape.GetFromCursor(platformCursor)
	platformPosition := blueprintspatial.Components.Position.GetFromCursor(platformCursor)
	platformDynamics := blueprintmotion.Components.Dynamics.GetFromCursor(platformCursor)

	if ok, collisionResult := spatial.Detector.Check(
		*playerShape, *platformShape, playerPosition.Two, platformPosition.Two,
	); ok {
		// Apply the collision resolution if conditions are met
		// We don't want to land on a platform unless falling
		// We don't want to collide with the sides
		if playerDynamics.Vel.Y > 0 && collisionResult.IsTopB() {
			// Use a vertical resolver since we cant collide with the sides
			motion.VerticalResolver.Resolve(
				&playerPosition.Two,
				&platformPosition.Two,
				playerDynamics,
				platformDynamics,
				collisionResult,
			)
			currentTick := scene.CurrentTick()
			playerAlreadyGrounded, onGround := components.OnGroundComponent.GetFromCursorSafe(playerCursor)

			// If not grounded, enqueue onGround with values
			if !playerAlreadyGrounded {
				playerEntity, _ := playerCursor.CurrentEntity()
				err := playerEntity.EnqueueAddComponentWithValue(
					components.OnGroundComponent,
					components.OnGround{LastTouch: currentTick, Landed: currentTick, SlopeNormal: collisionResult.Normal},
				)
				if err != nil {
					return err
				}
				// Otherwise update the existing OnGround
			} else {
				onGround.LastTouch = scene.CurrentTick()
				onGround.SlopeNormal = collisionResult.Normal
			}

		}
	}
	return nil
}
