package registry

import (
	"log"
	"pet-everyone/internal/websocket"
	"sync"
)

type HubRegistry struct {
	mu   sync.RWMutex
	Hubs map[string]*websocket.Hub
}

func NewHubRegistry() *HubRegistry {
	return &HubRegistry{
		mu:   sync.RWMutex{},
		Hubs: make(map[string]*websocket.Hub),
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

	h.Hubs[id] = websocket.NewHub(id)
	go h.Hubs[id].Run()

	log.Println("Created hub for pet{", id, "}")
	return h.Hubs[id], true
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
