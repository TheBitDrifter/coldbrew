package coldbrew

import (
	"errors"

	blueprintclient "github.com/TheBitDrifter/blueprint/client"
)

// MaterializeSprites converts a bundle of sprite blueprints into concrete Sprite objects
// It skips any blueprint with an empty location key
func MaterializeSprites(spriteBundle blueprintclient.SpriteBundle) []Sprite {
	var sprites []Sprite
	for _, spriteBlueprint := range spriteBundle.Blueprints {
		location := spriteBlueprint.Location
		if location.Key == "" {
			continue
		}
		spr := globalSpriteCache.GetItem32(location.Index)
		sprites = append(sprites, spr)
	}
	return sprites
}

// MaterializeSounds converts a collection of sound blueprints into concrete Sound objects
// It skips any blueprint with an empty location key
func MaterializeSounds(soundBundle blueprintclient.SoundBundle) []Sound {
	var sounds []Sound
	for _, soundBlueprint := range soundBundle.Blueprints {
		location := soundBlueprint.Location
		if location.Key == "" {
			continue
		}
		snd := globalSoundCache.GetItem32(location.Index)
		sounds = append(sounds, snd)
	}
	return sounds
}

func MaterializeSound(soundBundle blueprintclient.SoundBundle, sc blueprintclient.SoundConfig) (Sound, error) {
	for _, soundBlueprint := range soundBundle.Blueprints {
		location := soundBlueprint.Location
		if location.Key == "" {
			continue
		}
		if sc.Path == location.Key {
			return globalSoundCache.GetItem32(location.Index), nil
		}
	}
	return Sound{}, errors.New("sound not found")
}
