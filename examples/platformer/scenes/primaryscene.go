package scenes

import (
	"github.com/TheBitDrifter/blueprint"
	blueprintclient "github.com/TheBitDrifter/blueprint/client"
	blueprintinput "github.com/TheBitDrifter/blueprint/input"
	blueprintmotion "github.com/TheBitDrifter/blueprint/motion"
	blueprintspatial "github.com/TheBitDrifter/blueprint/spatial"
	"github.com/TheBitDrifter/blueprint/vector"
	"github.com/TheBitDrifter/coldbrew/examples/platformer/animations"
	"github.com/TheBitDrifter/coldbrew/examples/platformer/components"
	"github.com/TheBitDrifter/coldbrew/examples/platformer/sounds"
	"github.com/TheBitDrifter/warehouse"
)

var PrimarySceneName = "PrimaryScene"

func PrimaryScene(height, width int, sto warehouse.Storage) error {
	playerArchetype, err := sto.NewOrExistingArchetype(
		blueprintspatial.Components.Position,
		blueprintclient.Components.SpriteBundle,
		blueprintspatial.Components.Direction,
		blueprintinput.Components.InputBuffer,
		blueprintclient.Components.CameraIndex,
		blueprintspatial.Components.Shape,
		blueprintmotion.Components.Dynamics,
		blueprintclient.Components.SoundBundle,
	)
	if err != nil {
		return err
	}

	err = playerArchetype.Generate(1,
		blueprintspatial.NewPosition(180, 340),
		blueprintspatial.NewRectangle(18, 58),
		blueprintmotion.NewDynamics(10),
		blueprintspatial.NewDirectionRight(),
		blueprintinput.InputBuffer{ReceiverIndex: 0},
		blueprintclient.CameraIndex(0),
		blueprintclient.NewSpriteBundle().
			AddSprite("characters/box_man_main.png", true).
			WithAnimations(animations.IdleAnimation, animations.RunAnimation, animations.FallAnimation, animations.JumpAnimation).
			SetActiveAnimation(animations.IdleAnimation).
			WithOffset(vector.Two{X: -72, Y: -59}).
			WithCustomRenderer(),
		blueprintclient.NewSoundBundle().
			AddSoundFromConfig(sounds.Run).
			AddSoundFromConfig(sounds.Jump).
			AddSoundFromConfig(sounds.Land),
	)
	if err != nil {
		return err
	}

	err = playerArchetype.Generate(1,
		blueprintspatial.NewPosition(540, 180),
		blueprintspatial.NewRectangle(18, 58),
		blueprintmotion.NewDynamics(10),
		blueprintspatial.NewDirectionRight(),
		blueprintinput.InputBuffer{ReceiverIndex: 1},
		blueprintclient.CameraIndex(1),
		blueprintclient.NewSpriteBundle().
			AddSprite("characters/box_man_alt.png", true).
			WithAnimations(animations.IdleAnimation, animations.RunAnimation, animations.FallAnimation, animations.JumpAnimation).
			SetActiveAnimation(animations.RunAnimation).
			WithOffset(vector.Two{X: -72, Y: -59}).
			WithCustomRenderer(),
		blueprintclient.NewSoundBundle().
			AddSoundFromConfig(sounds.Run).
			AddSoundFromConfig(sounds.Jump).
			AddSoundFromConfig(sounds.Land),
	)
	if err != nil {
		return err
	}
	// Create a custom parallax background
	// Use 0.0 if you want a still background (when we add camera movement later)
	err = blueprint.NewParallaxBackgroundBuilder(sto).
		AddLayer("backgrounds/city/sky.png", 0.05, 0.05).
		AddLayer("backgrounds/city/far.png", 0.2, 0.1).
		AddLayer("backgrounds/city/mid.png", 0.4, 0.2).
		AddLayer("backgrounds/city/near.png", 0.6, 0.6).
		Build()

	addTerrain(sto)

	musicArche, err := sto.NewOrExistingArchetype(blueprintclient.Components.SoundBundle, components.MusicTag)
	if err != nil {
		return err
	}
	musicArche.Generate(1, blueprintclient.NewSoundBundle().AddSoundFromPath("music.wav"))

	return nil
}

func addTerrain(sto warehouse.Storage) error {
	// Creating the new terrain archetype
	terrainArchetype, err := sto.NewOrExistingArchetype(
		components.BlockTerrainTag, // The new tag
		blueprintclient.Components.SpriteBundle,
		blueprintspatial.Components.Shape,
		blueprintspatial.Components.Position,
		blueprintmotion.Components.Dynamics,
	)
	if err != nil {
		return err
	}
	// Wall left (invisible)
	err = terrainArchetype.Generate(1,
		blueprintspatial.NewRectangle(10, 860),
		blueprintspatial.NewPosition(0, 0),
	)
	if err != nil {
		return err
	}
	// Wall right (invisible)
	err = terrainArchetype.Generate(1,
		blueprintspatial.NewRectangle(10, 860),
		blueprintspatial.NewPosition(1600, 0),
	)
	if err != nil {
		return err
	}
	// Floor
	err = terrainArchetype.Generate(1,
		blueprintspatial.NewPosition(1500, 470),
		blueprintspatial.NewRectangle(4000, 50),
		blueprintclient.NewSpriteBundle().
			AddSprite("terrain/floor.png", true).
			WithOffset(vector.Two{X: -1500, Y: -25}),
	)
	// Block
	err = terrainArchetype.Generate(1,
		blueprintspatial.NewPosition(285, 390),
		blueprintspatial.NewRectangle(64, 75),
		blueprintclient.NewSpriteBundle().
			AddSprite("terrain/block.png", true).
			WithOffset(vector.Two{X: -33, Y: -38}),
	)
	if err != nil {
		return err
	}

	// Ramp
	err = terrainArchetype.Generate(1,
		blueprintspatial.NewPosition(465, 422),
		blueprintspatial.NewDoubleRamp(250, 46, 0.2),
		blueprintclient.NewSpriteBundle().
			AddSprite("terrain/ramp.png", true).
			WithOffset(vector.Two{X: -125, Y: -22}),
	)
	if err != nil {
		return err
	}

	// Creating the new platform archetype
	platformArchetype, err := sto.NewOrExistingArchetype(
		components.PlatformTerrainTag,        // The new tag
		blueprintspatial.Components.Rotation, // We will use rotation to make some of these platforms sloped
		blueprintclient.Components.SpriteBundle,
		blueprintspatial.Components.Shape,
		blueprintspatial.Components.Position,
		blueprintmotion.Components.Dynamics,
	)
	if err != nil {
		return err
	}

	// Platforms left to right

	// Platform one
	err = platformArchetype.Generate(1,
		blueprintspatial.NewTriangularPlatform(144, 16),
		blueprintspatial.NewPosition(130, 350),
		blueprintclient.NewSpriteBundle().
			AddSprite("terrain/platform.png", true).
			WithOffset(vector.Two{X: -72, Y: -8}),
	)
	//  Platform 2
	err = platformArchetype.Generate(1,
		blueprintspatial.NewTriangularPlatform(144, 16),
		blueprintspatial.NewPosition(280, 240),
		blueprintclient.NewSpriteBundle().
			AddSprite("terrain/platform.png", true).
			WithOffset(vector.Two{X: -72, Y: -8}),
	)
	// Platform 3
	err = platformArchetype.Generate(1,
		blueprintspatial.NewTriangularPlatform(144, 16),
		blueprintspatial.NewPosition(430, 350),
		blueprintclient.NewSpriteBundle().
			AddSprite("terrain/platform.png", true).
			WithOffset(vector.Two{X: -72, Y: -8}),
	)

	// Platform 4
	err = platformArchetype.Generate(1,
		blueprintspatial.NewTriangularPlatform(144, 16),
		blueprintspatial.NewPosition(610, 350),
		blueprintspatial.Rotation(0.2),
		blueprintclient.NewSpriteBundle().
			AddSprite("terrain/platform.png", true).
			WithOffset(vector.Two{X: -72, Y: -8}),
	)
	// Platform 5
	err = platformArchetype.Generate(1,
		blueprintspatial.NewTriangularPlatform(144, 16),
		blueprintspatial.NewPosition(690, 250),
		blueprintspatial.Rotation(-0.3),
		blueprintclient.NewSpriteBundle().
			AddSprite("terrain/platform.png", true).
			WithOffset(vector.Two{X: -72, Y: -8}),
	)
	// Platform 6
	err = platformArchetype.Generate(1,
		blueprintspatial.NewTriangularPlatform(144, 16),
		blueprintspatial.NewPosition(870, 220),
		blueprintspatial.Rotation(0.6),
		blueprintclient.NewSpriteBundle().
			AddSprite("terrain/platform.png", true).
			WithOffset(vector.Two{X: -72, Y: -8}),
	)
	return nil
}
