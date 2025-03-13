package rendersystems

import (
	"image"
	"log/slog"
	"math"

	"github.com/TheBitDrifter/bark"
	"github.com/TheBitDrifter/blueprint"
	blueprintclient "github.com/TheBitDrifter/blueprint/client"
	blueprintspatial "github.com/TheBitDrifter/blueprint/spatial"
	"github.com/TheBitDrifter/blueprint/vector"
	"github.com/TheBitDrifter/coldbrew"
	"github.com/TheBitDrifter/warehouse"
	"github.com/hajimehoshi/ebiten/v2"
)

// GlobalRenderer is the default render system for the coldbrew package
// It automatically handles sprites, sprite sheets, and parallax backgrounds
type GlobalRenderer struct {
	logger *slog.Logger
	sorted [][]RenderItem
}

// RenderItem contains all state needed for rendering an entity
type RenderItem struct {
	sprite    coldbrew.Sprite
	blueprint *blueprintclient.SpriteBlueprint
	pos       blueprintspatial.Position
	rot       blueprintspatial.Rotation
	scale     blueprintspatial.Scale
	direction blueprintspatial.Direction
}

// Render processes all renderable entities within a scene and presents them to screen
func (sys GlobalRenderer) Render(cli coldbrew.Client, screen coldbrew.Screen) {
	if sys.logger == nil {
		sys.logger = bark.For("GlobalRenderSystem")
	}

	for _, cam := range cli.Cameras() {
		if !cam.Active() {
			continue
		}

		entry, ok := cli.CameraSceneTracker()[cam]
		var scene coldbrew.Scene
		if ok {
			scene = entry.Scene
		} else {
			sys.logger.Debug("Camera not assigned, attempting active scene 0", "cameraIndex", cam.Index())
			if cli.SceneCount() == 0 {
				sys.logger.Debug("Camera not assigned, all scenes inactive, aborting render", "cameraIndex", cam.Index())
				continue
			}
			scene = cli.ActiveScene(0)
			sys.logger.Debug("Camera not assigned, assigning to active scene", "cameraIndex", cam.Index(), "active scene", scene.Name())
			cli.CameraSceneTracker()[cam] = coldbrew.CameraSceneRecord{
				Scene: scene,
				Tick:  cli.CurrentTick(),
			}
		}

		if !scene.Ready() || !cam.Ready(cli) {
			scene = cli.LoadingScenes()[0]
		}
		// Render backgrounds
		cursor := scene.NewCursor(blueprint.Queries.ParallaxBackground)
		for range cursor.Next() {
			if ok, bgConfig := blueprintclient.Components.ParallaxBackground.GetFromCursorSafe(cursor); ok {
				position := blueprintspatial.Components.Position.GetFromCursor(cursor)
				sprBundle := blueprintclient.Components.SpriteBundle.GetFromCursor(cursor)

				backgroundSprite := coldbrew.MaterializeSprites(sprBundle)[0]
				if sprBundle.Blueprints[0].Config.Active {
					RenderBackground(backgroundSprite, position.Two, bgConfig, cam, scene.Width())
				}
			}
		}

		cursor = scene.NewCursor(blueprint.Queries.SpriteBundle)
		for range cursor.Next() {
			if blueprintclient.Components.ParallaxBackground.CheckCursor(cursor) {
				continue
			}

			sprBundle := blueprintclient.Components.SpriteBundle.GetFromCursor(cursor)
			sprites := coldbrew.MaterializeSprites(sprBundle)
			for i, sprite := range sprites {
				bp := &sprBundle.Blueprints[i]
				pos := blueprintspatial.Components.Position.GetFromCursor(cursor)
				hasDirection, direction := blueprintspatial.Components.Direction.GetFromCursorSafe(cursor)
				if !hasDirection {
					directionV := blueprintspatial.NewDirectionRight()
					direction = &directionV
				}
				hasRot, rot := blueprintspatial.Components.Rotation.GetFromCursorSafe(cursor)
				if !hasRot {
					rotV := blueprintspatial.Rotation(0)
					rot = &rotV
				}
				hasScale, scale := blueprintspatial.Components.Scale.GetFromCursorSafe(cursor)
				if !hasScale {
					scaleV := blueprintspatial.NewScale(1, 1)
					scale = &scaleV
				}

				prio := bp.Config.Priority
				if len(sys.sorted) <= prio {
					newSorted := make([][]RenderItem, prio+1)
					copy(newSorted, sys.sorted)
					sys.sorted = newSorted
				}
				sys.sorted[prio] = append(sys.sorted[prio], RenderItem{sprite: sprite, blueprint: bp, pos: *pos, rot: *rot, scale: *scale, direction: *direction})
			}
		}
		for _, sortedItemBucket := range sys.sorted {
			for _, sortedItem := range sortedItemBucket {
				spr := sortedItem.sprite
				bp := sortedItem.blueprint
				config := bp.Config
				pos := sortedItem.pos
				direction := sortedItem.direction
				rot := sortedItem.rot
				scale := sortedItem.scale

				if config.IgnoreDefaultRenderer {
					continue
				}
				RenderEntity(pos.Two, float64(rot), scale.Two, direction, spr, bp, cam, scene.CurrentTick())
			}
		}
		sys.sorted = make([][]RenderItem, len(sys.sorted))

		cam.PresentToScreen(screen, coldbrew.ClientConfig.CameraBorderSize())
	}
}

// RenderBackground draws a parallax background with proper scrolling behavior
func RenderBackground(
	backgroundSprite coldbrew.Sprite,
	position vector.Two,
	bgConfig *blueprintclient.ParallaxBackground,
	cam coldbrew.Camera,
	sceneWidth int,
) {
	spriteWidth := float64(backgroundSprite.Image().Bounds().Dx())
	bgCount := int(math.Ceil(float64(sceneWidth)/spriteWidth)) + 2

	// Get camera-specific translation
	camIndex := cam.Index()
	currentTrans := &bgConfig.RelativeTranslations[camIndex]

	if bgConfig.DisableLooping {
		bgCount = 1
	}
	// Render the background layers
	for i := 0; i < bgCount; i++ {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(currentTrans.X, currentTrans.Y)
		opts.GeoM.Translate(float64(i)*spriteWidth, 0)
		cam.DrawImage(backgroundSprite.Image(), opts, position)
	}
}

// RenderEntity draws a single entity with all its transformations applied
func RenderEntity(
	position vector.Two,
	rotation float64,
	scale vector.Two,
	direction blueprintspatial.Direction,
	spr coldbrew.Sprite,
	blueprint *blueprintclient.SpriteBlueprint,
	cam coldbrew.Camera,
	currentTick int,
) {
	config := blueprint.Config
	if !config.Active {
		return
	}
	isSpriteSheet := blueprint.HasAnimations()
	if isSpriteSheet {
		RenderSpriteSheetAnimation(
			spr,
			blueprint,
			config.ActiveAnimIndex,
			position,
			rotation,
			scale,
			direction,
			config.Offset,
			config.Static,
			cam,
			currentTick,
			nil,
		)
	} else {
		RenderSprite(spr, position, rotation, scale, config.Offset, direction, config.Static, cam)
	}
}

// RenderEntityFromCursor renders an entity directly from a warehouse cursor
func RenderEntityFromCursor(cursor *warehouse.Cursor, cam coldbrew.Camera, currentTick int) {
	sprBundle := blueprintclient.Components.SpriteBundle.GetFromCursor(cursor)
	sprites := coldbrew.MaterializeSprites(sprBundle)
	for i, sprite := range sprites {
		bp := &sprBundle.Blueprints[i]
		pos := blueprintspatial.Components.Position.GetFromCursor(cursor)
		hasDirection, direction := blueprintspatial.Components.Direction.GetFromCursorSafe(cursor)
		if !hasDirection {
			directionV := blueprintspatial.NewDirectionRight()
			direction = &directionV
		}
		hasRot, rot := blueprintspatial.Components.Rotation.GetFromCursorSafe(cursor)
		if !hasRot {
			rotV := blueprintspatial.Rotation(0)
			rot = &rotV
		}
		hasScale, scale := blueprintspatial.Components.Scale.GetFromCursorSafe(cursor)
		if !hasScale {
			scaleV := blueprintspatial.NewScale(0, 0)
			scale = &scaleV
		}

		isSpriteSheet := bp.HasAnimations()
		if isSpriteSheet {
			RenderSpriteSheetAnimation(
				sprite,
				bp,
				bp.Config.ActiveAnimIndex,
				pos.Two,
				float64(*rot),
				scale.Two,
				*direction,
				bp.Config.Offset,
				bp.Config.Static,
				cam,
				currentTick,
				nil,
			)
		} else {
			RenderSprite(sprite, pos.Two, float64(*rot), scale.Two, bp.Config.Offset, *direction, bp.Config.Static, cam)
		}
	}
}

// RenderSprite draws a static sprite with transformations applied
func RenderSprite(
	sprite coldbrew.Sprite,
	position vector.Two,
	rotation float64,
	scale vector.Two,
	offset vector.Two,
	direction blueprintspatial.Direction,
	static bool,
	cam coldbrew.Camera,
) {
	if scale.X == 0 {
		scale.X = 1
	}
	if scale.Y == 0 {
		scale.Y = 1
	}
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(offset.X, offset.Y)
	if direction.IsLeft() {
		opts.GeoM.Scale(-1, 1)
	}
	// Scale before rotation for upscaled sprites
	opts.GeoM.Scale(scale.X, scale.Y)
	if rotation != 0 {
		opts.GeoM.Rotate(rotation)
	}

	if static {
		cam.DrawImageStatic(sprite.Image(), opts, position)
		return
	}
	cam.DrawImage(sprite.Image(), opts, position)
}

// RenderSpriteSheetAnimation draws an animated sprite with proper frame selection
func RenderSpriteSheetAnimation(
	sheet coldbrew.Sprite,
	spriteBlueprint *blueprintclient.SpriteBlueprint,
	index int,
	position vector.Two,
	rotation float64,
	scale vector.Two,
	direction blueprintspatial.Direction,
	offset vector.Two,
	static bool,
	cam coldbrew.Camera,
	tick int,
	logger *slog.Logger,
) {
	anim := &spriteBlueprint.Animations[index]

	durationInTicks := anim.FrameCount * anim.Speed
	if anim.StartTick == 0 {
		anim.StartTick = tick
	}
	animFinished := tick-durationInTicks >= anim.StartTick

	if animFinished && !anim.Freeze {
		anim.StartTick = tick
	}
	var frameIndex int
	if animFinished && anim.Freeze {
		frameIndex = anim.FrameCount - 1
	} else {
		frameIndex = ((tick - anim.StartTick) / anim.Speed) % anim.FrameCount
	}
	if scale.X == 0 {
		scale.X = 1
	}
	if scale.Y == 0 {
		scale.Y = 1
	}
	frame := GetAnimationFrame(sheet, *anim, frameIndex, logger)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(offset.X, offset.Y)
	opts.GeoM.Translate(anim.PositionOffset.X, anim.PositionOffset.Y)

	if direction.IsLeft() {
		opts.GeoM.Scale(-1, 1)
	}
	// Scale before rotation for upscaled sprites
	opts.GeoM.Scale(scale.X, scale.Y)
	if rotation != 0 {
		opts.GeoM.Rotate(rotation)
	}
	if static {
		cam.DrawImageStatic(frame, opts, position)
		return
	}
	cam.DrawImage(frame, opts, position)
}

// GetAnimationFrame extracts a single frame from a sprite sheet based on animation data
func GetAnimationFrame(sheet coldbrew.Sprite, anim blueprintclient.AnimationData, frameIndex int, logger *slog.Logger) *ebiten.Image {
	sx := frameIndex * anim.FrameWidth
	sy := anim.RowIndex * anim.FrameHeight
	frame := sheet.Image().SubImage(image.Rect(sx, sy, sx+anim.FrameWidth, sy+anim.FrameHeight)).(*ebiten.Image)
	return frame
}
