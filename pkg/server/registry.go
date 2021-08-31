package server

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

// Registry is an in memory structure to hold websocket connection
type Registry struct {
	data map[string]*websocket.Conn

	mu sync.Mutex
}

// NewRegistry creates a new registry
func NewRegistry() *Registry {
	return &Registry{
		data: make(map[string]*websocket.Conn),
	}
}

// Add adds an item to the registry in a thread safe way
func (r *Registry) Add(key string, t *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[key] = t
}

// GetByID returns an item give its registry key
func (r *Registry) GetByKey(key string) (*websocket.Conn, error) {
	if val, ok := r.data[key]; ok {
		return val, nil
	}
	return nil, errors.New("item not found")
}

// Delete removes an item from registry
func (r *Registry) Delete(key string) error {
	if _, err := r.GetByKey(key); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.data, key)
	return nil
}
