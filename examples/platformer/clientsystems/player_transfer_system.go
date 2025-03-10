package clientsystems

import (
	"fmt"
	"log"

	"github.com/TheBitDrifter/blueprint"
	blueprint_input "github.com/TheBitDrifter/blueprint/input"
	"github.com/TheBitDrifter/coldbrew"
	"github.com/TheBitDrifter/coldbrew/examples/platformer/actions"
	"github.com/TheBitDrifter/coldbrew/examples/platformer/scenes"
	"github.com/TheBitDrifter/warehouse"
)

// playerTransfer holds information needed to transfer a player between scenes
type playerTransfer struct {
	origin       coldbrew.Scene   // The scene the player is currently in
	target       coldbrew.Scene   // The scene the player will be transferred to
	playerEntity warehouse.Entity // The player entity to transfer
}

// PlayerTransferSystem handles moving player entities between different game scenes
type PlayerTransferSystem struct{}

// Run executes the player transfer system for each game tick
func (PlayerTransferSystem) Run(cli coldbrew.Client) error {
	// Store pending transfers to process after scanning all scenes
	var pending []playerTransfer

	// Get scene cache for looking up available scenes
	sceneCache := cli.Cache()

	// --- Step 1: Detect player transfer requests ---
	// Iterate through all active scenes
	for _, activeScene := range cli.ActiveScenes() {
		// Skip scenes that aren't fully loaded yet
		if !activeScene.Ready() {
			continue
		}

		// Get all entities with input buffers in this scene
		cursor := activeScene.NewCursor(blueprint.Queries.InputBuffer)
		for cursor.Next() {
			// Check for player transfer input
			inputBuffer := blueprint_input.Components.InputBuffer.GetFromCursor(cursor)
			_, transferRequested := inputBuffer.ConsumeInput(actions.PlayerTransfer)

			if transferRequested {
				// --- Apply cooldown to prevent rapid transfers ---
				const cooldownTicks = 60 // ~1 second at 60 FPS

				// Check if player is still in cooldown period
				inCooldown := (cli.CurrentTick() - activeScene.LastSelectedTick()) <= cooldownTicks
				if inCooldown {
					continue
				}

				// Get the player entity that requested the transfer
				currentPlayerEntity, err := cursor.CurrentEntity()
				if err != nil {
					return err
				}

				// --- Determine target scene ---
				// Simple toggle between scenes
				var sceneTargetName string
				if activeScene.Name() == scenes.PrimarySceneName {
					sceneTargetName = scenes.SecondarySceneName
				} else {
					sceneTargetName = scenes.PrimarySceneName
				}

				// Look up the target scene in cache
				targetSceneIndex, found := sceneCache.GetIndex(sceneTargetName)
				if !found {
					log.Println("Target scene not found:", sceneTargetName)
					return fmt.Errorf("target scene '%s' not found in cache", sceneTargetName)
				}

				// Get the actual scene object
				targetScene := sceneCache.GetItem(targetSceneIndex)

				// Create and store the transfer request
				transfer := playerTransfer{
					origin:       activeScene,
					target:       targetScene,
					playerEntity: currentPlayerEntity,
				}
				pending = append(pending, transfer)
			}
		}
	}

	// --- Step 2: Process all player transfers ---
	// We do this after collecting all requests to avoid modifying
	// the scenes while we're iterating through them
	for _, transfer := range pending {
		cli.ActivateScene(transfer.target, transfer.playerEntity)
	}

	// --- Step 3: Clean up empty scenes ---
	// Deactivate any scenes that no longer have player entities
	for _, activeScene := range cli.ActiveScenes() {
		cursor := activeScene.NewCursor(blueprint.Queries.InputBuffer)
		if cursor.TotalMatched() == 0 {
			// No player entities with input buffers left in this scene
			cli.DeactivateScene(activeScene)
		}
	}

	return nil
}
