package coldbrew

import (
	"github.com/TheBitDrifter/blueprint"
	blueprint_client "github.com/TheBitDrifter/blueprint/client"
	"github.com/TheBitDrifter/warehouse"
)

var SpriteQuery = func() warehouse.Query {
	query := warehouse.Factory.NewQuery()
	query.Or(
		blueprint_client.Components.SpriteLocations.Small,
		blueprint_client.Components.SpriteLocations.Med,
		blueprint_client.Components.SpriteLocations.Large,
		blueprint_client.Components.SpriteLocations.XL,
	)
	return query
}()

var ActiveSpriteQuery = func() warehouse.Query {
	query := warehouse.Factory.NewQuery()
	query.Or(
		query.And(blueprint_client.Components.SpriteLocations.Small, blueprint.Components.Position, activeSpriteComponent),
		query.And(blueprint_client.Components.SpriteLocations.Med, blueprint.Components.Position, activeSpriteComponent),
		query.And(blueprint_client.Components.SpriteLocations.Large, blueprint.Components.Position, activeSpriteComponent),
		query.And(blueprint_client.Components.SpriteLocations.XL, blueprint.Components.Position, activeSpriteComponent),
	)
	return query
}()

var CameraQuery = func() warehouse.Query {
	query := warehouse.Factory.NewQuery()
	query.And(cameraComponent)
	return query
}()

var ParallaxQuery = func() warehouse.Query {
	query := warehouse.Factory.NewQuery()
	query.And(blueprint_client.Components.ParallaxBackground)
	return query
}()
