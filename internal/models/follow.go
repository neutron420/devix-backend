package models

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	FollowerID  uuid.UUID `gorm:"type:uuid;primaryKey;index"`
	FollowingID uuid.UUID `gorm:"type:uuid;primaryKey;index"`
	CreatedAt   time.Time
}
