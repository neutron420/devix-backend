package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username     string    `gorm:"uniqueIndex;not null;size:30"`
	Email        string    `gorm:"uniqueIndex;not null;size:255"`
	PasswordHash string    `gorm:"not null"`
	DisplayName  *string   `gorm:"size:100"`
	Bio          *string   `gorm:"type:text"`
	AvatarURL    *string   `gorm:"type:text"`
	Role         string    `gorm:"default:user;size:20"`
	IsActive     bool      `gorm:"default:true"`
	IsVerified   bool      `gorm:"default:false"`
	PostCount    int       `gorm:"default:0"`
	Reputation   int       `gorm:"default:0"`
	LastLoginAt  *time.Time
	WebsiteURL   string `gorm:"size:255"`
	GitHubURL    string `gorm:"size:255"`
	TwitterURL   string `gorm:"size:255"`
	Location     string `gorm:"size:100"`
	Preferences  string `gorm:"type:text;default:'{}'"` // JSON string
	VerificationToken *string `gorm:"size:255"`
	ResetToken        *string `gorm:"size:255"`
	ResetExpiresAt   *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type UserPublicProfile struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	DisplayName *string   `json:"display_name"`
	AvatarURL   *string   `json:"avatar_url"`
}
