package websocket

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return false
		}

		u, err := url.Parse(origin)
		if err != nil {
			return false
		}

		host := u.Hostname()
		return host == "localhost" || host == "127.0.0.1" || host == "peteveryone.com" || host == "www.peteveryone.com"
	},
}

func ServeWs(hub *Hub, userID string, isGuest bool, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), UserID: userID, IsGuest: isGuest}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
