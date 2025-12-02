package registry

import (
	"log"
	"pet-everyone/internal/db/models"
	"pet-everyone/internal/service"
	"pet-everyone/internal/transport"
	"pet-everyone/internal/websocket"
	"sync"
)

type HubRegistry struct {
	serializer transport.Serializer
	mu         sync.RWMutex
	Hubs       map[string]*websocket.Hub
	petModel   *models.PetModel
}

func NewHubRegistry(petModel *models.PetModel) *HubRegistry {
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
	serializer := transport.NewJSONSerializer()

	petService := service.NewPetService(id, h.petModel, serializer.In())
	chatService := service.NewChatService(serializer.In())

	router := transport.NewRouter(petService, chatService)
	go router.Route()

	hub := websocket.NewHub(id, router.In())
	serializer.Subscribe(hub.GetBroadcastChannel())

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
