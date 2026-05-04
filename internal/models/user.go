package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents the system user model.
type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username     string         `gorm:"uniqueIndex;not null;size:30"`
	Email        string         `gorm:"uniqueIndex;not null;size:255"`
	PasswordHash string         `gorm:"not null"`
	DisplayName  *string        `gorm:"size:100"`
	Bio          *string        `gorm:"type:text"`
	AvatarURL    *string        `gorm:"type:text"`
	Role         string         `gorm:"default:user;size:20"`
	IsVerified   bool           `gorm:"default:false"`
	PostCount    int            `gorm:"default:0"`
	Reputation   int            `gorm:"default:0"`
	LastLoginAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// UserPublicProfile is what other users can see.
type UserPublicProfile struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	DisplayName *string   `json:"display_name"`
	AvatarURL   *string   `json:"avatar_url"`
}
