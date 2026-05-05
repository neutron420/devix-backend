package vote

type VoteRequest struct {
	VoteType int `json:"vote_type" binding:"required,vote_type"`
}

type VoteResponse struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	VoteType int    `json:"vote_type"`
	NewCount int    `json:"new_count"`
}
