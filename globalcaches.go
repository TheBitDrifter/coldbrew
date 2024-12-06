package coldbrew

import "github.com/TheBitDrifter/warehouse"

var (
	globalSoundCache  = warehouse.FactoryNewCache[Sound](ClientConfig.maxSoundsCached)
	globalSpriteCache = warehouse.FactoryNewCache[Sprite](ClientConfig.maxSpritesCached)
)
