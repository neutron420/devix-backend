package websocket

import (
	"encoding/json"
	"sync"

	"github.com/google/uuid"
)

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type Hub struct {
	clients map[*Client]bool

	broadcast chan []byte

	register chan *Client

	unregister chan *Client

	userClients map[uuid.UUID][]*Client

	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:   make(chan []byte),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		userClients: make(map[uuid.UUID][]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.userClients[client.UserID] = append(h.userClients[client.UserID], client)
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				cls := h.userClients[client.UserID]
				for i, c := range cls {
					if c == client {
						h.userClients[client.UserID] = append(cls[:i], cls[i+1:]...)
						break
					}
				}
				if len(h.userClients[client.UserID]) == 0 {
					delete(h.userClients, client.UserID)
				}
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:

				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) Broadcast(msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	h.broadcast <- data
}

func (h *Hub) NotifyUser(userID uuid.UUID, msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	h.mu.RLock()
	clients, ok := h.userClients[userID]
	h.mu.RUnlock()

	if ok {
		for _, client := range clients {
			select {
			case client.send <- data:
			default:
			}
		}
	}
}
