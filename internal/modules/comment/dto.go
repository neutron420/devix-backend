package comment

type CreateCommentRequest struct {
	Content  string  `json:"content" binding:"required,min=1,max=5000"`
	ParentID *string `json:"parent_id" binding:"omitempty,uuid"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=5000"`
}

type CommentResponse struct {
	ID        string            `json:"id"`
	PostID    string            `json:"post_id"`
	Author    *AuthorResponse   `json:"author"`
	ParentID  *string           `json:"parent_id"`
	Content   string            `json:"content"`
	Depth     int               `json:"depth"`
	VoteCount int               `json:"vote_count"`
	IsDeleted bool              `json:"is_deleted"`
	Replies   []CommentResponse `json:"replies,omitempty"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}

type AuthorResponse struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	DisplayName *string `json:"display_name"`
	AvatarURL   *string `json:"avatar_url"`
}
