package models

import (
	"time"
	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null"` 
	ActorID   uuid.UUID          `gorm:"type:uuid;not null"`
	Actor     *UserPublicProfile `gorm:"-"`
	Action    string             `gorm:"size:50;not null"`
	TargetID  uuid.UUID `gorm:"type:uuid;not null"`       
	IsRead    bool      `gorm:"default:false;index"`
	CreatedAt time.Time
}
