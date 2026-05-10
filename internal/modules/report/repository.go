package report

import (
	"context"
	"errors"

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

func (r *Repository) Create(ctx context.Context, report *models.Report) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.Report, error) {
	var report models.Report
	err := r.db.WithContext(ctx).First(&report, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &report, err
}

func (r *Repository) ListPending(ctx context.Context, limit, offset int) ([]models.Report, int64, error) {
	var reports []models.Report
	var total int64
	db := r.db.WithContext(ctx).Model(&models.Report{}).Where("status = ?", "pending")
	db.Count(&total)
	err := db.Order("created_at ASC").Limit(limit).Offset(offset).Find(&reports).Error
	return reports, total, err
}

func (r *Repository) ListAll(ctx context.Context, limit, offset int) ([]models.Report, int64, error) {
	var reports []models.Report
	var total int64
	db := r.db.WithContext(ctx).Model(&models.Report{})
	db.Count(&total)
	err := db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&reports).Error
	return reports, total, err
}

func (r *Repository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, reviewedBy uuid.UUID, reviewNote string) error {
	return r.db.WithContext(ctx).Model(&models.Report{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      status,
		"reviewed_by": reviewedBy,
		"review_note": reviewNote,
	}).Error
}

func (r *Repository) HasUserReported(ctx context.Context, reporterID uuid.UUID, targetType string, targetID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Report{}).
		Where("reporter_id = ? AND target_type = ? AND target_id = ?", reporterID, targetType, targetID).
		Count(&count).Error
	return count > 0, err
}
