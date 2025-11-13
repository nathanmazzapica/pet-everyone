package websocket

import "log"

type Hub struct {
	id        string
	clients   map[*Client]bool
	broadcast chan []byte

	register   chan *Client
	unregister chan *Client
}

func NewHub(id string) *Hub {
	return &Hub{
		id:         id,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("[HUB %s]: REGISTERED NEW CLIENT", h.id)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("[HUB %s]: DEREGISTERED CLIENT", h.id)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
