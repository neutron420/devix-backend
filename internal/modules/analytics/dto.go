package analytics

type PostAnalyticsResponse struct {
	TotalViews int64            `json:"total_views"`
	Countries  map[string]int64 `json:"countries"`
	Devices    map[string]int64 `json:"devices"`
	Browsers   map[string]int64 `json:"browsers"`
	OS         map[string]int64 `json:"os"`
	Referrers  map[string]int64 `json:"referrers"`
}
