package coldbrew

import (
	"fmt"

	blueprintclient "github.com/TheBitDrifter/blueprint/client"
)

// MaterializeSprites converts a bundle of sprite blueprints into concrete Sprite objects
// It skips any blueprint with an empty location key
func MaterializeSprites(spriteBundle *blueprintclient.SpriteBundle) []Sprite {
	var sprites []Sprite
	for i := range spriteBundle.Blueprints {
		spriteBlueprint := &spriteBundle.Blueprints[i]
		if spriteBlueprint.Location.Index != 0 {
			spr := globalSpriteCache.GetItem32(spriteBlueprint.Location.Index)
			sprites = append(sprites, spr)
			continue
		}
		if spriteBlueprint.Location.Key != "" {
			if idx, ok := globalSpriteCache.GetIndex(spriteBlueprint.Location.Key); ok {
				spriteBlueprint.Location.Index = uint32(idx)
				spr := globalSpriteCache.GetItem(idx)
				sprites = append(sprites, spr)
				continue
			}
		}
	}
	return sprites
}

// MaterializeSounds converts a collection of sound blueprints into concrete Sound objects
// It skips any blueprint with an empty location key
func MaterializeSounds(soundBundle *blueprintclient.SoundBundle) []Sound {
	var sounds []Sound
	for i := range soundBundle.Blueprints {
		soundBlueprint := &soundBundle.Blueprints[i]
		if soundBlueprint.Location.Index != 0 {
			snd := globalSoundCache.GetItem32(soundBlueprint.Location.Index)
			sounds = append(sounds, snd)
			continue
		}
		if soundBlueprint.Location.Key != "" {
			if idx, ok := globalSoundCache.GetIndex(soundBlueprint.Location.Key); ok {
				soundBlueprint.Location.Index = uint32(idx)
				snd := globalSoundCache.GetItem(idx)
				sounds = append(sounds, snd)
				continue
			}
		}
	}
	return sounds
}

// MaterializeSound finds and returns a specific Sound object from a bundle based on the provided SoundConfig
// Returns an error if the sound is not found
func MaterializeSound(soundBundle *blueprintclient.SoundBundle, sc blueprintclient.SoundConfig) (Sound, error) {
	for i := range soundBundle.Blueprints {
		soundBlueprint := &soundBundle.Blueprints[i]
		location := soundBlueprint.Location

		if location.Index != 0 && sc.Path == location.Key {
			return globalSoundCache.GetItem32(location.Index), nil
		}

		if location.Key != "" && sc.Path == location.Key {
			if idx, ok := globalSoundCache.GetIndex(location.Key); ok {
				soundBlueprint.Location.Index = uint32(idx)
				return globalSoundCache.GetItem(idx), nil
			}
		}
	}
	return Sound{}, fmt.Errorf("%v not found", sc)
}
