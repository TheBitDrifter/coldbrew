package coldbrew

import "github.com/TheBitDrifter/warehouse"

var globalSpriteCache = warehouse.FactoryNewCache[Sprite](ClientConfig.MaxSpritesCached).(*warehouse.SimpleCache[Sprite])
