package models

import (
	"time"
	"github.com/google/uuid"
)

type Bookmark struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_post_bookmark"`
	PostID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_post_bookmark"`
	CreatedAt time.Time
	Post *Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}
