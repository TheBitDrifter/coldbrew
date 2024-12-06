package coldbrew

import "github.com/TheBitDrifter/warehouse"

type scene struct {
	index         int
	loaded        bool
	height, width int
	storage       warehouse.Storage
	systems       struct {
		renderers []RenderSystem
	}
}
