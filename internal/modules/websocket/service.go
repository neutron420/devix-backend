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

func (s *Service) BroadcastTyping(userID, otherUserID uuid.UUID, isTyping bool) {
	s.hub.NotifyTyping(userID, otherUserID, isTyping)
}

func (s *Service) IsOnline(ctx context.Context, userID uuid.UUID) bool {
	return s.hub.IsUserOnline(ctx, userID)
}

func (s *Service) NotifyRoom(ctx context.Context, roomName string, eventType string, payload interface{}) {
	s.log.Debug().Str("type", eventType).Str("room", roomName).Msg("notifying room via websocket")
	s.hub.NotifyRoom(roomName, Message{
		Type:    eventType,
		Payload: payload,
	})
}

func (s *Service) JoinRoom(client *Client, roomName string) {
	s.hub.JoinRoom(client, roomName)
}

func (s *Service) LeaveRoom(client *Client, roomName string) {
	s.hub.LeaveRoom(client, roomName)
}
