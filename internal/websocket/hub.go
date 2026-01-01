package websocket

import (
	"context"
	"log"
	"pet-everyone/internal/transport"
)

type Hub struct {
	id        string // Hub ID = associated petID
	clients   map[*Client]bool
	broadcast chan []byte

	register   chan *Client
	unregister chan *Client

	commands chan<- transport.Envelope
	cancel   context.CancelFunc
}

func NewHub(id string, cmds chan<- transport.Envelope, cancel context.CancelFunc) *Hub {
	return &Hub{
		id:         id,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		commands:   cmds,
		cancel:     cancel,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("[HUB %s]: REGISTERED NEW CLIENT {%s}", h.id, client.UserID)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("[HUB %s]: DEREGISTERED CLIENT", h.id)
				if len(h.clients) == 0 {
					h.cancel()
				}
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

func (h *Hub) Clean() {
	for client := range h.clients {
		client.send <- []byte("shutdown")
		h.unregister <- client
	}
	log.Printf("[HUB %s]: shutting down", h.id)
}

func (h *Hub) GetBroadcastChannel() chan []byte {
	return h.broadcast
}
