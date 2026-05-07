package bookmark

import "devix-backend/internal/modules/post"

type BookmarkListResponse struct {
	Posts   []post.PostResponse `json:"posts"`
	Cursor  string             `json:"cursor,omitempty"`
	HasMore bool               `json:"has_more"`
}
