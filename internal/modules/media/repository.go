package media

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

func (r *Repository) Create(ctx context.Context, media *models.Media) error {
	return r.db.WithContext(ctx).Create(media).Error
}

func (r *Repository) GetByPostID(ctx context.Context, postID uuid.UUID) ([]models.Media, error) {
	var media []models.Media
	err := r.db.WithContext(ctx).Where("post_id = ?", postID).Order("sort_order asc, created_at asc").Find(&media).Error
	return media, err
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Media{}, "id = ?", id).Error
}

func (r *Repository) DeleteByPostID(ctx context.Context, postID uuid.UUID) ([]string, error) {
	var media []models.Media
	if err := r.db.WithContext(ctx).Where("post_id = ?", postID).Find(&media).Error; err != nil {
		return nil, err
	}

	urls := make([]string, len(media))
	for i, m := range media {
		urls[i] = m.FileURL
	}

	err := r.db.WithContext(ctx).Where("post_id = ?", postID).Delete(&models.Media{}).Error
	return urls, err
}

func (r *Repository) CountByPostAndType(ctx context.Context, postID uuid.UUID, fileType models.MediaType) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Media{}).Where("post_id = ? AND file_type = ?", postID, fileType).Count(&count).Error
	return count, err
}

func (r *Repository) GetNextSortOrder(ctx context.Context, postID uuid.UUID) (int, error) {
	var maxOrder int
	err := r.db.WithContext(ctx).Model(&models.Media{}).Where("post_id = ?", postID).Select("COALESCE(MAX(sort_order), -1)").Scan(&maxOrder).Error
	return maxOrder + 1, err
}
