package websocket

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

const (
	RedisChannelWS = "devix:ws:events"
)

type Hub struct {
	clients map[*Client]bool

	broadcast chan []byte

	register chan *Client

	unregister chan *Client

	userClients map[uuid.UUID][]*Client
	roomClients map[string][]*Client

	redis *redis.Client
	log   zerolog.Logger

	mu sync.RWMutex
}

func NewHub(redisClient *redis.Client, log zerolog.Logger) *Hub {
	return &Hub{
		broadcast:   make(chan []byte),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		userClients: make(map[uuid.UUID][]*Client),
		roomClients: make(map[string][]*Client),
		redis:       redisClient,
		log:         log.With().Str("component", "websocket_hub").Logger(),
	}
}

func (h *Hub) Run(ctx context.Context) {
	if h.redis != nil {
		go h.listenToRedis(ctx)
	}

	for {
		select {
		case <-ctx.Done():
			return
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

				// Remove from all rooms
				for roomName, cls := range h.roomClients {
					for i, c := range cls {
						if c == client {
							h.roomClients[roomName] = append(cls[:i], cls[i+1:]...)
							break
						}
					}
					if len(h.roomClients[roomName]) == 0 {
						delete(h.roomClients, roomName)
					}
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

type RedisWSMessage struct {
	TargetUserID *uuid.UUID `json:"target_user_id,omitempty"`
	RoomName     *string    `json:"room_name,omitempty"`
	Data         []byte     `json:"data"`
}

func (h *Hub) listenToRedis(ctx context.Context) {
	pubsub := h.redis.Subscribe(ctx, RedisChannelWS)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		var redisMsg RedisWSMessage
		if err := json.Unmarshal([]byte(msg.Payload), &redisMsg); err != nil {
			continue
		}

		if redisMsg.TargetUserID != nil {
			h.localNotifyUser(*redisMsg.TargetUserID, redisMsg.Data)
		} else if redisMsg.RoomName != nil {
			h.localNotifyRoom(*redisMsg.RoomName, redisMsg.Data)
		} else {
			h.localBroadcast(redisMsg.Data)
		}
	}
}

func (h *Hub) localBroadcast(data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		select {
		case client.send <- data:
		default:
		}
	}
}

func (h *Hub) localNotifyRoom(roomName string, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	clients, ok := h.roomClients[roomName]
	if ok {
		for _, client := range clients {
			select {
			case client.send <- data:
			default:
			}
		}
	}
}

func (h *Hub) localNotifyUser(userID uuid.UUID, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	clients, ok := h.userClients[userID]
	if ok {
		for _, client := range clients {
			select {
			case client.send <- data:
			default:
			}
		}
	}
}

func (h *Hub) JoinRoom(client *Client, roomName string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.roomClients[roomName] = append(h.roomClients[roomName], client)
}

func (h *Hub) LeaveRoom(client *Client, roomName string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	cls, ok := h.roomClients[roomName]
	if ok {
		for i, c := range cls {
			if c == client {
				h.roomClients[roomName] = append(cls[:i], cls[i+1:]...)
				break
			}
		}
		if len(h.roomClients[roomName]) == 0 {
			delete(h.roomClients, roomName)
		}
	}
}

func (h *Hub) NotifyRoom(roomName string, msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	if h.redis != nil {
		redisMsg := RedisWSMessage{RoomName: &roomName, Data: data}
		payload, _ := json.Marshal(redisMsg)
		h.redis.Publish(context.Background(), RedisChannelWS, payload)
	} else {
		h.localNotifyRoom(roomName, data)
	}
}

func (h *Hub) Broadcast(msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	if h.redis != nil {
		redisMsg := RedisWSMessage{Data: data}
		payload, _ := json.Marshal(redisMsg)
		h.redis.Publish(context.Background(), RedisChannelWS, payload)
	} else {
		h.broadcast <- data
	}
}

func (h *Hub) NotifyUser(userID uuid.UUID, msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	if h.redis != nil {
		redisMsg := RedisWSMessage{TargetUserID: &userID, Data: data}
		payload, _ := json.Marshal(redisMsg)
		h.redis.Publish(context.Background(), RedisChannelWS, payload)
	} else {
		h.localNotifyUser(userID, data)
	}
}
