package websocket

import (
	"log"
	"sync"
)

type HubRegistry struct {
	mu   sync.RWMutex
	Hubs map[string]*Hub
}

func NewHubRegistry() *HubRegistry {
	return &HubRegistry{
		mu:   sync.RWMutex{},
		Hubs: make(map[string]*Hub),
	}
}

func (h *HubRegistry) GetOrCreateHub(id string) *Hub {
	if hub, ok := h.Hubs[id]; ok {
		return hub
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	h.Hubs[id] = NewHub(id)
	go h.Hubs[id].Run()

	log.Println("Created hub for pet{", id, "}")
	return h.Hubs[id]
}

func (h *HubRegistry) RemoveHub(id string) {
	if _, ok := h.Hubs[id]; ok {
		// clean up hub

		// delete hub
		h.mu.Lock()
		delete(h.Hubs, id)
		h.mu.Unlock()
	}
}
