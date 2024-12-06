package coldbrew

import (
	"github.com/TheBitDrifter/warehouse"
	"github.com/hajimehoshi/ebiten/v2"
)

func (cli *client) Draw(ebitenScreen *ebiten.Image) {
	if !cli.currentScene.loaded {
		return
	}

	screen := Screen{
		Sprite{Name: "screen", Image: ebitenScreen},
	}
	for _, renderSys := range cli.currentScene.systems.renderers {
		renderSys.Render(cli, screen)
	}
	if ClientConfig.ShowClientData {
		cli.dataLog(screen)
	}
}

func (cli *client) Update() error {
	if !cli.currentScene.loaded {
		cli.load()
	}

	cursor := warehouse.Factory.NewCursor(CameraQuery, cli.currentScene.storage)
	for cursor.Next() {
		cam := cameraComponent.GetFromCursor(cursor)
		cam.Positions.Local.X += 0.88888888
		cam.Positions.Local.Y += 0.88888888
	}
	return nil
}

func (cli *client) load() error {
	for _, asset := range cli.assets {
		locations := asset.Locator.AllMutableLocations(cli)
		err := asset.Loader.Load(locations)
		if err != nil {
			return err
		}
	}
	cli.currentScene.loaded = true
	return nil
}

func (cli *client) Layout(int, int) (int, int) {
	return ClientConfig.Resolution.X, ClientConfig.Resolution.Y
}
