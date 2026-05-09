package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID        uuid.UUID          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PostID    uuid.UUID          `gorm:"type:uuid;not null;index;index:idx_comments_post_created"`
	AuthorID  uuid.UUID          `gorm:"type:uuid;not null;index;index:idx_comments_author_created"`
	Author    *UserPublicProfile `gorm:"-"`
	ParentID  *uuid.UUID         `gorm:"type:uuid;index"`
	Content   string             `gorm:"type:text;not null"`
	Depth     int                `gorm:"default:0"`
	VoteCount int                `gorm:"default:0"`
	IsDeleted bool               `gorm:"default:false"`
	CreatedAt time.Time          `gorm:"index;index:idx_comments_post_created;index:idx_comments_author_created"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
