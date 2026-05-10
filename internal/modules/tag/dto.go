package tag

type CreateTagRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=50"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	ParentID    *string `json:"parent_id" binding:"omitempty,uuid"`
	Category    string  `json:"category" binding:"omitempty,max=50"`
	Synonyms    string  `json:"synonyms" binding:"omitempty,max=500"`
}

type TagResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	ParentID    *string `json:"parent_id,omitempty"`
	Category    string  `json:"category"`
	Synonyms    string  `json:"synonyms,omitempty"`
	PostCount   int     `json:"post_count"`
}

type TagTreeResponse struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Slug        string             `json:"slug"`
	Description *string            `json:"description"`
	Category    string             `json:"category"`
	PostCount   int                `json:"post_count"`
	Children    []TagTreeResponse  `json:"children,omitempty"`
}
