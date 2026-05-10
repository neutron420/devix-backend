package activity

import (
	"context"
	"time"

	"devix-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, log *models.ActivityLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *Repository) ListByUser(ctx context.Context, userID uuid.UUID, action string, since time.Time, limit, offset int) ([]models.ActivityLog, int64, error) {
	var logs []models.ActivityLog
	var total int64

	db := r.db.WithContext(ctx).Model(&models.ActivityLog{}).Where("user_id = ?", userID)

	if action != "" {
		db = db.Where("action = ?", action)
	}
	if !since.IsZero() {
		db = db.Where("created_at >= ?", since)
	}

	db.Count(&total)
	err := db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}
