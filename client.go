package coldbrew

import (
	"fmt"

	"github.com/TheBitDrifter/blueprint"
	"github.com/TheBitDrifter/table"
	"github.com/TheBitDrifter/warehouse"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type client struct {
	currentScene *scene
	scenes       warehouse.Cache[*scene]
	assets       []Asset
}

func NewClient() Client {
	return &client{}
}

func (cc *client) Start() error {
	err := ebiten.RunGame(cc)
	if err != nil {
		return err
	}
	return nil
}

func (cli *client) NewScene(name string, height, width int, s blueprint.Plan, renderSystems []RenderSystem) error {
	cli.scenes = warehouse.FactoryNewCache[*scene](ClientConfig.MaxScenesCached)
	schema := table.Factory.NewSchema()
	storage := warehouse.Factory.NewStorage(schema)
	err := s(storage)
	if err != nil {
		return err
	}
	newScene := &scene{storage: storage}
	if ClientConfig.DefaultRenderer {
		renderSystemsWithDefaultRenderer := make([]RenderSystem, len(renderSystems)+2)
		renderSystemsWithDefaultRenderer[0] = baseRenderSystem{}
		renderSystemsWithDefaultRenderer[1] = parallaxRenderSystem{}
		renderSystemsWithDefaultRenderer = append(renderSystemsWithDefaultRenderer, renderSystems...)
		renderSystems = renderSystemsWithDefaultRenderer
	}
	newScene.systems.renderers = renderSystems
	newScene.index, err = cli.scenes.Register(name, newScene)
	newScene.height = height
	newScene.width = width
	if cli.currentScene == nil {
		cli.currentScene = newScene
	}
	cli.assets = make([]Asset, 1)
	cli.assets[imageK] = Asset{
		Name:    "Image",
		Locator: spriteLocator{},
		Loader:  spriteLoader{},
	}
	return nil
}

func (cli client) NewCursor(query warehouse.Query) *warehouse.Cursor {
	return warehouse.Factory.NewCursor(query, cli.currentScene.storage)
}

func (cli *client) GetSprites(cursor *warehouse.Cursor) []*Sprite {
	var images []*Sprite
	locations := cli.assets[imageK].Locator.Locations(cursor)
	for _, location := range locations {
		if location.Key == "" {
			continue
		}
		img := globalSpriteCache.GetItem32(location.Index)
		images = append(images, img)
	}
	return images
}

func (cli *client) GetSprite(idx int, cursor *warehouse.Cursor) (*Sprite, error) {
	location, err := cli.assets[imageK].Locator.LocationFromIndex(idx, cursor)
	if err != nil {
		return nil, err
	}
	sprite := globalSpriteCache.GetItem32(location.Index)
	return sprite, nil
}

func (*client) dataLog(screen Screen) {
	stats := fmt.Sprintf("FRAMES: %v\nTICKS: %v", ebiten.ActualFPS(), ebiten.ActualTPS())
	ebitenutil.DebugPrint(screen.Sprite.Image, stats)
}
