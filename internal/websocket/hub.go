package websocket

import (
	"context"
	"log"
	"pet-everyone/internal/transport"
	"sync"
	"time"
)

type Hub struct {
	id        string // Hub ID = associated petID
	clients   map[*Client]bool
	mu        sync.RWMutex
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
		mu:            sync.RWMutex{},
		commands:      cmds,
		cancel:        cancel,
		shutdownDelay: shutdownDelay,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

			log.Printf("[HUB %s]: REGISTERED NEW CLIENT {%s}", h.id, client.UserID)

			// Cancel shutdown if a client reconnects
			if h.shutdownTimer != nil {
				if h.shutdownTimer.Stop() {
					log.Printf("[HUB %s]: CANCELLED SHUTDOWN - client reconnected", h.id)
				}
				h.shutdownTimer = nil
			}

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				if err := client.conn.Close(); err != nil {
					log.Printf("[HUB %s]: error closing client connection: %v", h.id, err)
				}
				close(client.send)
				delete(h.clients, client)
				log.Printf("[HUB %s]: DEREGISTERED CLIENT", h.id)

			}

			if len(h.clients) == 0 {
				// Start shutdown timer instead of immediate shutdown
				log.Printf("[HUB %s]: NO CLIENTS - scheduling shutdown in %v", h.id, h.shutdownDelay)
				h.shutdownTimer = time.AfterFunc(h.shutdownDelay, func() {
					log.Printf("[HUB %s]: SHUTDOWN TIMER EXPIRED - triggering shutdown", h.id)
					h.cancel()
				})
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					select {
					case h.unregister <- client:
					default:
						log.Printf("[HUB %s]: failed to unregister client", h.id)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) Clean() {
	// Stop shutdown timer if it's running
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.shutdownTimer != nil {
		h.shutdownTimer.Stop()
		h.shutdownTimer = nil
	}

	for client := range h.clients {
		client.send <- []byte("shutdown")
		select {
		case h.unregister <- client:
		default:
			log.Printf("[HUB %s]: failed to unregister client", h.id)
		}
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
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.UserID == exceptUserID {
			continue
		}
		select {
		case client.send <- message:
		default:
			select {
			case h.unregister <- client:
			default:
				log.Printf("[HUB %s]: failed to unregister client", h.id)
			}
		}
	}
}

func (h *Hub) SendToUser(message []byte, userID string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.UserID == userID {
			select {
			case client.send <- message:
			default:
				select {
				case h.unregister <- client:
				default:
					log.Printf("[HUB %s]: failed to unregister client", h.id)
				}
			}
			return
		}
	}
	log.Printf("[HUB %s]: user %s not found for targeted message", h.id, userID)
}
