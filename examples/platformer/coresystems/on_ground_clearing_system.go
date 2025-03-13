package coresystems

import (
	"github.com/TheBitDrifter/blueprint"
	"github.com/TheBitDrifter/coldbrew/examples/platformer/components"
	"github.com/TheBitDrifter/warehouse"
)

type OnGroundClearingSystem struct{}

func (OnGroundClearingSystem) Run(scene blueprint.Scene, dt float64) error {
	// Query any entity that has onGround
	onGroundQuery := warehouse.Factory.NewQuery().And(components.OnGroundComponent)
	onGroundCursor := scene.NewCursor(onGroundQuery)

	// Iterate through matched entities
	for range onGroundCursor.Next() {
		// Get the onGround component state
		onGround := components.OnGroundComponent.GetFromCursor(onGroundCursor)

		// 8 tick coyote time window
		if scene.CurrentTick()-onGround.LastTouch > 8 {
			groundedEntity, _ := onGroundCursor.CurrentEntity()

			// We can't mutate while iterating so we enqueue the changes instead
			groundedEntity.EnqueueRemoveComponent(components.OnGroundComponent)
		}
	}
	return nil
}
