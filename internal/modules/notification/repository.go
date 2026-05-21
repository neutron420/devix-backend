package notification

import (
	"context"
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

func (r *Repository) Create(ctx context.Context, notification *models.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *Repository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
	if err != nil {
		return nil, err
	}

	if len(notifications) == 0 {
		return notifications, nil
	}

	actorIDs := make([]uuid.UUID, 0, len(notifications))
	seen := make(map[uuid.UUID]bool)
	for _, n := range notifications {
		if !seen[n.ActorID] {
			seen[n.ActorID] = true
			actorIDs = append(actorIDs, n.ActorID)
		}
	}

	var actors []models.User
	if err := r.db.WithContext(ctx).Where("id IN ?", actorIDs).Find(&actors).Error; err != nil {
		return nil, err
	}

	actorMap := make(map[uuid.UUID]*models.User)
	for i := range actors {
		actorMap[actors[i].ID] = &actors[i]
	}

	for i := range notifications {
		if actor, ok := actorMap[notifications[i].ActorID]; ok {
			notifications[i].Actor = &models.UserPublicProfile{
				ID:          actor.ID,
				Username:    actor.Username,
				DisplayName: actor.DisplayName,
				AvatarURL:   actor.AvatarURL,
			}
		}
	}

	return notifications, nil
}

func (r *Repository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (r *Repository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *Repository) MarkAsRead(ctx context.Context, userID, notificationID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true).Error
}

func (r *Repository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ?", userID).
		Update("is_read", true).Error
}
