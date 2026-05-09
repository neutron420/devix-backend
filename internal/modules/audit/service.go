package audit

import (
	"context"
	"time"

	"devix-backend/internal/models"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Service struct {
	db  *gorm.DB
	log zerolog.Logger
}

func NewService(db *gorm.DB, log zerolog.Logger) *Service {
	return &Service{
		db:  db,
		log: log.With().Str("module", "audit").Logger(),
	}
}

func (s *Service) Log(ctx context.Context, actorID uuid.UUID, action, target, targetID, details, ip, userAgent string) {
	entry := &models.AuditLog{
		ID:        uuid.New(),
		ActorID:   actorID,
		Action:    action,
		Target:    target,
		TargetID:  targetID,
		Details:   details,
		IPAddress: ip,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
	}

	go func() {
		if err := s.db.Create(entry).Error; err != nil {
			s.log.Error().Err(err).Str("action", action).Msg("failed to write audit log")
		}
	}()
}

func (s *Service) GetLogs(ctx context.Context, limit, offset int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := s.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

func (s *Service) GetLogsByActor(ctx context.Context, actorID uuid.UUID, limit, offset int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := s.db.WithContext(ctx).
		Where("actor_id = ?", actorID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

func (s *Service) GetLogsByAction(ctx context.Context, action string, limit, offset int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := s.db.WithContext(ctx).
		Where("action = ?", action).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}
