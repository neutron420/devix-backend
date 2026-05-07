package notification

import (
	"time"
)

type NotificationResponse struct {
	ID        string        `json:"id"`
	Actor     ActorResponse `json:"actor"`
	Action    string        `json:"action"`
	TargetID  string        `json:"target_id"`
	IsRead    bool          `json:"is_read"`
	CreatedAt time.Time     `json:"created_at"`
}

type ActorResponse struct {
	ID        string  `json:"id"`
	Username  string  `json:"username"`
	AvatarURL *string `json:"avatar_url"`
}

type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	UnreadCount   int64                 `json:"unread_count"`
}
