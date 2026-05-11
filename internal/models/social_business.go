package models

import (
	"time"

	"github.com/google/uuid"
)


type Conversation struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	User1ID   uuid.UUID `gorm:"type:uuid;not null;index"`
	User2ID   uuid.UUID `gorm:"type:uuid;not null;index"`
	LastMsgAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Message struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ConversationID uuid.UUID `gorm:"type:uuid;not null;index"`
	SenderID       uuid.UUID `gorm:"type:uuid;not null;index"`
	Content        string    `gorm:"type:text;not null"`
	IsRead         bool      `gorm:"default:false"`
	CreatedAt      time.Time
}

type Organization struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `gorm:"size:100;not null;uniqueIndex"`
	Slug        string    `gorm:"size:100;not null;uniqueIndex"`
	Bio         string    `gorm:"type:text"`
	AvatarURL   string    `gorm:"type:text"`
	WebsiteURL  string    `gorm:"size:255"`
	OwnerID     uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OrgMember struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrgID   uuid.UUID `gorm:"type:uuid;not null;index"`
	UserID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Role    string    `gorm:"size:20;default:'member'"` 
	JoinedAt time.Time
}

type Poll struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PostID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	Question  string    `gorm:"size:255;not null"`
	ExpiresAt time.Time
	CreatedAt time.Time
}

type PollOption struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PollID uuid.UUID `gorm:"type:uuid;not null;index"`
	Text   string    `gorm:"size:255;not null"`
	Votes  int       `gorm:"default:0"`
}

type PollVote struct {
	PollID   uuid.UUID `gorm:"type:uuid;primaryKey;not null"`
	OptionID uuid.UUID `gorm:"type:uuid;primaryKey;not null"`
	UserID   uuid.UUID `gorm:"type:uuid;primaryKey;not null"`
}


type AnalyticsEvent struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TargetID  uuid.UUID `gorm:"type:uuid;not null;index"` // PostID or ProfileID
	Type      string    `gorm:"size:50;index"`           // view, click, etc.
	Country   string    `gorm:"size:100"`
	Device    string    `gorm:"size:100"`
	OS        string    `gorm:"size:100"`
	Browser   string    `gorm:"size:100"`
	Referrer  string    `gorm:"size:255"`
	IPHash    string    `gorm:"size:64;index"`
	CreatedAt time.Time
}
