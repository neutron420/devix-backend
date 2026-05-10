package report

type CreateReportRequest struct {
	TargetType  string `json:"target_type" binding:"required,oneof=post comment user"`
	TargetID    string `json:"target_id" binding:"required,uuid"`
	Reason      string `json:"reason" binding:"required,oneof=spam harassment misinformation inappropriate other"`
	Description string `json:"description" binding:"omitempty,max=1000"`
}

type ReviewReportRequest struct {
	Status     string `json:"status" binding:"required,oneof=reviewed resolved rejected"`
	ReviewNote string `json:"review_note" binding:"omitempty,max=1000"`
}

type ReportResponse struct {
	ID          string  `json:"id"`
	ReporterID  string  `json:"reporter_id"`
	TargetType  string  `json:"target_type"`
	TargetID    string  `json:"target_id"`
	Reason      string  `json:"reason"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	ReviewedBy  *string `json:"reviewed_by,omitempty"`
	ReviewNote  string  `json:"review_note,omitempty"`
	CreatedAt   string  `json:"created_at"`
}
