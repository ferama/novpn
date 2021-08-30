package server

import (
	"github.com/gorilla/websocket"
)

type hubItem struct {
	key  string
	conn *websocket.Conn
}

type Hub struct {
	// Registered clients.
	clients map[string]*websocket.Conn
	// Register requests from the clients.
	register chan *hubItem

	// Unregister requests from clients.
	unregister chan *hubItem
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[string]*websocket.Conn),
		register:   make(chan *hubItem),
		unregister: make(chan *hubItem),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			// log.Println("register ", client.key)
			h.clients[client.key] = client.conn
		case client := <-h.unregister:
			// log.Println("unregister ", client.key)
			delete(h.clients, client.key)
		}
	}
}
