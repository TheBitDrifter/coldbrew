package clientsystems

import (
	"math"

	blueprintclient "github.com/TheBitDrifter/blueprint/client"
	blueprintinput "github.com/TheBitDrifter/blueprint/input"
	blueprintmotion "github.com/TheBitDrifter/blueprint/motion"
	"github.com/TheBitDrifter/coldbrew"
	"github.com/TheBitDrifter/coldbrew/examples/platformer/animations"
	"github.com/TheBitDrifter/coldbrew/examples/platformer/components"
	"github.com/TheBitDrifter/coldbrew/examples/platformer/sounds"
	"github.com/TheBitDrifter/warehouse"
)

// PlayerSoundSystem handles all sound effects related to player movement
type PlayerSoundSystem struct{}

// Run executes the sound system for each game tick
func (sys PlayerSoundSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	// Create a query for players that can play sounds and are on the ground
	playersWithSoundsOnTheGround := warehouse.Factory.NewQuery()
	playersWithSoundsOnTheGround.And(
		blueprintclient.Components.SoundBundle,  // Has sounds
		blueprintclient.Components.SpriteBundle, // Has sprites (for animation sync)
		blueprintinput.Components.InputBuffer,   // Can receive input
		blueprintmotion.Components.Dynamics,     // Has physics properties
		components.OnGroundComponent,            // Is on the ground
	)

	// Get all entities that match the query
	cursor := scene.NewCursor(playersWithSoundsOnTheGround)

	// Iterate over all matching players
	for cursor.Next() {
		// --- Get required components ---
		soundBundle := blueprintclient.Components.SoundBundle.GetFromCursor(cursor)
		dyn := blueprintmotion.Components.Dynamics.GetFromCursor(cursor)
		onGround := components.OnGroundComponent.GetFromCursor(cursor)
		currentTick := scene.CurrentTick()

		// --- Get consistent player index ---
		// We use the camera index to ensure consistent sound players
		idx := int(*blueprintclient.Components.CameraIndex.GetFromCursor(cursor))

		// --- Handle Landing Sound ---
		// Play landing sound when player just landed this tick
		if onGround.Landed == currentTick {
			landingSound, _ := coldbrew.MaterializeSound(soundBundle, sounds.Land)
			player := landingSound.GetPlayer(idx)

			// If not already playing, play from start
			if !player.IsPlaying() {
				player.Rewind()
				player.Play()
			}
		}

		// --- Handle Jump Sound ---
		// Play jump sound when player has upward velocity and just jumped
		if dyn.Vel.Y < 5 && onGround.LastJump == currentTick {
			jumpSound, _ := coldbrew.MaterializeSound(soundBundle, sounds.Jump)
			player := jumpSound.GetPlayer(idx)

			// If not already playing, play from start
			if !player.IsPlaying() {
				player.Rewind()
				player.Play()
			}
		}

		// --- Skip run sounds if not moving horizontally or just touched ground ---
		// This prevents run sounds when player is standing still or just landed
		const minMovementSpeed = 10.0
		if math.Abs(dyn.Vel.X) <= minMovementSpeed && onGround.LastTouch == currentTick {
			continue
		}

		// --- Handle Run Sound ---
		// Sync run sounds with the walk animation
		spriteBundle := blueprintclient.Components.SpriteBundle.GetFromCursor(cursor)

		// In this demo we only use one slot in the spriteBundle so we index 0
		runAnimation, _ := spriteBundle.Blueprints[0].GetAnim(animations.RunAnimation)

		// Get client's current tick (might differ from scene tick in multiplayer)
		clientTick := cli.CurrentTick()
		walkStartTick := runAnimation.StartTick

		// Calculate which frame of the animation we're on
		ticksSinceStart := clientTick - walkStartTick
		framesSinceStart := ticksSinceStart / runAnimation.Speed
		currentFrame := framesSinceStart % runAnimation.FrameCount

		// --- Play footstep sounds on specific frames ---
		// Only play sounds on "contact frames" (0 and 4)
		// These are the frames where feet touch the ground in the animation
		if currentFrame == 0 || currentFrame == 4 {
			// Only play on the first tick of the frame to avoid rapid repetition
			if ticksSinceStart%runAnimation.Speed == 0 {
				runSound, _ := coldbrew.MaterializeSound(soundBundle, sounds.Run)

				// Use any available player to avoid interrupting previous steps
				// This allows overlapping footstep sounds for natural running
				player, _ := runSound.GetAnyAvailable()

				if !player.IsPlaying() {
					player.Rewind()
				}

				// Set slightly lower volume for footsteps
				player.SetVolume(0.8)
				player.Play()
			}
		}
	}

	return nil
}
