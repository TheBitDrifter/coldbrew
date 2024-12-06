package coldbrew

import (
	"image/color"
	"log"

	"github.com/TheBitDrifter/blueprint"
	blueprint_client "github.com/TheBitDrifter/blueprint/client"
	"github.com/hajimehoshi/ebiten/v2"
)

var cameraSurfaces [MaxScreenSplit]*Sprite

var colors = []color.RGBA{
	{255, 0, 0, 255},   // Red
	{0, 255, 0, 255},   // Green
	{0, 0, 255, 255},   // Blue
	{255, 255, 0, 255}, // Yellow
	{255, 0, 255, 255}, // Magenta
	{0, 255, 255, 255}, // Cyan
	{255, 128, 0, 255}, // Orange
	{128, 0, 255, 255}, // Purple
}

type baseRenderSystem struct{}

func (baseRenderSystem) Render(cli Client, screen Screen) {
	spriteCursor := cli.NewCursor(ActiveSpriteQuery)
	cameraCursor := cli.NewCursor(CameraQuery)

	i := 0
	for cameraCursor.Next() {
		cam := blueprint_client.Components.Camera.GetFromCursor(cameraCursor)
		renderCam := renderableCamera{
			Camera:  cam,
			surface: cameraSurfaces[i],
		}
		if renderCam.surface == nil {
			surface := &Sprite{
				Image: ebiten.NewImage(cam.Height, cam.Width),
			}
			surface.Image.Fill(colors[i])
			renderCam.surface = surface
		}

		for spriteCursor.Next() {
			spritePos := blueprint.Components.Position.GetFromCursor(spriteCursor)
			opts := &blueprint_client.RenderOptions{}
			if renderOptionsComponent.Check(spriteCursor) {
				// Dereference opts to make a value/copy type.
				// This prevents repeated mutations to the opts.
				existingOptsClone := *renderOptionsComponent.GetFromCursor(spriteCursor)
				opts = &existingOptsClone
				opts.TranslateX += spritePos.X
				opts.TranslateY += spritePos.Y
			}
			activeSpriteIndex := activeSpriteComponent.GetFromCursor(spriteCursor).Index
			sprite, err := cli.GetSprite(activeSpriteIndex, spriteCursor)
			if err != nil {
				log.Fatal("todo")
			}
			renderCam.RenderSprite(*sprite, spritePos.X, spritePos.Y, *opts)
		}
		renderCam.Render(screen)
		i++
	}
}

type parallaxRenderSystem struct{}

func (parallaxRenderSystem) Render(cli Client, screen Screen) {
	backgroundCursor := cli.NewCursor(ParallaxQuery)
	for backgroundCursor.Next() {
		if !renderOptionsComponent.Check(backgroundCursor) {
			backgroundEntity, err := backgroundCursor.CurrentEntity()
			if err != nil {
				log.Fatal(err)
			}
			err = backgroundEntity.EnqueueAddComponent(renderOptionsComponent)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	for backgroundCursor.Next() {
		parallaxConfig := blueprint_client.Components.ParallaxBackground.GetFromCursor(backgroundCursor)
		opts := renderOptionsComponent.GetFromCursor(backgroundCursor)

		camCursor := cli.NewCursor(CameraQuery)
		for camCursor.Next() {
			cam := blueprint_client.Components.Camera.GetFromCursor(camCursor)
			opts.TranslateX = cam.Positions.Local.X * parallaxConfig.SpeedX
			opts.TranslateY = cam.Positions.Local.Y * parallaxConfig.SpeedY
		}
	}
}
