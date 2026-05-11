package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
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

func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *Repository) GetUserByVerificationToken(ctx context.Context, token string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("verification_token = ?", token).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *Repository) SetUserVerified(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"is_verified":        true,
			"verification_token": nil,
		}).Error
}

func (r *Repository) SetUserResetToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"reset_token":      &token,
			"reset_expires_at": &expiresAt,
		}).Error
}

func (r *Repository) GetUserByResetToken(ctx context.Context, token string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("reset_token = ? AND reset_expires_at > ?", token, time.Now()).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *Repository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"password_hash":    passwordHash,
			"reset_token":      nil,
			"reset_expires_at": nil,
		}).Error
}

func (r *Repository) StoreRefreshToken(ctx context.Context, userID uuid.UUID, hash string, expiresAt time.Time) error {
	token := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: hash,
		ExpiresAt: expiresAt,
	}
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *Repository) GetRefreshToken(ctx context.Context, hash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := r.db.WithContext(ctx).Where("token_hash = ? AND revoked = ?", hash, false).First(&token).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &token, err
}

func (r *Repository) RevokeRefreshToken(ctx context.Context, hash string) error {
	return r.db.WithContext(ctx).Model(&models.RefreshToken{}).
		Where("token_hash = ?", hash).
		Update("revoked", true).Error
}

func (r *Repository) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("revoked", true).Error
}

func (r *Repository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("last_login_at", &now).Error
}

func HashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}
