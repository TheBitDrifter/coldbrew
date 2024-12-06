package coldbrew

import (
	blueprint_client "github.com/TheBitDrifter/blueprint/client"
	"github.com/hajimehoshi/ebiten/v2"
)

type renderableCamera struct {
	*blueprint_client.Camera
	surface *Sprite
}

func (c *renderableCamera) RenderSprite(sprite Sprite, spritePosX, spritePosY float64, opts blueprint_client.RenderOptions) {
	opts = c.translateOntoCamera(spritePosX, spritePosY, opts)
	c.surface.Render(sprite, opts)
}

func (c *renderableCamera) translateOntoCamera(spritePosX, spritePosY float64, opts blueprint_client.RenderOptions) blueprint_client.RenderOptions {
	opts = blueprint_client.RenderOptions{
		TranslateX: opts.TranslateX + (-c.Positions.Local.X - spritePosX),
		TranslateY: opts.TranslateY + (-c.Positions.Local.Y - spritePosY),
	}
	return opts
}

func (c *renderableCamera) Render(screen Screen) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(c.Positions.Screen.X, c.Positions.Screen.Y)
	screen.Image.DrawImage(c.surface.Image, opts)
}

func Clear(c *renderableCamera) {
	c.surface.Image.Clear()
}
