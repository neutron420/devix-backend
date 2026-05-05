package websocket

import (
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

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

func (s *Service) BroadcastEvent(ctx context.Context, eventType string, payload interface{}) {
	s.log.Debug().Str("type", eventType).Msg("broadcasting websocket event")
	s.hub.Broadcast(Message{
		Type:    eventType,
		Payload: payload,
	})
}

func (s *Service) NotifyUser(ctx context.Context, userID uuid.UUID, eventType string, payload interface{}) {
	s.log.Debug().Str("type", eventType).Str("user_id", userID.String()).Msg("notifying user via websocket")
	s.hub.NotifyUser(userID, Message{
		Type:    eventType,
		Payload: payload,
	})
}
