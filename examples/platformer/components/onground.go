package components

import (
	"github.com/TheBitDrifter/blueprint/vector"
	"github.com/TheBitDrifter/warehouse"
)

var OnGroundComponent = warehouse.FactoryNewComponent[OnGround]()

type OnGround struct {
	LastTouch   int
	Landed      int
	LastJump    int
	SlopeNormal vector.Two
}
