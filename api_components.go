package coldbrew

import (
	blueprint_client "github.com/TheBitDrifter/blueprint/client"
	"github.com/hajimehoshi/ebiten/v2"
)

// How to expose this to our scene without dep?
// Do we abuse our generic cache?

type Sprite struct {
	Name  string
	Image *ebiten.Image
}

type ISprite interface {
	Render(target Sprite, options blueprint_client.RenderOptions)
}

type IRenderOptions interface {
	Translation() (float64, float64)
}
