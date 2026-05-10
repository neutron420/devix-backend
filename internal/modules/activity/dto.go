package activity

import "time"

type ActivityResponse struct {
	ID         string `json:"id"`
	Action     string `json:"action"`
	TargetType string `json:"target_type"`
	TargetID   string `json:"target_id"`
	Metadata   string `json:"metadata,omitempty"`
	CreatedAt  string `json:"created_at"`
}

type ActivityListResponse struct {
	Activities []ActivityResponse `json:"activities"`
	Total      int64              `json:"total"`
}

type ActivityQuery struct {
	Limit  int       `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset int       `form:"offset" binding:"omitempty,min=0"`
	Action string    `form:"action" binding:"omitempty"`
	Since  time.Time `form:"since" binding:"omitempty"`
}
