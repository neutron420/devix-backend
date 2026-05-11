package poll

import "time"

type CreatePollRequest struct {
	PostID    string    `json:"post_id" binding:"required"`
	Question  string    `json:"question" binding:"required,max=255"`
	Options   []string  `json:"options" binding:"required,min=2,max=10"`
	ExpiresAt time.Time `json:"expires_at" binding:"required"`
}

type PollResponse struct {
	ID        string         `json:"id"`
	Question  string         `json:"question"`
	Options   []OptionResp   `json:"options"`
	ExpiresAt time.Time      `json:"expires_at"`
	TotalVotes int           `json:"total_votes"`
	HasVoted  bool           `json:"has_voted"`
	VotedID   string         `json:"voted_option_id,omitempty"`
}

type OptionResp struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Votes int    `json:"votes"`
}

type VoteRequest struct {
	OptionID string `json:"option_id" binding:"required"`
}
