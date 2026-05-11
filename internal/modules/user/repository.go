package user

import (
	"context"
	"errors"

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

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "username = ?", username).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *Repository) UsernameExists(ctx context.Context, username string, excludeID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("username = ? AND id != ?", username, excludeID).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *Repository) UpdateProfile(ctx context.Context, userID uuid.UUID, req *UpdateProfileRequest) error {
	updates := make(map[string]interface{})
	if req.DisplayName != nil {
		updates["display_name"] = *req.DisplayName
	}
	if req.Bio != nil {
		updates["bio"] = *req.Bio
	}
	if req.Username != nil {
		updates["username"] = *req.Username
	}
	if req.WebsiteURL != nil {
		updates["website_url"] = *req.WebsiteURL
	}
	if req.GitHubURL != nil {
		updates["github_url"] = *req.GitHubURL
	}
	if req.TwitterURL != nil {
		updates["twitter_url"] = *req.TwitterURL
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.Preferences != nil {
		updates["preferences"] = *req.Preferences
	}
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}

func (r *Repository) Delete(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, userID).Error
}

func (r *Repository) UpdateStatus(ctx context.Context, userID uuid.UUID, isActive bool) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("is_active", isActive).Error
}

func (r *Repository) UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarURL string) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("avatar_url", avatarURL).Error
}

func (r *Repository) AdjustReputation(ctx context.Context, userID uuid.UUID, points int) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("reputation", gorm.Expr("reputation + ?", points)).Error
}
