package websocket

import (
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// Service provides methods to trigger websocket events from other modules.
type Service struct {
	hub *Hub
	log zerolog.Logger
}

func NewService(hub *Hub, log zerolog.Logger) *Service {
	return &Service{
		hub: hub,
		log: log.With().Str("module", "websocket_service").Logger(),
	}
}

// BroadcastEvent sends an event to all connected clients.
func (s *Service) BroadcastEvent(ctx context.Context, eventType string, payload interface{}) {
	s.log.Debug().Str("type", eventType).Msg("broadcasting websocket event")
	s.hub.Broadcast(Message{
		Type:    eventType,
		Payload: payload,
	})
}

// NotifyUser sends an event to a specific user.
func (s *Service) NotifyUser(ctx context.Context, userID uuid.UUID, eventType string, payload interface{}) {
	s.log.Debug().Str("type", eventType).Str("user_id", userID.String()).Msg("notifying user via websocket")
	s.hub.NotifyUser(userID, Message{
		Type:    eventType,
		Payload: payload,
	})
}
