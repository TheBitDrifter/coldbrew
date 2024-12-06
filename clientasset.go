package coldbrew

import "io/fs"

type assetManager struct {
	*spriteLoader
	*soundLoader
}

func newAssetManager(embeddedFS fs.FS) *assetManager {
	return &assetManager{
		spriteLoader: NewSpriteLoader(embeddedFS),
		soundLoader:  NewSoundLoader(embeddedFS),
	}
}
