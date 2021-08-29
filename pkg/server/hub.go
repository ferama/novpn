package server

import "github.com/gorilla/websocket"

type Hub struct {
	// Registered clients.
	clients map[*websocket.Conn]bool
	// Register requests from the clients.
	register chan *websocket.Conn

	// Unregister requests from clients.
	unregister chan *websocket.Conn
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			delete(h.clients, client)
		}
	}
}
