package post

type CreatePostRequest struct {
	Title    string   `json:"title" binding:"required,min=5,max=300"`
	Content  string   `json:"content" binding:"required,min=10"`
	PostType string   `json:"post_type" binding:"required,post_type"`
	Status   string   `json:"status" binding:"omitempty,oneof=draft published"`
	Tags     []string `json:"tags" binding:"omitempty,max=10,dive,min=1,max=50"`
}

type UpdatePostRequest struct {
	Title   *string  `json:"title" binding:"omitempty,min=5,max=300"`
	Content *string  `json:"content" binding:"omitempty,min=10"`
	Status  *string  `json:"status" binding:"omitempty,oneof=draft published archived"`
	Tags    []string `json:"tags" binding:"omitempty,max=10,dive,min=1,max=50"`
}

type FeedQuery struct {
	Cursor   string `form:"cursor"`
	Limit    int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Type     string `form:"type" binding:"omitempty,post_type"`
	Tag      string `form:"tag"`
	AuthorID string `form:"author_id"`
	Search   string `form:"q"`
	Sort     string `form:"sort" binding:"omitempty,oneof=latest trending"`
}

type PostResponse struct {
	ID           string          `json:"id"`
	Author       *AuthorResponse `json:"author"`
	Title        string          `json:"title"`
	Slug         string          `json:"slug"`
	Content      string          `json:"content"`
	PostType     string          `json:"post_type"`
	Status       string          `json:"status"`
	ViewCount    int             `json:"view_count"`
	VoteCount    int             `json:"vote_count"`
	CommentCount int             `json:"comment_count"`
	IsPinned     bool            `json:"is_pinned"`
	Tags         []TagResponse   `json:"tags"`
	Media        []MediaResponse `json:"media"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
}

type AuthorResponse struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	DisplayName *string `json:"display_name"`
	AvatarURL   *string `json:"avatar_url"`
}

type TagResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type MediaResponse struct {
	ID           string `json:"id"`
	FileURL      string `json:"file_url"`
	FileType     string `json:"file_type"`
	FileSize     int64  `json:"file_size"`
	MimeType     string `json:"mime_type"`
	OriginalName string `json:"original_name,omitempty"`
}

type PostListResponse struct {
	Posts   []PostResponse `json:"posts"`
	Cursor  string         `json:"cursor,omitempty"`
	HasMore bool           `json:"has_more"`
}
