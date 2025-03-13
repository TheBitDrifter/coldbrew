package clientsystems

import (
	"github.com/TheBitDrifter/blueprint"
	blueprintclient "github.com/TheBitDrifter/blueprint/client"
	"github.com/TheBitDrifter/coldbrew"
)

// BackgroundScrollSystem handles parallax scrolling effects for backgrounds based on camera movement
type BackgroundScrollSystem struct{}

// Run updates the position of parallax backgrounds relative to camera positions
// It calculates offsets based on configured speed values and camera world position
func (BackgroundScrollSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	activeCameras := cli.ActiveCamerasFor(scene)
	for _, cam := range activeCameras {
		cursor := scene.NewCursor(blueprint.Queries.ParallaxBackground)
		_, worldPosition := cam.Positions()
		for range cursor.Next() {
			config := blueprintclient.Components.ParallaxBackground.GetFromCursor(cursor)
			// Apply inverse movement based on camera position and configured speeds
			config.RelativeTranslations[cam.Index()].X = worldPosition.X * config.SpeedX * -1
			config.RelativeTranslations[cam.Index()].Y = worldPosition.Y * config.SpeedY * -1
		}
	}
	return nil
}
