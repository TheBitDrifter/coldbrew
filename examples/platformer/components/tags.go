package components

import "github.com/TheBitDrifter/warehouse"

var (
	BlockTerrainTag    = warehouse.FactoryNewComponent[struct{}]()
	PlatformTerrainTag = warehouse.FactoryNewComponent[struct{}]()
	MusicTag           = warehouse.FactoryNewComponent[struct{}]()
)
