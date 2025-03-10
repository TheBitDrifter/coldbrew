package main

import (
	"embed"
	"log"

	"github.com/TheBitDrifter/blueprint"
	"github.com/TheBitDrifter/coldbrew"
	"github.com/TheBitDrifter/coldbrew/examples/platformer/actions"
	"github.com/TheBitDrifter/coldbrew/examples/platformer/clientsystems"
	"github.com/TheBitDrifter/coldbrew/examples/platformer/coresystems"
	"github.com/TheBitDrifter/coldbrew/examples/platformer/rendersystems"
	"github.com/TheBitDrifter/coldbrew/examples/platformer/scenes"
	"github.com/hajimehoshi/ebiten/v2"

	coldbrew_clientsystems "github.com/TheBitDrifter/coldbrew/clientsystems"
	coldbrew_rendersystems "github.com/TheBitDrifter/coldbrew/rendersystems"
	tteo_coresystems "github.com/TheBitDrifter/tteokbokki/coresystems"
)

//go:embed assets/*
var assets embed.FS

func main() {
	// Kinda weird resolution, but works for web demo
	resolutionX := 640
	resolutionY := 720
	maxSpritesCached := 200
	maxSoundsCached := 200
	maxScenesCached := 10

	// Create the client
	client := coldbrew.NewClient(
		resolutionX,
		resolutionY,
		maxSpritesCached,
		maxSoundsCached,
		maxScenesCached,
		assets,
	)

	// Configure the client window
	client.SetTitle("Bappa Platformer!")
	client.SetWindowSize(640, 720)
	client.SetResizable(true)

	// Loader settingd
	client.SetMinimumLoadTime(30)
	client.SetEnforceMinOnActive(true)

	// Common system slices
	renderSystems := []coldbrew.RenderSystem{
		rendersystems.PlayerCameraPriorityRenderer{},
	}

	clientSystems := []coldbrew.ClientSystem{
		coldbrew_clientsystems.BackgroundScrollSystem{},
		clientsystems.CameraFollowerSystem{},
		clientsystems.PlayerAnimationSystem{},
		clientsystems.MusicSystem{},
		clientsystems.PlayerSoundSystem{},
	}

	coreSystems := []blueprint.CoreSystem{
		coresystems.GravitySystem{},                 // Apply gravity forces
		coresystems.PlayerMovementSystem{},          // Apply player input forces
		tteo_coresystems.IntegrationSystem{},        // Update velocities and positions
		tteo_coresystems.TransformSystem{},          // Update collision shapes
		coresystems.PlayerBlockCollisionSystem{},    // Handle solid block collisions
		coresystems.PlayerPlatformCollisionSystem{}, // Handle one-way platforms
		coresystems.OnGroundClearingSystem{},        // Update grounded state/coyote time
	}

	// Register Scenes
	primarySceneWidth, primarySceneHeight := 1600, 520
	err := client.RegisterScene(
		scenes.PrimarySceneName,
		primarySceneWidth,
		primarySceneHeight,
		scenes.PrimaryScene,
		renderSystems,
		clientSystems,
		coreSystems,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = client.RegisterScene(
		scenes.SecondarySceneName,
		primarySceneWidth,
		primarySceneHeight,
		scenes.SecondaryScene,
		renderSystems,
		clientSystems,
		coreSystems,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Register Global Systems
	client.RegisterGlobalRenderSystem(coldbrew_rendersystems.GlobalRenderer{}, &coldbrew_rendersystems.DebugRenderer{})
	client.RegisterGlobalClientSystem(
		coldbrew_clientsystems.InputBufferSystem{},
		&coldbrew_clientsystems.CameraSceneAssignerSystem{},
		clientsystems.PlayerTransferSystem{},
		clientsystems.FullScreenSystem{},
	)

	// Activate cameras
	cameraOne, err := client.ActivateCamera()
	if err != nil {
		log.Fatal(err)
	}
	cameraTwo, err := client.ActivateCamera()
	if err != nil {
		log.Fatal(err)
	}

	// Position Cameras
	halfScreenTall := resolutionY / 2

	cameraOne.SetDimensions(resolutionX, halfScreenTall)
	cameraTwo.SetDimensions(resolutionX, halfScreenTall)

	// Get camera two position reference and position it appropriately
	cameraTwoPosition, _ := cameraTwo.Positions()
	cameraTwoPosition.Y = float64(halfScreenTall)

	client.SetCameraBorderSize(8)

	// Register Receivers
	registerReceivers(client)
	// Run the client
	if err := client.Start(); err != nil {
		log.Fatal(err)
	}
}

func registerReceivers(client coldbrew.Client) {
	receiver1, _ := client.ActivateReceiver()
	receiver1.RegisterKey(ebiten.KeyW, actions.Up)
	receiver1.RegisterKey(ebiten.KeyS, actions.Down)
	receiver1.RegisterKey(ebiten.KeyA, actions.Left)
	receiver1.RegisterKey(ebiten.KeyD, actions.Right)

	// Key 1 (number key 1) to transfer player 1
	receiver1.RegisterKey(ebiten.Key1, actions.PlayerTransfer)

	receiver1.RegisterPad(0)
	receiver1.RegisterGamepadAxes(false, actions.StickMovement)

	// Receiver 2 is arrow keys.
	receiver2, _ := client.ActivateReceiver()
	receiver2.RegisterKey(ebiten.KeyUp, actions.Up)
	receiver2.RegisterKey(ebiten.KeyDown, actions.Down)
	receiver2.RegisterKey(ebiten.KeyLeft, actions.Left)
	receiver2.RegisterKey(ebiten.KeyRight, actions.Right)

	// Key 2 (number key 2) to transfer player 2
	receiver2.RegisterKey(ebiten.Key2, actions.PlayerTransfer)
}
