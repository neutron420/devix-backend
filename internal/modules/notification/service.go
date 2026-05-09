package notification

import (
	"context"
	"devix-backend/internal/models"
	"devix-backend/internal/modules/websocket"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	repo      *Repository
	wsService *websocket.Service
	log       zerolog.Logger
}

func NewService(repo *Repository, wsService *websocket.Service, log zerolog.Logger) *Service {
	return &Service{
		repo:      repo,
		wsService: wsService,
		log:       log.With().Str("module", "notification").Logger(),
	}
}

func (s *Service) GetUserNotifications(ctx context.Context, userID uuid.UUID, page, limit int) (*NotificationListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}
	offset := (page - 1) * limit

	notifications, err := s.repo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	unreadCount, _ := s.repo.GetUnreadCount(ctx, userID)

	responses := make([]NotificationResponse, 0, len(notifications))
	for _, n := range notifications {
		responses = append(responses, s.toResponse(&n))
	}

	return &NotificationListResponse{
		Notifications: responses,
		UnreadCount:   unreadCount,
	}, nil
}

func (s *Service) CreateNotification(ctx context.Context, userID, actorID, targetID uuid.UUID, action string) error {
	// Don't notify if user is the actor
	if userID == actorID {
		return nil
	}

	notification := &models.Notification{
		ID:        uuid.New(),
		UserID:    userID,
		ActorID:   actorID,
		TargetID:  targetID,
		Action:    action,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, notification); err != nil {
		return err
	}

	// Trigger real-time update
	if s.wsService != nil {
		s.wsService.NotifyUser(ctx, userID, "new_notification", s.toResponse(notification))
	}

	return nil
}

func (s *Service) MarkAsRead(ctx context.Context, userID, notificationID uuid.UUID) error {
	return s.repo.MarkAsRead(ctx, userID, notificationID)
}

func (s *Service) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}

func (s *Service) toResponse(n *models.Notification) NotificationResponse {
	resp := NotificationResponse{
		ID:        n.ID.String(),
		Action:    n.Action,
		TargetID:  n.TargetID.String(),
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt,
	}

	if n.Actor != nil {
		resp.Actor = ActorResponse{
			ID:        n.Actor.ID.String(),
			Username:  n.Actor.Username,
			AvatarURL: n.Actor.AvatarURL,
		}
	}

	return resp
}
