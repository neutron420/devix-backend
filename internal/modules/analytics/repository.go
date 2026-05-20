package analytics

import (
	"context"

	"devix-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) LogEvent(ctx context.Context, event *models.AnalyticsEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *Repository) GetPostStats(ctx context.Context, postID uuid.UUID) (*PostAnalyticsResponse, error) {
	var stats PostAnalyticsResponse
	stats.Countries = make(map[string]int64)
	stats.Devices = make(map[string]int64)
	stats.Browsers = make(map[string]int64)
	stats.OS = make(map[string]int64)
	stats.Referrers = make(map[string]int64)

	if err := r.db.WithContext(ctx).Model(&models.AnalyticsEvent{}).
		Where("target_id = ? AND type = ?", postID, "view").
		Count(&stats.TotalViews).Error; err != nil {
		return nil, err
	}

	if err := r.groupCounts(ctx, postID, "country", stats.Countries); err != nil {
		return nil, err
	}
	if err := r.groupCounts(ctx, postID, "device", stats.Devices); err != nil {
		return nil, err
	}
	if err := r.groupCounts(ctx, postID, "browser", stats.Browsers); err != nil {
		return nil, err
	}
	if err := r.groupCounts(ctx, postID, "os", stats.OS); err != nil {
		return nil, err
	}
	if err := r.groupCounts(ctx, postID, "referrer", stats.Referrers); err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *Repository) groupCounts(ctx context.Context, targetID uuid.UUID, column string, dest map[string]int64) error {
	type groupedCount struct {
		Key   string
		Count int64
	}

	var rows []groupedCount
	if err := r.db.WithContext(ctx).Model(&models.AnalyticsEvent{}).
		Select(column+" AS key, count(*) AS count").
		Where("target_id = ? AND type = ?", targetID, "view").
		Group(column).
		Scan(&rows).Error; err != nil {
		return err
	}

	for _, row := range rows {
		key := row.Key
		if key == "" {
			key = "unknown"
		}
		dest[key] = row.Count
	}

	return nil
}
