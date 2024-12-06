package coldbrew

import (
	"fmt"
	"log"

	"github.com/TheBitDrifter/warehouse"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var _ AssetLoader = &spriteLoader{}

type spriteLoader struct{}

func (loader spriteLoader) Load(locations []*warehouse.CacheLocation) error {
	for _, location := range locations {
		if location.Key == "" {
			continue
		}
		imageIndex, ok := globalSpriteCache.GetIndex(location.Key)
		if ok {
			location.Index = uint32(imageIndex)
			continue
		}
		spr, err := loader.loadSpriteFromPath(location.Key)
		if err != nil {
			log.Println("ey yo", err)
			return err
		}
		index, err := globalSpriteCache.Register(location.Key, spr)
		if err != nil {
			log.Println("ey yo", err)
			return err
		}
		location.Index = uint32(index)
	}
	return nil
}

func (loader spriteLoader) loadSpriteFromPath(path string) (Sprite, error) {
	// Todo: change for build?
	// if BuildTime == "true" {
	//
	// 	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	//
	// 	if err != nil {
	// 		log.Fatalf("Asset Error: %v\n", err)
	// 	}
	//
	// 	path = dir + path
	// 	path = strings.ReplaceAll(path, "./", "/")
	// }

	img, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("images/%s.png", path))
	if err != nil {
		return Sprite{}, err
	}

	return Sprite{
		Image: img,
	}, nil
}
