package coldbrew

// SystemManager manages registration of various system types
type SystemManager interface {
	RegisterGlobalRenderSystem(...GlobalRenderSystem)
	RegisterGlobalClientSystem(...GlobalClientSystem)
	RegisterGlobalClientSecondarySystem(...GlobalClientSystem)
}

type systemManager struct {
	globalRenderers              []GlobalRenderSystem
	globalClientSystems          []GlobalClientSystem
	globalClientSecondarySystems []GlobalClientSystem
}

// RegisterGlobalRenderSystem adds one or more render systems to the global renderers
func (sm *systemManager) RegisterGlobalRenderSystem(renderers ...GlobalRenderSystem) {
	for _, r := range renderers {
		sm.globalRenderers = append(sm.globalRenderers, r)
	}
}

// RegisterGlobalClientSystem adds one or more client systems to the global client systems
func (sm *systemManager) RegisterGlobalClientSystem(systems ...GlobalClientSystem) {
	for _, s := range systems {
		sm.globalClientSystems = append(sm.globalClientSystems, s)
	}
}

// RegisterGlobalClientSecondarySystem adds one or more client systems to run at the end of the update loop
func (sm *systemManager) RegisterGlobalClientSecondarySystem(systems ...GlobalClientSystem) {
	for _, s := range systems {
		sm.globalClientSecondarySystems = append(sm.globalClientSecondarySystems, s)
	}
}
