package models

import (
	"time"

	"github.com/google/uuid"
)

type Vote struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_post;uniqueIndex:idx_user_comment"`
	PostID    *uuid.UUID `gorm:"type:uuid;index;uniqueIndex:idx_user_post"`
	CommentID *uuid.UUID `gorm:"type:uuid;index;uniqueIndex:idx_user_comment"`
	VoteType  int        `gorm:"type:smallint;not null"`
	CreatedAt time.Time
}
