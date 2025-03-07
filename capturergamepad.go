package coldbrew

import (
	"log/slog"
	"math"

	"github.com/TheBitDrifter/bark"
	blueprintinput "github.com/TheBitDrifter/blueprint/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var _ InputCapturer = &gamepadCapturer{}

// gamepadCapturer handles gamepad input detection and processing
type gamepadCapturer struct {
	client  *client
	logger  *slog.Logger
	buttons map[ebiten.GamepadID][]ebiten.GamepadButton
	axes    map[ebiten.GamepadID][]float64
	sticks  map[ebiten.GamepadID]stickState
	active  bool
}

// stick represents an analog stick's position
type stick struct {
	X float64
	Y float64
}

// stickState tracks analog stick positions for a gamepad
type stickState struct {
	Left  stick
	Right stick
}

// newGamepadCapturer creates a new gamepad input capturer
func newGamepadCapturer(client *client) *gamepadCapturer {
	return &gamepadCapturer{
		client:  client,
		logger:  bark.For("gamepad"),
		buttons: make(map[ebiten.GamepadID][]ebiten.GamepadButton),
		axes:    make(map[ebiten.GamepadID][]float64),
		sticks:  make(map[ebiten.GamepadID]stickState),
	}
}

// Capture processes all gamepad inputs
func (h *gamepadCapturer) Capture() {
	h.handleDisconnections()
	h.handleNewConnections()
	h.captureInputs()
	h.processReceiverInputs()
}

// handleDisconnections cleans up state for disconnected gamepads
func (h *gamepadCapturer) handleDisconnections() {
	for id := range h.buttons {
		if inpututil.IsGamepadJustDisconnected(id) {
			h.logger.Info("gamepad disconnected",
				"id", id,
				"name", ebiten.GamepadName(id),
				"sdl_id", ebiten.GamepadSDLID(id),
			)
			delete(h.buttons, id)
			delete(h.axes, id)
			delete(h.sticks, id)
		}
	}
}

// handleNewConnections initializes state for newly connected gamepads
func (h *gamepadCapturer) handleNewConnections() {
	newIds := inpututil.AppendJustConnectedGamepadIDs([]ebiten.GamepadID{})
	for _, id := range newIds {
		h.buttons[id] = []ebiten.GamepadButton{}
		h.axes[id] = []float64{}
		h.sticks[id] = stickState{}
		h.logger.Info("gamepad connected",
			"id", id,
			"name", ebiten.GamepadName(id),
			"sdl_id", ebiten.GamepadSDLID(id),
		)
	}
}

// captureInputs reads the current state of all connected gamepads
func (h *gamepadCapturer) captureInputs() {
	for id := range h.buttons {
		h.captureButtonState(id)
		h.captureAxesState(id)
	}
}

// captureButtonState records which buttons are currently pressed
func (h *gamepadCapturer) captureButtonState(id ebiten.GamepadID) {
	var pressedButtons []ebiten.GamepadButton
	pressedButtons = inpututil.AppendPressedGamepadButtons(id, pressedButtons)
	h.buttons[id] = pressedButtons
}

// captureAxesState captures analog input state, preferring standard gamepad layout
// when available, falling back to raw axis values
func (h *gamepadCapturer) captureAxesState(id ebiten.GamepadID) {
	if ebiten.IsStandardGamepadLayoutAvailable(id) {
		h.sticks[id] = stickState{
			Left: processStickInput(
				ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal),
				ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical),
			),
			Right: processStickInput(
				ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisRightStickHorizontal),
				ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisRightStickVertical),
			),
		}
		return
	}

	// Fallback to raw axes if standard layout isn't available
	maxAxis := ebiten.GamepadAxisCount(id)
	axisValues := make([]float64, maxAxis)

	for a := 0; a < maxAxis; a++ {
		axisValues[a] = ebiten.GamepadAxisValue(id, ebiten.GamepadAxisType(a))
	}
	h.axes[id] = axisValues

	if len(axisValues) >= 4 {
		h.sticks[id] = stickState{
			Left:  processStickInput(axisValues[0], axisValues[1]),
			Right: processStickInput(axisValues[2], axisValues[3]),
		}
	} else {
		h.logger.Debug("insufficient axes for stick mapping",
			"gamepad_id", id,
			"axes_count", len(axisValues),
		)
	}
}

// processReceiverInputs forwards gamepad inputs to active receivers
func (h *gamepadCapturer) processReceiverInputs() {
	x, y := ebiten.CursorPosition()

	for i, receiver := range h.client.receivers {
		if !receiver.active {
			continue
		}

		h.processButtonInputs(i, receiver, x, y)
		h.processStickInputs(i, receiver, x, y)
	}
}

// processButtonInputs handles button presses for a specific receiver
func (h *gamepadCapturer) processButtonInputs(receiverIndex int, receiver *receiver, x, y int) {
	buttons := h.buttons[ebiten.GamepadID(receiver.padID)]
	for _, btn := range buttons {
		if !receiver.padLayout.mask.Contains(uint32(btn)) {
			continue
		}
		input := receiver.padLayout.buttons[btn]
		h.client.receivers[receiverIndex].inputs.pad = append(
			h.client.receivers[receiverIndex].inputs.pad,
			blueprintinput.StampedInput{
				Tick: tick,
				X:    x,
				Y:    y,
				Val:  input,
			},
		)
	}
}

// processStickInputs handles analog stick movement for a specific receiver
func (h *gamepadCapturer) processStickInputs(receiverIndex int, receiver *receiver, x, y int) {
	stickState, ok := h.sticks[ebiten.GamepadID(receiver.padID)]
	if !ok {
		return
	}

	// Process left stick if enabled
	if (stickState.Left.X != 0 || stickState.Left.Y != 0) && h.client.receivers[receiverIndex].leftAxes {
		h.client.receivers[receiverIndex].inputs.pad = append(
			h.client.receivers[receiverIndex].inputs.pad,
			blueprintinput.StampedInput{
				Tick: tick,
				// Use stick X/Y values directly as the vector components
				X:   int(stickState.Left.X * 100),  // Scale to a reasonable range if needed
				Y:   int(-stickState.Left.Y * 100), // Invert Y and scale to a reasonable range
				Val: h.client.receivers[receiverIndex].leftAxesInput,
			},
		)
		h.logger.Debug("gamepad left stick processed",
			"x", stickState.Left.X,
			"y", -stickState.Left.Y, // Showing inverted Y value in logs
			"val", h.client.receivers[receiverIndex].leftAxesInput,
		)
	}

	// Process right stick if enabled
	if (stickState.Right.X != 0 || stickState.Right.Y != 0) && h.client.receivers[receiverIndex].rightAxes {
		h.client.receivers[receiverIndex].inputs.pad = append(
			h.client.receivers[receiverIndex].inputs.pad,
			blueprintinput.StampedInput{
				Tick: tick,
				// Use stick X/Y values directly as the vector components
				X:   int(stickState.Right.X * 100),  // Scale to a reasonable range if needed
				Y:   int(-stickState.Right.Y * 100), // Invert Y and scale to a reasonable range
				Val: h.client.receivers[receiverIndex].rightAxesInput,
			},
		)
		h.logger.Debug("gamepad right stick processed",
			"x", stickState.Right.X,
			"y", -stickState.Right.Y, // Showing inverted Y value in logs
			"val", h.client.receivers[receiverIndex].rightAxesInput,
		)
	}
}

// processStickInput processes raw stick input with deadzone and normalization
// Returns normalized stick position after deadzone processing
func processStickInput(x, y float64) stick {
	const deadzone = 0.25
	x = math.Max(math.Min(x, 1.0), -1.0)
	y = math.Max(math.Min(y, 1.0), -1.0)

	magnitude := math.Sqrt(x*x + y*y)
	if magnitude < deadzone || magnitude == 0 {
		return stick{X: 0, Y: 0}
	}

	normalizedX := x / magnitude
	normalizedY := y / magnitude

	adjustedMagnitude := (magnitude - deadzone) / (1 - deadzone)
	adjustedMagnitude = math.Max(0, math.Min(1, adjustedMagnitude))
	return stick{
		X: normalizedX * adjustedMagnitude,
		Y: normalizedY * adjustedMagnitude,
	}
}
