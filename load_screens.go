package coldbrew

type AnimationConfig struct {
	FrameWidth,
	FrameHeight,
	FrameCount,
	Speed int // Frames per tick.
}

type AnimationState struct {
	StartTick int
}
