package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostType string

const (
	PostTypeQuestion PostType = "question"
	PostTypeConcept  PostType = "concept"
	PostTypeBuildLog PostType = "build-log"
)

type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusArchived  PostStatus = "archived"
	PostStatusDeleted   PostStatus = "deleted"
)

// Post represents a knowledge-sharing post or question.
type Post struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	AuthorID     uuid.UUID      `gorm:"type:uuid;not null;index"`
	Author       *UserPublicProfile `gorm:"-"` // Not persisted in posts table
	Title        string         `gorm:"size:300;not null"`
	Slug         string         `gorm:"uniqueIndex;not null;size:350"`
	Content      string         `gorm:"type:text;not null"`
	PostType     PostType       `gorm:"size:20;not null;index"`
	Status       PostStatus     `gorm:"size:20;default:published;index"`
	ViewCount    int            `gorm:"default:0"`
	VoteCount    int            `gorm:"default:0"`
	CommentCount int            `gorm:"default:0"`
	IsPinned     bool           `gorm:"default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	// Associations
	Tags  []Tag   `gorm:"many2many:post_tags;"`
	Media []Media `gorm:"foreignKey:PostID"`
}
