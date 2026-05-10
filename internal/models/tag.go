package models

import (
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string     `gorm:"size:50;uniqueIndex;not null"`
	Slug        string     `gorm:"size:60;uniqueIndex;not null"`
	Description *string    `gorm:"type:text"`
	ParentID    *uuid.UUID `gorm:"type:uuid;index"`
	Category    string     `gorm:"size:50;index;default:general"`
	Synonyms    string     `gorm:"type:text"`
	PostCount   int        `gorm:"default:0"`
	CreatedAt   time.Time
}
