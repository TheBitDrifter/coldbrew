package clientsystems

import (
	"github.com/TheBitDrifter/coldbrew"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type FullScreenSystem struct{}

func (FullScreenSystem) Run(cli coldbrew.Client) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	return nil
}
