package coldbrew

import (
	"fmt"

	blueprint_client "github.com/TheBitDrifter/blueprint/client"
	"github.com/TheBitDrifter/warehouse"
)

type spriteLocator struct{}

var _ AssetLocator = spriteLocator{}

func (locator spriteLocator) LocationFromIndex(idx int, cursor *warehouse.Cursor) (warehouse.CacheLocation, error) {
	if blueprint_client.Components.SpriteLocations.Small.Check(cursor) {
		sprites := blueprint_client.Components.SpriteLocations.Small.GetFromCursor(cursor)
		if idx < 0 || idx >= len(sprites) {
			return warehouse.CacheLocation{}, fmt.Errorf("todo", "proper error")
		}
		return sprites[idx], nil
	}
	if blueprint_client.Components.SpriteLocations.Med.Check(cursor) {
		sprites := blueprint_client.Components.SpriteLocations.Med.GetFromCursor(cursor)
		if idx < 0 || idx >= len(sprites) {
			return warehouse.CacheLocation{}, fmt.Errorf("todo", "proper error")
		}
		return sprites[idx], nil
	}
	if blueprint_client.Components.SpriteLocations.Large.Check(cursor) {
		sprites := blueprint_client.Components.SpriteLocations.Large.GetFromCursor(cursor)
		if idx < 0 || idx >= len(sprites) {
			return warehouse.CacheLocation{}, fmt.Errorf("todo", "proper error")
		}
		return sprites[idx], nil
	}
	if blueprint_client.Components.SpriteLocations.XL.Check(cursor) {
		sprites := blueprint_client.Components.SpriteLocations.XL.GetFromCursor(cursor)
		if idx < 0 || idx >= len(sprites) {
			return warehouse.CacheLocation{}, fmt.Errorf("todo", "proper error")
		}
		return sprites[idx], nil
	}

	return warehouse.CacheLocation{}, fmt.Errorf("todo")
}

func (locator spriteLocator) Locations(cursor *warehouse.Cursor) []warehouse.CacheLocation {
	var locations []warehouse.CacheLocation

	if blueprint_client.Components.SpriteLocations.Small.Check(cursor) {
		sLocations := blueprint_client.Components.SpriteLocations.Small.GetFromCursor(cursor)
		for _, location := range sLocations {
			locations = append(locations, location)
		}
	}
	if blueprint_client.Components.SpriteLocations.Med.Check(cursor) {
		mLocations := blueprint_client.Components.SpriteLocations.Med.GetFromCursor(cursor)
		for _, location := range mLocations {
			locations = append(locations, location)
		}
	}
	if blueprint_client.Components.SpriteLocations.Large.Check(cursor) {
		lLocations := blueprint_client.Components.SpriteLocations.Large.GetFromCursor(cursor)
		for _, location := range lLocations {
			locations = append(locations, location)
		}
	}
	if blueprint_client.Components.SpriteLocations.XL.Check(cursor) {
		xlLocations := blueprint_client.Components.SpriteLocations.XL.GetFromCursor(cursor)
		for _, location := range xlLocations {
			locations = append(locations, location)
		}
	}
	return locations
}

func (locator spriteLocator) MutableLocations(cursor *warehouse.Cursor) []*warehouse.CacheLocation {
	var locations []*warehouse.CacheLocation

	if blueprint_client.Components.SpriteLocations.Small.Check(cursor) {
		sLocations := blueprint_client.Components.SpriteLocations.Small.GetFromCursor(cursor)
		for i := range sLocations {
			locations = append(locations, &sLocations[i])
		}
	}
	if blueprint_client.Components.SpriteLocations.Med.Check(cursor) {
		mLocations := blueprint_client.Components.SpriteLocations.Med.GetFromCursor(cursor)
		for i := range mLocations {
			locations = append(locations, &mLocations[i])
		}
	}
	if blueprint_client.Components.SpriteLocations.Large.Check(cursor) {
		lLocations := blueprint_client.Components.SpriteLocations.Large.GetFromCursor(cursor)
		for i := range lLocations {
			locations = append(locations, &lLocations[i])
		}
	}
	if blueprint_client.Components.SpriteLocations.XL.Check(cursor) {
		xlLocations := blueprint_client.Components.SpriteLocations.XL.GetFromCursor(cursor)
		for i := range xlLocations {
			locations = append(locations, &xlLocations[i])
		}
	}

	return locations
}

func (locator spriteLocator) AllMutableLocations(cli Client) []*warehouse.CacheLocation {
	cur := cli.NewCursor(SpriteQuery)
	var locations []*warehouse.CacheLocation
	for cur.Next() {
		locations = append(locations, locator.MutableLocations(cur)...)
	}

	return locations
}
