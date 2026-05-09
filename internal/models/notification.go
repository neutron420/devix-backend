package models

import (
	"time"
	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null;index:idx_notifications_user_read_created"` 
	ActorID   uuid.UUID          `gorm:"type:uuid;not null"`
	Actor     *UserPublicProfile `gorm:"-"`
	Action    string             `gorm:"size:50;not null"`
	TargetID  uuid.UUID `gorm:"type:uuid;not null"`       
	IsRead    bool      `gorm:"default:false;index;index:idx_notifications_user_read_created"`
	CreatedAt time.Time `gorm:"index:idx_notifications_user_read_created"`
}
