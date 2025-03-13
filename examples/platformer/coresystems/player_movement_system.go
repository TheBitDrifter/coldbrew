package coresystems

import (
	"math"

	"github.com/TheBitDrifter/blueprint"
	blueprintinput "github.com/TheBitDrifter/blueprint/input"
	blueprintmotion "github.com/TheBitDrifter/blueprint/motion"
	blueprintspatial "github.com/TheBitDrifter/blueprint/spatial"
	"github.com/TheBitDrifter/blueprint/vector"
	"github.com/TheBitDrifter/coldbrew/examples/platformer/actions"
	"github.com/TheBitDrifter/coldbrew/examples/platformer/components"
	"github.com/TheBitDrifter/warehouse"
)

// PlayerMovementSystem handles all player movement logic including horizontal movement and jumping
type PlayerMovementSystem struct{}

// Run executes the movement system for each game tick
func (sys PlayerMovementSystem) Run(scene blueprint.Scene, dt float64) error {
	// Handle horizontal movement first, then jumping
	sys.handleHorizontal(scene)
	sys.handleJump(scene)
	return nil
}

// handleHorizontal processes left/right movement and slope interactions
func (PlayerMovementSystem) handleHorizontal(scene blueprint.Scene) {
	// Get all entities with input buffers
	cursor := scene.NewCursor(blueprint.Queries.InputBuffer)
	for range cursor.Next() {
		// --- Gather required components ---
		dyn := blueprintmotion.Components.Dynamics.GetFromCursor(cursor)              // Physics properties
		incomingInputs := blueprintinput.Components.InputBuffer.GetFromCursor(cursor) // User inputs
		direction := blueprintspatial.Components.Direction.GetFromCursor(cursor)      // Facing direction

		// Player's horizontal movement speed
		const speedX = 120.0

		// --- Process left/right inputs ---
		// Check and consume left input
		_, pressedLeft := incomingInputs.ConsumeInput(actions.Left)
		if pressedLeft {
			direction.SetLeft()
		}

		// Check and consume right input
		_, pressedRight := incomingInputs.ConsumeInput(actions.Right)
		if pressedRight {
			direction.SetRight()
		}

		// Track if player is moving horizontally this frame
		isMovingHorizontal := pressedLeft || pressedRight

		// --- Check if player is on ground ---
		playerIsGrounded, onGround := components.OnGroundComponent.GetFromCursorSafe(cursor)
		currentTick := scene.CurrentTick()

		// --- Ground snapping logic ---
		// If player is on ground, apply a small downward force to keep them attached to slopes
		if playerIsGrounded {
			ticksOnGround := currentTick - onGround.LastJump
			if ticksOnGround > 10 { // Only apply after being on ground for a while
				const snapForce = 40.0
				dyn.Vel.Y = math.Max(dyn.Vel.Y, snapForce) // Apply downward force
			}
		}

		// --- Handle in-air movement ---
		if !playerIsGrounded {
			// Air movement is simpler - just move in the direction pressed or stop
			if isMovingHorizontal {
				dyn.Vel.X = speedX * direction.AsFloat() // direction.AsFloat() returns -1 for left, 1 for right
			} else {
				dyn.Vel.X = 0 // No horizontal movement when no keys pressed
			}
			continue // Skip the rest of the loop for airborne players
		}

		// --- Handle flat ground movement ---
		// Check if the player is on a flat surface (normal pointing straight up)
		flat := onGround.SlopeNormal.X == 0 && onGround.SlopeNormal.Y == 1
		if flat {
			// Same as air movement on flat ground
			if isMovingHorizontal {
				dyn.Vel.X = speedX * direction.AsFloat()
			} else {
				dyn.Vel.X = 0
			}
			continue // Skip slope handling
		}

		// --- Handle slope movement ---
		if isMovingHorizontal && playerIsGrounded {
			// Calculate tangent vector along the slope
			// The tangent is perpendicular to the normal, so we swap X/Y and negate Y
			tangent := vector.Two{X: onGround.SlopeNormal.Y, Y: -onGround.SlopeNormal.X}

			// Determine if player is moving uphill
			// This is determined by checking if the direction and normal X have the same sign
			isUphill := (direction.AsFloat() * onGround.SlopeNormal.X) > 0

			// Scale tangent by movement direction for correct slope alignment
			slopeDir := tangent.Scale(direction.AsFloat())

			if isUphill {
				// When going uphill, only set X velocity and let physics handle Y
				dyn.Vel.X = slopeDir.X * speedX
			} else {
				// When going downhill, help player follow the slope
				dyn.Vel.X = slopeDir.X * speedX

				// Only apply after being on ground for a while
				ticksOnGround := currentTick - onGround.LastJump
				if ticksOnGround > 10 {
					// Apply downward velocity along the slope
					dyn.Vel.Y = slopeDir.Y * speedX
				}
			}
		} else {
			// No horizontal input while on a slope = stop moving
			dyn.Vel.X = 0
		}
	}
}

// handleJump processes jump input with support for coyote time and input buffering
func (PlayerMovementSystem) handleJump(scene blueprint.Scene) {
	// Create query for players eligible to jump (have ground and input components)
	playersEligibleToJumpQuery := warehouse.Factory.NewQuery()
	playersEligibleToJumpQuery.And(components.OnGroundComponent, blueprintinput.Components.InputBuffer)

	// Get all entities that match the query
	cursor := scene.NewCursor(playersEligibleToJumpQuery)
	currentTick := scene.CurrentTick()

	for range cursor.Next() {
		// Get required components
		dyn := blueprintmotion.Components.Dynamics.GetFromCursor(cursor)
		incomingInputs := blueprintinput.Components.InputBuffer.GetFromCursor(cursor)
		onGround := components.OnGroundComponent.GetFromCursor(cursor)

		// Jump strength
		const jumpForce = 350.0

		// Check for jump input (Up action)
		if stampedInput, inputReceived := incomingInputs.ConsumeInput(actions.Up); inputReceived {
			// --- Jump Eligibility Checks ---

			// Coyote time: Allow jumping within 8 ticks of leaving ground
			const coyoteTimeTicks = 8
			playerGroundedWithinCoyoteTime := currentTick-onGround.LastTouch <= coyoteTimeTicks

			// Input buffer: Allow jump inputs from shortly before landing to register
			const inputBufferTicks = 8
			jumpInputWithinBufferWindow := math.Abs(float64(stampedInput.Tick-onGround.LastTouch)) <= inputBufferTicks

			// Prevent double jumps: Make sure we haven't jumped since last touching ground
			playerHasNotJumpedSinceGroundTouch := onGround.LastJump < onGround.LastTouch

			// All conditions must be met to jump
			canJump := playerGroundedWithinCoyoteTime &&
				jumpInputWithinBufferWindow &&
				playerHasNotJumpedSinceGroundTouch

			if canJump {
				// Apply upward velocity and acceleration for jump
				dyn.Vel.Y = -jumpForce
				dyn.Accel.Y = -jumpForce
				// Record jump time
				onGround.LastJump = currentTick
			}
		}
	}
}
