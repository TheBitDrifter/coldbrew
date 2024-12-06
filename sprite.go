package coldbrew

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Ensure sprite implements Sprite interface
var _ Sprite = &sprite{}

// Sprite manages image state and rendering properties
// It contains a reference to the drawable *ebiten.Image
type Sprite interface {
	// Name returns the sprite's identifier
	Name() string

	// Image returns the underlying ebiten.Image
	Image() *ebiten.Image

	// Draw renders the sprite to the target image with the provided options
	Draw(*ebiten.Image, *ebiten.DrawImageOptions)
}

// sprite implements the Sprite interface
type sprite struct {
	name  string
	image *ebiten.Image
}

// Image returns the underlying ebiten.Image
func (s *sprite) Image() *ebiten.Image {
	return s.image
}

// Name returns the sprite's identifier
func (s *sprite) Name() string {
	return s.name
}

// Draw renders the sprite to the target image
func (s *sprite) Draw(target *ebiten.Image, options *ebiten.DrawImageOptions) {
	target.DrawImage(s.image, options)
}
