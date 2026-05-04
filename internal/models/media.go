package models

import (
	"time"

	"github.com/google/uuid"
)

type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
)

// Media represents an attachment to a post.
type Media struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PostID       uuid.UUID `gorm:"type:uuid;not null;index"`
	FileURL      string    `gorm:"type:text;not null"`
	FileType     MediaType `gorm:"size:10;not null"`
	FileSize     int64     `gorm:"not null"`
	MimeType     string    `gorm:"size:100;not null"`
	OriginalName *string   `gorm:"size:255"`
	Width        *int
	Height       *int
	SortOrder    int       `gorm:"default:0"`
	CreatedAt    time.Time
}
