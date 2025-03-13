package scenes

import (
	"github.com/TheBitDrifter/blueprint"
	"github.com/TheBitDrifter/warehouse"
)

var TertiarySceneName = "Tert"

func TertiaryScene(height, width int, sto warehouse.Storage) error {
	err := blueprint.NewParallaxBackgroundBuilder(sto).
		AddLayer("backgrounds/city/far2.png", 0.0, 0.0).
		Build()
	if err != nil {
		return nil
	}
	addTerrain(sto)
	return nil
}
