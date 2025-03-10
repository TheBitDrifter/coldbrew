package coldbrew

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/TheBitDrifter/bark"
	"github.com/TheBitDrifter/blueprint"
	"github.com/TheBitDrifter/table"
	"github.com/TheBitDrifter/warehouse"
)

var _ SceneManager = &sceneManager{}

// SceneManager handles scene lifecycle, transitions, and state management
type SceneManager interface {
	// Scene access

	// IsActive checks if the given scene is currently active
	IsActive(Scene) bool
	// LoadingScenes returns all scenes currently being loaded
	LoadingScenes() []Scene
	// Cache returns the scene cache
	Cache() warehouse.Cache[Scene]

	// Scene lifecycle
	ActiveScenes() []Scene

	// RegisterScene creates and registers a new scene with the provided configuration
	RegisterScene(string, int, int, blueprint.Plan, []RenderSystem, []ClientSystem, []blueprint.CoreSystem) error
	// ChangeScene transitions to a target scene, transferring specified entities
	ChangeScene(target Scene, entities ...warehouse.Entity) error
	// ActivateScene activates a target scene while keeping the origin scene active
	ActivateScene(target Scene, entities ...warehouse.Entity) error
	// DeactivateScene removes a scene from the active scenes list
	DeactivateScene(target Scene)
}

type sceneManager struct {
	mu            sync.RWMutex
	activeScenes  []Scene
	loadingScenes []Scene
	cache         warehouse.Cache[Scene]
	cacheBust     atomic.Bool
}

// newSceneManager creates a scene manager with specified cache size
func newSceneManager(maxScenesCached int) *sceneManager {
	return &sceneManager{
		cache: warehouse.FactoryNewCache[Scene](maxScenesCached),
	}
}

func (m *sceneManager) ActiveScenes() []Scene {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create a copy to avoid race conditions on the returned slice
	result := make([]Scene, len(m.activeScenes))
	copy(result, m.activeScenes)
	return result
}

// Scene access methods
func (m *sceneManager) IsActive(check Scene) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, match := range m.activeScenes {
		if check == match {
			return true
		}
	}
	return false
}

func (m *sceneManager) LoadingScenes() []Scene {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create a copy to avoid race conditions on the returned slice
	result := make([]Scene, len(m.loadingScenes))
	copy(result, m.loadingScenes)
	return result
}

func (m *sceneManager) Cache() warehouse.Cache[Scene] {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.cache
}

// ActivateScene adds a target scene to active scenes and transfers specified entities from origin
func (m *sceneManager) ActivateScene(target Scene, entities ...warehouse.Entity) error {
	targetStorage := target.Storage()
	for _, en := range entities {
		originStorage := en.Storage()
		if err := originStorage.TransferEntities(targetStorage, en); err != nil {
			return bark.AddTrace(err)
		}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, scene := range m.activeScenes {
		if scene == target {
			target.SetSelectedTick()
			return nil
		}
	}
	target.SetActivatedTick()
	target.SetSelectedTick()
	m.activeScenes = append(m.activeScenes, target)
	return nil
}

// ChangeScene replaces the current active scene with the target scene
// Only works when exactly one scene is active
func (m *sceneManager) ChangeScene(target Scene, entities ...warehouse.Entity) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.activeScenes) > 1 {
		return bark.AddTrace(
			errors.New("cannot use change scene api when multiple scenes are active — use activate scene api instead"),
		)
	}
	if len(m.activeScenes) == 0 {
		return bark.AddTrace(errors.New("no scenes are active"))
	}
	origin := m.activeScenes[0]
	originStorage := origin.Storage()
	targetStorage := target.Storage()
	for _, en := range entities {
		if err := originStorage.TransferEntities(targetStorage, en); err != nil {
			return bark.AddTrace(err)
		}
	}
	target.SetSelectedTick()
	m.activeScenes[0] = target
	return nil
}

// DeactivateScene removes the target scene from the active scenes list
func (m *sceneManager) DeactivateScene(target Scene) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, scene := range m.activeScenes {
		if scene == target {
			lastIdx := len(m.activeScenes) - 1
			m.activeScenes[i] = m.activeScenes[lastIdx]
			m.activeScenes = m.activeScenes[:lastIdx]
			break
		}
	}
}

// RegisterScene creates and registers a new scene with the provided configuration
// If no scenes are active, the new scene becomes active
func (m *sceneManager) RegisterScene(
	name string,
	width, height int,
	plan blueprint.Plan,
	renderSystems []RenderSystem,
	clientSystems []ClientSystem,
	coreSystems []blueprint.CoreSystem,
) error {
	newScene, err := m.newScene(name, width, height, plan, renderSystems, clientSystems, coreSystems)
	if err != nil {
		return bark.AddTrace(err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.cache.Register(name, newScene)
	if len(m.activeScenes) == 0 {
		m.activeScenes = append(m.activeScenes, newScene)
	}
	return nil
}

// newScene creates a new scene with the provided configuration
func (m *sceneManager) newScene(
	name string,
	width, height int,
	plan blueprint.Plan,
	renderSystems []RenderSystem,
	clientSystems []ClientSystem,
	coreSystems []blueprint.CoreSystem,
) (Scene, error) {
	schema := table.Factory.NewSchema()
	storage := warehouse.Factory.NewStorage(schema)
	newScene := &scene{
		plan:    plan,
		storage: storage,
		name:    name,
		height:  height,
		width:   width,
		systems: struct {
			renderers []RenderSystem
			client    []ClientSystem
			core      []blueprint.CoreSystem
		}{
			renderers: renderSystems,
			core:      coreSystems,
			client:    clientSystems,
		},
	}
	return newScene, nil
}
