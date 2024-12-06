package coldbrew

type Asset struct {
	Name    string
	Locator AssetLocator
	Loader  AssetLoader
}

type assetKey int

const (
	imageK assetKey = 0
	soundK
)
