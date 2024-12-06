package coldbrew

type config struct {
	Title               string
	MaxSpritesCached    int
	MaxScenesCached     int
	MandatoryLoadScreen bool
	DefaultRenderer     bool
	ShowClientData      bool
	Resolution          struct {
		X, Y int
	}
	WindowSize struct {
		X, Y int
	}
}

var ClientConfig = config{
	Title:               "Hello!",
	MaxSpritesCached:    1000,
	DefaultRenderer:     true,
	ShowClientData:      true,
	MandatoryLoadScreen: true,
	Resolution: struct {
		X, Y int
	}{
		X: 640,
		Y: 320,
	},
}

const MaxScreenSplit = 8
