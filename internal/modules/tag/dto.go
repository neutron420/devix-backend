package tag

type CreateTagRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=50"`
	Description *string `json:"description" binding:"omitempty,max=500"`
}

type TagResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	PostCount   int     `json:"post_count"`
}
