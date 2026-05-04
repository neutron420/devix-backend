package models

import (
	"time"

	"github.com/google/uuid"
)

// Vote represents a user's upvote/downvote on a post or comment.
type Vote struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_post;uniqueIndex:idx_user_comment"`
	PostID    *uuid.UUID `gorm:"type:uuid;index;uniqueIndex:idx_user_post"`
	CommentID *uuid.UUID `gorm:"type:uuid;index;uniqueIndex:idx_user_comment"`
	VoteType  int        `gorm:"type:smallint;not null"` // 1 for upvote, -1 for downvote
	CreatedAt time.Time
}
