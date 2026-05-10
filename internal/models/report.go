package models

import (
	"time"

	"github.com/google/uuid"
)

type ReportStatus string

const (
	ReportStatusPending  ReportStatus = "pending"
	ReportStatusReviewed ReportStatus = "reviewed"
	ReportStatusResolved ReportStatus = "resolved"
	ReportStatusRejected ReportStatus = "rejected"
)

type Report struct {
	ID           uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ReporterID   uuid.UUID    `gorm:"type:uuid;not null;index"`
	TargetType   string       `gorm:"size:20;not null;index"`
	TargetID     uuid.UUID    `gorm:"type:uuid;not null;index"`
	Reason       string       `gorm:"size:50;not null"`
	Description  string       `gorm:"type:text"`
	Status       ReportStatus `gorm:"size:20;default:pending;index"`
	ReviewedBy   *uuid.UUID   `gorm:"type:uuid"`
	ReviewNote   string       `gorm:"type:text"`
	CreatedAt    time.Time    `gorm:"index"`
	UpdatedAt    time.Time
}

type ActivityLog struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Action     string    `gorm:"size:50;not null;index"`
	TargetType string    `gorm:"size:20;not null"`
	TargetID   uuid.UUID `gorm:"type:uuid;not null"`
	Metadata   string    `gorm:"type:text"`
	CreatedAt  time.Time `gorm:"index"`
}
