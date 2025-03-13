package clientsystems

import (
	blueprintclient "github.com/TheBitDrifter/blueprint/client"
	"github.com/TheBitDrifter/coldbrew"
	"github.com/TheBitDrifter/coldbrew/examples/platformer/components"
	"github.com/TheBitDrifter/warehouse"
)

type MusicSystem struct{}

func (sys MusicSystem) Run(lc coldbrew.LocalClient, scene coldbrew.Scene) error {
	// Setup query and cursor for music
	musicQuery := warehouse.Factory.NewQuery().And(components.MusicTag)
	cursor := scene.NewCursor(musicQuery)

	// There's only one but iterate nonetheless
	for range cursor.Next() {
		soundBundle := blueprintclient.Components.SoundBundle.GetFromCursor(cursor)

		// Get the actual sounds from the bundle
		sounds := coldbrew.MaterializeSounds(soundBundle)

		// Since we only have one song and one audio player we use 0 for both
		player := sounds[0].GetPlayer(0)

		// Loop if needed
		if !player.IsPlaying() {
			player.Rewind()
			player.Play()
		}
	}
	return nil
}
