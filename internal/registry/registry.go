package registry

import (
	"context"
	"log"
	"pet-everyone/internal/data/model"
	"pet-everyone/internal/service"
	"pet-everyone/internal/transport"
	"pet-everyone/internal/websocket"
	"sync"
)

type HubRegistry struct {
	mu       sync.RWMutex
	Hubs     map[string]*websocket.Hub
	petModel *model.PetModel
}

func NewHubRegistry(petModel *model.PetModel) *HubRegistry {
	return &HubRegistry{
		mu:       sync.RWMutex{},
		Hubs:     make(map[string]*websocket.Hub),
		petModel: petModel,
	}
}

// GetOrCreateHub retrieves an existing Hub by ID or creates and initializes a new one if it doesn't exist.
// Returns the Hub and a boolean indicating whether it was newly created.
func (h *HubRegistry) GetOrCreateHub(id string) (*websocket.Hub, bool) {
	if hub, ok := h.Hubs[id]; ok {
		return hub, false
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	newHub := h.initializeHub(id)
	h.Hubs[id] = newHub
	go h.Hubs[id].Run()

	log.Println("Created hub for pet{", id, "}")
	return h.Hubs[id], true
}

func (h *HubRegistry) initializeHub(id string) *websocket.Hub {
	ctx, cancel := context.WithCancel(context.Background())
	serializer := transport.NewJSONSerializer()

	petService := service.NewPetService(ctx, id, h.petModel, serializer.In())
	chatService := service.NewChatService(serializer.In())

	router := transport.NewRouter(petService, chatService)
	go router.Route()

	hub := websocket.NewHub(id, router.In(), cancel)
	serializer.Subscribe(hub.GetBroadcastChannel())

	go func() {
		for {
			select {
			case <-ctx.Done():
				h.RemoveHub(id)
			}
		}
	}()

	return hub
}

func (h *HubRegistry) RemoveHub(id string) {
	if hub, ok := h.Hubs[id]; ok {
		// clean up hub
		hub.Clean()

		// delete hub
		h.mu.Lock()
		delete(h.Hubs, id)
		h.mu.Unlock()
	}
}
