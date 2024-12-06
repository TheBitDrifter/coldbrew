package coldbrew

import (
	"github.com/TheBitDrifter/blueprint"
	"github.com/TheBitDrifter/warehouse"
	"github.com/hajimehoshi/ebiten/v2"
)

type Client interface {
	Start() error
	NewScene(string, int, int, blueprint.Plan, []RenderSystem) error
	NewCursor(warehouse.Query) *warehouse.Cursor
	GetSprites(cursor *warehouse.Cursor) []*Sprite
	GetSprite(int, *warehouse.Cursor) (*Sprite, error)
	// ChangeSceneByIndex(int) error
	// ChangeSceneByName(string) error
	clientConfigurator
	ebiten.Game
}

type clientConfigurator interface {
	SetWindowSize(int, int)
	SetResolution(int, int)
	SetResizable(bool)
	SetFullScreen(bool)

	SetTitle(string)
}

type RenderSystem interface {
	Render(cli Client, screen Screen)
}

type AssetLocator interface {
	LocationFromIndex(int, *warehouse.Cursor) (warehouse.CacheLocation, error)
	Locations(*warehouse.Cursor) []warehouse.CacheLocation
	MutableLocations(*warehouse.Cursor) []*warehouse.CacheLocation
	AllMutableLocations(Client) []*warehouse.CacheLocation
}

type AssetLoader interface {
	Load([]*warehouse.CacheLocation) error
}

type Screen struct {
	Sprite
}
