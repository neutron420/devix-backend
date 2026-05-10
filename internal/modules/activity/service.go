package activity

import (
	"context"
	"time"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/models"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	repo *Repository
	log  zerolog.Logger
}

func NewService(repo *Repository, log zerolog.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log.With().Str("module", "activity").Logger(),
	}
}

func (s *Service) Track(ctx context.Context, userID uuid.UUID, action, targetType string, targetID uuid.UUID, metadata string) {
	log := &models.ActivityLog{
		ID:         uuid.New(),
		UserID:     userID,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Metadata:   metadata,
		CreatedAt:  time.Now(),
	}
	if err := s.repo.Create(ctx, log); err != nil {
		s.log.Error().Err(err).Str("action", action).Msg("failed to track activity")
	}
}

func (s *Service) GetUserActivity(ctx context.Context, userID uuid.UUID, query ActivityQuery) (*ActivityListResponse, error) {
	if query.Limit <= 0 {
		query.Limit = 20
	}

	logs, total, err := s.repo.ListByUser(ctx, userID, query.Action, query.Since, query.Limit, query.Offset)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	activities := make([]ActivityResponse, 0, len(logs))
	for _, l := range logs {
		activities = append(activities, ActivityResponse{
			ID:         l.ID.String(),
			Action:     l.Action,
			TargetType: l.TargetType,
			TargetID:   l.TargetID.String(),
			Metadata:   l.Metadata,
			CreatedAt:  l.CreatedAt.Format(time.RFC3339),
		})
	}

	return &ActivityListResponse{Activities: activities, Total: total}, nil
}
