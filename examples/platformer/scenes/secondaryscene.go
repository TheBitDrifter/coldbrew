package scenes

import (
	"github.com/TheBitDrifter/blueprint"
	"github.com/TheBitDrifter/warehouse"
)

var SecondarySceneName = "SecondaryScene"

func SecondaryScene(height, width int, sto warehouse.Storage) error {
	err := blueprint.NewParallaxBackgroundBuilder(sto).
		AddLayer("backgrounds/city/sky.png", 0.0, 0.0).
		Build()
	if err != nil {
		return nil
	}
	addTerrain(sto)
	return nil
}
