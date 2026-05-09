package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ActorID   uuid.UUID `gorm:"type:uuid;index"`
	Action    string    `gorm:"not null;size:50;index"`
	Target    string    `gorm:"size:100"`
	TargetID  string    `gorm:"size:100"`
	Details   string    `gorm:"type:text"`
	IPAddress string    `gorm:"size:45"`
	UserAgent string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"index"`
}
