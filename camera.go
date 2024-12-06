package coldbrew

import (
	"github.com/TheBitDrifter/blueprint/vector"
	"github.com/hajimehoshi/ebiten/v2"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
)

var _ Camera = &camera{}

// Camera manages viewport rendering and coordinate transformation between
// screen space and scene space
type Camera interface {
	// Ready checks if the camera is ready for rendering based on client state
	Ready(Client) bool
	// Active returns whether the camera is currently active
	Active() bool
	// Activate enables the camera
	Activate()
	// Deactivate clears the surface and disables the camera
	Deactivate()
	// Surface returns the camera's rendering surface
	Surface() *ebiten.Image
	// Localize transforms global scene coordinates to camera-local coordinates
	Localize(pos vector.Two) vector.Two
	// DrawImage renders an image at the specified scene position, localized to the camera
	DrawImage(img *ebiten.Image, opts *ebiten.DrawImageOptions, pos vector.Two)
	// DrawImageStatic renders an image at an absolute position without transformation (on the camera)
	DrawImageStatic(img *ebiten.Image, opts *ebiten.DrawImageOptions, pos vector.Two)
	// DrawTextBasic renders basic text at the specified scene position, localized to the camera
	DrawTextBasic(text string, opts *text.DrawOptions, fontFace *text.GoXFace, pos vector.Two)
	// DrawText renders text at the specified scene position, localized to the camera
	DrawText(text string, opts *text.DrawOptions, fontFace *text.GoTextFace, pos vector.Two)
	// DrawTextBasicStatic renders basic text at an absolute position (on the camera)
	DrawTextBasicStatic(text string, opts *text.DrawOptions, fontFace *text.GoXFace, pos vector.Two)
	// DrawTextStatic renders text at an absolute position (on the camera)
	DrawTextStatic(text string, opts *text.DrawOptions, fontFace *text.GoTextFace, pos vector.Two)
	// PresentToScreen renders the camera's contents to the provided screen
	PresentToScreen(screen Screen)
	// SetDimensions updates the camera's width and height
	SetDimensions(width, height int)
	// Dimensions returns the camera's current width and height
	Dimensions() (width, height int)
	// Positions returns the camera's global (screen) and local (scene/world) positions
	Positions() (screen, scene *vector.Two)
	// Index returns the camera's rendering priority
	Index() int
	// setIndex sets the camera's rendering priority
	setIndex(int)
}

type camera struct {
	active                        bool
	surface                       Sprite
	height, width                 int
	screenPosition, worldPosition vector.Two
	index                         int
}

func (c *camera) Ready(cli Client) bool {
	entry, ok := cli.CameraSceneTracker()[c]
	cameraLastChanged := entry.Tick
	if !ok {
		return false
	}
	cutoff := 0
	if ClientConfig.enforceMinOnActive {
		cutoff = cameraLastChanged
	} else {
		cutoff = entry.Scene.LastActivatedTick()
	}
	return tick-cutoff >= ClientConfig.minimumLoadTime && c.active
}

func (c *camera) Active() bool {
	return c.active
}

func (c *camera) Activate() {
	c.active = true
}

func (c *camera) Deactivate() {
	c.Surface().Clear()
	c.active = false
}

func (c *camera) Surface() *ebiten.Image { return c.surface.Image() }

// DrawImage draws an image at the given scene position, applying camera transformation
func (c *camera) DrawImage(img *ebiten.Image, opts *ebiten.DrawImageOptions, pos vector.Two) {
	localPos := c.Localize(pos)
	opts.GeoM.Translate(localPos.X, localPos.Y)
	c.surface.Image().DrawImage(img, opts)
}

// DrawImageStatic draws an image at the given position without camera transformation
func (c *camera) DrawImageStatic(img *ebiten.Image, opts *ebiten.DrawImageOptions, pos vector.Two) {
	opts.GeoM.Translate(pos.X, pos.Y)
	c.surface.Image().DrawImage(img, opts)
}

// DrawText draws text using the advanced text renderer with camera transformation
func (c *camera) DrawText(content string, opts *text.DrawOptions, fontFace *text.GoTextFace, pos vector.Two) {
	localPos := c.Localize(pos)
	opts.GeoM.Translate(localPos.X, localPos.Y)
	text.Draw(c.Surface(), content, fontFace, opts)
}

// DrawTextBasic draws text using the basic text renderer with camera transformation
func (c *camera) DrawTextBasic(content string, opts *text.DrawOptions, fontFace *text.GoXFace, pos vector.Two) {
	localPos := c.Localize(pos)
	opts.GeoM.Translate(localPos.X, localPos.Y)
	text.Draw(c.Surface(), content, fontFace, opts)
}

// DrawTextStatic draws text using the advanced text renderer without camera transformation
func (c *camera) DrawTextStatic(content string, opts *text.DrawOptions, fontFace *text.GoTextFace, pos vector.Two) {
	opts.GeoM.Translate(pos.X, pos.Y)
	text.Draw(c.Surface(), content, fontFace, opts)
}

// DrawTextBasicStatic draws text using the basic text renderer without camera transformation
func (c *camera) DrawTextBasicStatic(content string, opts *text.DrawOptions, fontFace *text.GoXFace, pos vector.Two) {
	opts.GeoM.Translate(pos.X, pos.Y)
	text.Draw(c.Surface(), content, fontFace, opts)
}

// PresentToScreen draws the camera's surface to the target screen at the camera's position
func (c *camera) PresentToScreen(screen Screen) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(c.screenPosition.X, c.screenPosition.Y)
	screen.Image().DrawImage(c.surface.Image(), opts)
}

// SetDimensions updates the camera dimensions and recreates the surface with the new size
func (c *camera) SetDimensions(width, height int) {
	c.width = width
	c.height = height

	newImg := ebiten.NewImage(width, height)
	c.surface = &sprite{name: c.surface.Name(), image: newImg}
}

func (c *camera) Dimensions() (width, height int) { return c.width, c.height }

func (c *camera) Positions() (*vector.Two, *vector.Two) {
	return &c.screenPosition, &c.worldPosition
}

// SetPositions updates the camera's screen and world positions
func (c *camera) SetPositions(screen, world vector.Two) {
	c.screenPosition = screen
	c.worldPosition = world
}

// Localize converts world coordinates to camera-local coordinates
func (c *camera) Localize(pos vector.Two) vector.Two {
	return vector.Two{X: -c.worldPosition.X + pos.X, Y: -c.worldPosition.Y + pos.Y}
}

func (c *camera) Index() int {
	return c.index
}

func (c *camera) setIndex(index int) {
	c.index = index
}
