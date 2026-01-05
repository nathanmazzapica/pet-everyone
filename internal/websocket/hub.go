package websocket

import (
	"context"
	"log"
	"pet-everyone/internal/transport"
	"time"
)

type Hub struct {
	id        string // Hub ID = associated petID
	clients   map[*Client]bool
	broadcast chan []byte

	register   chan *Client
	unregister chan *Client

	commands chan<- transport.Envelope
	cancel   context.CancelFunc

	shutdownTimer *time.Timer
	shutdownDelay time.Duration
}

func NewHub(id string, cmds chan<- transport.Envelope, cancel context.CancelFunc, shutdownDelay time.Duration) *Hub {
	return &Hub{
		id:            id,
		broadcast:     make(chan []byte),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		clients:       make(map[*Client]bool),
		commands:      cmds,
		cancel:        cancel,
		shutdownDelay: shutdownDelay,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("[HUB %s]: REGISTERED NEW CLIENT {%s}", h.id, client.UserID)

			// Cancel shutdown if a client reconnects
			if h.shutdownTimer != nil {
				if h.shutdownTimer.Stop() {
					log.Printf("[HUB %s]: CANCELLED SHUTDOWN - client reconnected", h.id)
				}
				h.shutdownTimer = nil
			}

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("[HUB %s]: DEREGISTERED CLIENT", h.id)

				if len(h.clients) == 0 {
					// Start shutdown timer instead of immediate shutdown
					log.Printf("[HUB %s]: NO CLIENTS - scheduling shutdown in %v", h.id, h.shutdownDelay)
					h.shutdownTimer = time.AfterFunc(h.shutdownDelay, func() {
						log.Printf("[HUB %s]: SHUTDOWN TIMER EXPIRED - triggering shutdown", h.id)
						h.cancel()
					})
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
	// Stop shutdown timer if it's running
	if h.shutdownTimer != nil {
		h.shutdownTimer.Stop()
		h.shutdownTimer = nil
	}

	for client := range h.clients {
		client.send <- []byte("shutdown")
		h.unregister <- client
	}
	log.Printf("[HUB %s]: shutting down", h.id)
}

func (h *Hub) GetBroadcastChannel() chan []byte {
	return h.broadcast
}

func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

func (h *Hub) BroadcastExcept(message []byte, exceptUserID string) {
	for client := range h.clients {
		if client.UserID == exceptUserID {
			continue
		}
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

func (h *Hub) SendToUser(message []byte, userID string) {
	for client := range h.clients {
		if client.UserID == userID {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
			}
			return
		}
	}
	log.Printf("[HUB %s]: user %s not found for targeted message", h.id, userID)
}
