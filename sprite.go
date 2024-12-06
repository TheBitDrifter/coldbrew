package coldbrew

import (
	blueprint_client "github.com/TheBitDrifter/blueprint/client"
	"github.com/hajimehoshi/ebiten/v2"
)

var _ ISprite = &Sprite{}

func (s Sprite) Render(target Sprite, options blueprint_client.RenderOptions) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(options.TranslateX, options.TranslateY)
	s.Image.DrawImage(target.Image, opts)
}
