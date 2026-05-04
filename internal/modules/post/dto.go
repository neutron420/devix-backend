package post

// --- Request DTOs ---

// CreatePostRequest represents the post creation request.
type CreatePostRequest struct {
	Title    string   `json:"title" binding:"required,min=5,max=300"`
	Content  string   `json:"content" binding:"required,min=10"`
	PostType string   `json:"post_type" binding:"required,post_type"`
	Status   string   `json:"status" binding:"omitempty,oneof=draft published"`
	Tags     []string `json:"tags" binding:"omitempty,max=10,dive,min=1,max=50"`
}

// UpdatePostRequest represents the post update request.
type UpdatePostRequest struct {
	Title   *string  `json:"title" binding:"omitempty,min=5,max=300"`
	Content *string  `json:"content" binding:"omitempty,min=10"`
	Status  *string  `json:"status" binding:"omitempty,oneof=draft published archived"`
	Tags    []string `json:"tags" binding:"omitempty,max=10,dive,min=1,max=50"`
}

// FeedQuery represents feed/list query parameters.
type FeedQuery struct {
	Cursor   string `form:"cursor"`
	Limit    int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Type     string `form:"type" binding:"omitempty,post_type"`
	Tag      string `form:"tag"`
	AuthorID string `form:"author_id"`
	Search   string `form:"q"`
	Sort     string `form:"sort" binding:"omitempty,oneof=latest trending"`
}

// --- Response DTOs ---

// PostResponse is the standard post response.
type PostResponse struct {
	ID           string              `json:"id"`
	Author       *AuthorResponse     `json:"author"`
	Title        string              `json:"title"`
	Slug         string              `json:"slug"`
	Content      string              `json:"content"`
	PostType     string              `json:"post_type"`
	Status       string              `json:"status"`
	ViewCount    int                 `json:"view_count"`
	VoteCount    int                 `json:"vote_count"`
	CommentCount int                 `json:"comment_count"`
	IsPinned     bool                `json:"is_pinned"`
	Tags         []TagResponse       `json:"tags"`
	Media        []MediaResponse     `json:"media"`
	CreatedAt    string              `json:"created_at"`
	UpdatedAt    string              `json:"updated_at"`
}

// AuthorResponse is the embedded author in post responses.
type AuthorResponse struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	DisplayName *string `json:"display_name"`
	AvatarURL   *string `json:"avatar_url"`
}

// TagResponse is the embedded tag in post responses.
type TagResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// MediaResponse is the embedded media in post responses.
type MediaResponse struct {
	ID           string `json:"id"`
	FileURL      string `json:"file_url"`
	FileType     string `json:"file_type"`
	FileSize     int64  `json:"file_size"`
	MimeType     string `json:"mime_type"`
	OriginalName string `json:"original_name,omitempty"`
}

// PostListResponse is the paginated list of posts.
type PostListResponse struct {
	Posts   []PostResponse `json:"posts"`
	Cursor string         `json:"cursor,omitempty"`
	HasMore bool          `json:"has_more"`
}
