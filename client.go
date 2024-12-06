package coldbrew

import (
	"errors"
	"io/fs"
	"log"

	"github.com/TheBitDrifter/blueprint"
	blueprintclient "github.com/TheBitDrifter/blueprint/client"
	"github.com/TheBitDrifter/blueprint/vector"
	"github.com/TheBitDrifter/table"
	"github.com/TheBitDrifter/warehouse"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/basicfont"
)

var (
	tick        = 0
	_    Client = &client{}
)

// Client manages game state, rendering, and input
type Client interface {
	LocalClient
	SceneManager
	CameraManager
}

type LocalClient interface {
	Start() error
	CameraUtility
	TickManager
	InputManager
	CameraManager
	SystemManager
	ConfigManager
	ebiten.Game
}

type client struct {
	*tickManager
	*inputManager
	*cameraUtility
	*systemManager
	*sceneManager
	*configManager
	*assetManager
}

// NewClient creates a new client with specified resolution and cache settings
func NewClient(baseResX, baseResY, maxScenesCached int, embeddedFS fs.FS) Client {
	cli := &client{
		tickManager:   newTickManager(),
		cameraUtility: newCameraUtility(),
		systemManager: &systemManager{},
		configManager: newConfigManager(),
		sceneManager:  newSceneManager(maxScenesCached),
		assetManager:  newAssetManager(embeddedFS),
	}
	cli.inputManager = newInputManager(cli)
	ClientConfig.baseResolution.x = baseResX
	ClientConfig.baseResolution.y = baseResY
	ClientConfig.windowSize.x = baseResX
	ClientConfig.windowSize.y = baseResY
	ebiten.SetWindowSize(baseResX, baseResY)
	return cli
}

// Start initializes and runs the game loop.
func (cli *client) Start() error {
	if len(cli.loadingScenes) == 0 {
		cli.loadingScenes = append(cli.loadingScenes, defaultLoadingScene)
	}
	err := ebiten.RunGame(cli)
	if err != nil {
		return err
	}
	return nil
}

// LoadScenes loads active scenes
func (cli *client) LoadScenes() error {
	for _, scene := range cli.activeScenes {
		if !scene.IsLoaded() && !scene.IsLoading() {
			if err := cli.load(scene, globalSpriteCache, globalSoundCache); err != nil {
				cli.cacheBust.Store(true)
			}
		}
	}
	return nil
}

func (cli *client) Update() error {
	if inpututil.IsKeyJustReleased(ClientConfig.DebugKey()) && !isProd {
		ClientConfig.DebugVisual = !ClientConfig.DebugVisual
	}
	for _, s := range cli.activeScenes {
		_, err := s.ExecutePlan()
		if err != nil {
			return err
		}
	}
	go cli.LoadScenes()

	if cli.cacheBust.Load() {
		cli.cacheBust.Store(false)
		swapCacheSpr := warehouse.FactoryNewCache[Sprite](ClientConfig.maxSpritesCached)
		swapCacheSnd := warehouse.FactoryNewCache[Sound](ClientConfig.maxSoundsCached)
		for _, s := range cli.activeScenes {
			err := cli.load(s, swapCacheSpr, swapCacheSnd)
			if err != nil {
				log.Println("scene cannot fit  assets given current cache config, maxSpritesCached", ClientConfig.maxSpritesCached)
				log.Println("scene cannot fit  assets given current cache config, maxSpritesCached", ClientConfig.maxSoundsCached)
				return err
			}
		}
		globalSpriteCache = swapCacheSpr
		globalSoundCache = swapCacheSnd
	}
	cli.capturers.keyboard.Capture()
	cli.capturers.mouse.Capture()
	cli.capturers.gamepad.Capture()
	cli.capturers.touch.Capture()

	for _, globalClientSystem := range cli.globalClientSystems {
		err := globalClientSystem.Run(cli)
		if err != nil {
			return err
		}
	}
	for _, activeScene := range cli.activeScenes {
		cameraReady := true
		cameras := cli.ActiveCamerasFor(activeScene)
		for _, cam := range cameras {
			if !cam.Ready(cli) {
				cameraReady = false
			}
		}
		if !cameraReady || !activeScene.Ready() {
			loadingScene := cli.LoadingScenes()[0]
			for _, coreSys := range loadingScene.CoreSystems() {
				err := coreSys.Run(loadingScene, 1.0/float64(ClientConfig.tps))
				if err != nil {
					return err
				}
			}
			for _, clientSys := range loadingScene.ClientSystems() {
				clientSys.Run(cli, loadingScene)
			}
		}
		if activeScene.Ready() {
			for _, coreSys := range activeScene.CoreSystems() {
				err := coreSys.Run(activeScene, 1.0/float64(ClientConfig.tps))
				if err != nil {
					return err
				}
			}
			for _, clientSys := range activeScene.ClientSystems() {
				clientSys.Run(cli, activeScene)
			}
		}
	}
	for _, clientGlobalSecondarySystem := range cli.globalClientSecondarySystems {
		err := clientGlobalSecondarySystem.Run(cli)
		if err != nil {
			return err
		}
	}
	tick++
	return nil
}

func (cli *client) Layout(int, int) (int, int) {
	return ClientConfig.resolution.x, ClientConfig.resolution.y
}

func (cli *client) Draw(image *ebiten.Image) {
	screen := Screen{
		sprite{name: "screen", image: image},
	}
	for _, renderSys := range cli.globalRenderers {
		renderSys.Render(cli, screen)
	}

	for _, activeScene := range cli.activeScenes {
		renderers := activeScene.Renderers()
		cameraReady := true
		cameras := cli.ActiveCamerasFor(activeScene)
		for _, cam := range cameras {
			if !cam.Ready(cli) {
				cameraReady = false
			}
		}
		if !activeScene.Ready() || !cameraReady {
			loadingScene := cli.LoadingScenes()[0]
			for _, renderSys := range loadingScene.Renderers() {
				renderSys.Render(activeScene, screen, cli)
			}
		}
		for _, renderSys := range renderers {
			renderSys.Render(activeScene, screen, cli)
		}
	}
}

func (cli client) CameraSceneTracker() CameraSceneTracker {
	return cli.cameraSceneTracker
}

func (cli client) Cameras() [MaxSplit]Camera {
	return cli.cameras
}

func (cli client) ActivateCamera() (Camera, error) {
	for _, cam := range cli.cameras {
		if !cam.Active() {
			cam.Activate()

			return cam, nil
		}
	}
	return nil, errors.New("all cameras occupied")
}

func (cli *client) load(scene Scene, spriteCache warehouse.Cache[Sprite], soundCache warehouse.Cache[Sound]) error {
	if !scene.TryStartLoading() {
		return nil
	}
	defer func() {
		scene.SetLoading(false)
	}()
	sto := scene.Storage()
	cursor := warehouse.Factory.NewCursor(blueprint.Queries.SpriteBundle, sto)
	for cursor.Next() {
		bundle := blueprintclient.Components.SpriteBundle.GetFromCursor(cursor)
		err := cli.spriteLoader.Load(bundle, spriteCache)
		if err != nil {
			return err
		}
	}
	cursor = warehouse.Factory.NewCursor(blueprint.Queries.SoundBundle, sto)
	for cursor.Next() {
		bundle := blueprintclient.Components.SoundBundle.GetFromCursor(cursor)
		err := cli.soundLoader.Load(bundle, soundCache)
		if err != nil {
			return err
		}
	}
	scene.SetLoaded(true)
	return nil
}

var defaultLoadingScene = func() *scene {
	ls := &scene{}
	ls.name = "default loading scene"
	schema := table.Factory.NewSchema()
	ls.storage = warehouse.Factory.NewStorage(schema)
	ls.systems.renderers = append(ls.systems.renderers, defaultLoaderTextSystem{"Loading!"})
	return ls
}()

type defaultLoaderTextSystem struct {
	LoadingText string
}

func (sys defaultLoaderTextSystem) Render(scene Scene, screen Screen, cameraUtil CameraUtility) {
	loadingText := sys.LoadingText
	if loadingText == "" {
		loadingText = "Loading!"
	}
	for _, cam := range cameraUtil.ActiveCamerasFor(scene) {
		if cameraUtil.Ready(cam) {
			continue
		}
		textFace := text.NewGoXFace(basicfont.Face7x13)
		textBoundsX, textBoundsY := text.Measure(loadingText, textFace, 0)
		width, height := cam.Dimensions()
		centerX := float64((width - int(textBoundsX)) / 2)
		centerY := float64((height - int(textBoundsY)) / 2)
		cam.DrawTextBasicStatic(loadingText, &text.DrawOptions{}, textFace, vector.Two{
			X: centerX,
			Y: centerY + textBoundsY,
		})
		cam.PresentToScreen(screen)
	}
}
