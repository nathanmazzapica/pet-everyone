package registry

import (
	"context"
	"log"
	"pet-everyone/internal/data/model"
	"pet-everyone/internal/service"
	"pet-everyone/internal/transport"
	"pet-everyone/internal/websocket"
	"sync"
	"time"
)

type HubRegistry struct {
	mu           sync.RWMutex
	Hubs         map[string]*websocket.Hub
	petModel     *model.PetModel
	userModel    *model.UserModel
	visitorModel *model.VisitorModel
}

func NewHubRegistry(petModel *model.PetModel) *HubRegistry {
	return &HubRegistry{
		mu:       sync.RWMutex{},
		Hubs:     make(map[string]*websocket.Hub),
		petModel: petModel,
	}
}

// SetResolverDeps injects user and visitor models for display-name resolution.
func (h *HubRegistry) SetResolverDeps(userModel *model.UserModel, visitorModel *model.VisitorModel) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.userModel = userModel
	h.visitorModel = visitorModel
}

// GetOrCreateHub retrieves an existing Hub by ID or creates and initializes a new one if it doesn't exist.
// Returns the Hub and a boolean indicating whether it was newly created.
func (h *HubRegistry) GetOrCreateHub(id string) (*websocket.Hub, bool) {
	// First check with read lock (fast path for existing hubs)
	h.mu.RLock()
	if hub, ok := h.Hubs[id]; ok {
		h.mu.RUnlock()
		return hub, false
	}
	h.mu.RUnlock()

	// Not found, acquire write lock and check again (double-check pattern)
	h.mu.Lock()
	defer h.mu.Unlock()

	// Check again in case another goroutine created it while we were waiting for the lock
	if hub, ok := h.Hubs[id]; ok {
		return hub, false
	}

	// Still doesn't exist, create it
	newHub := h.initializeHub(id)
	h.Hubs[id] = newHub
	go newHub.Run()

	log.Println("Created hub for pet{", id, "}")
	return newHub, true
}

func (h *HubRegistry) initializeHub(id string) *websocket.Hub {
	ctx, cancel := context.WithCancel(context.Background())

	if h.userModel == nil || h.visitorModel == nil {
		log.Panic("hub registry resolver dependencies not set")
	}

	petService := service.NewPetService(ctx, id, h.petModel)
	resolver := service.NewDisplayNameResolver(h.userModel, h.visitorModel, service.DisplayNameResolverConfig{})
	chatService := service.NewChatService(resolver)

	router := transport.NewRouter(petService, chatService)

	go router.Route(ctx)

	// 30 second delay before shutting down after last client disconnects
	// This prevents unnecessary hub recreation on page refreshes
	const shutdownDelay = 30 * time.Second
	hub := websocket.NewHub(id, router.In(), cancel, shutdownDelay)

	serializer := transport.NewJSONSerializer(hub)

	// Wire Router events to Serializer
	go func() {
		for {
			select {
			case event, ok := <-router.Events():
				if !ok {
					return
				}
				serializer.In() <- event
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		<-ctx.Done()
		h.RemoveHub(id)
	}()

	return hub
}

func (h *HubRegistry) RemoveHub(id string) {
	// Lock while accessing the map
	h.mu.Lock()
	hub, ok := h.Hubs[id]
	if ok {
		// Remove from registry first to prevent concurrent access
		delete(h.Hubs, id)
	}
	h.mu.Unlock()

	// Clean up the hub after releasing the lock
	if ok {
		hub.Clean()
		log.Printf("[REGISTRY]: Removed hub for pet{%s}", id)
	}
}
