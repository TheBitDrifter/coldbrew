package coldbrew

import "github.com/hajimehoshi/ebiten/v2"

func (cc *client) SetWindowSize(x, y int) {
	ClientConfig.WindowSize.X = x
	ClientConfig.WindowSize.Y = y
	ebiten.SetWindowSize(x, y)
}

func (cc *client) SetResolution(x, y int) {
	ClientConfig.Resolution.X = x
	ClientConfig.Resolution.Y = y
}

func (cc *client) SetResizable(isResizable bool) {
	if isResizable {
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
		return
	}
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
}

func (cc *client) SetFullScreen(fullscreen bool) {
	ebiten.SetFullscreen(fullscreen)
}

func (cc *client) SetTitle(title string) {
	ClientConfig.Title = title
	ebiten.SetWindowTitle(title)
}
