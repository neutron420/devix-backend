package tag

import (
	"context"
	"errors"
	"time"

	"devix-backend/internal/models"
	"devix-backend/internal/pkg/sanitize"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(ctx context.Context) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.WithContext(ctx).Order("post_count desc, name asc").Find(&tags).Error
	return tags, err
}

func (r *Repository) GetBySlug(ctx context.Context, slug string) (*models.Tag, error) {
	var t models.Tag
	err := r.db.WithContext(ctx).First(&t, "slug = ?", slug).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &t, err
}

func (r *Repository) GetByName(ctx context.Context, name string) (*models.Tag, error) {
	var t models.Tag
	err := r.db.WithContext(ctx).Where("LOWER(name) = LOWER(?)", name).First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &t, err
}

func (r *Repository) Create(ctx context.Context, name string, description *string) (*models.Tag, error) {
	slug := sanitize.Slug(name)
	t := &models.Tag{
		ID:          uuid.New(),
		Name:        name,
		Slug:        slug,
		Description: description,
		CreatedAt:   time.Now(),
	}
	err := r.db.WithContext(ctx).Create(t).Error
	return t, err
}

func (r *Repository) GetTrending(ctx context.Context, limit int, since time.Duration) ([]models.Tag, error) {
	var tags []models.Tag
	sinceTime := time.Now().Add(-since)

	err := r.db.WithContext(ctx).
		Table("tags").
		Select("tags.*, COUNT(post_tags.post_id) as recent_count").
		Joins("JOIN post_tags ON post_tags.tag_id = tags.id").
		Joins("JOIN posts ON posts.id = post_tags.post_id").
		Where("posts.created_at >= ?", sinceTime).
		Group("tags.id").
		Order("recent_count DESC").
		Limit(limit).
		Find(&tags).Error

	return tags, err
}

func (r *Repository) GetChildren(ctx context.Context, parentID uuid.UUID) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Order("name ASC").Find(&tags).Error
	return tags, err
}

func (r *Repository) GetByCategory(ctx context.Context, category string) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.WithContext(ctx).Where("category = ?", category).Order("post_count DESC, name ASC").Find(&tags).Error
	return tags, err
}

func (r *Repository) GetTopLevel(ctx context.Context) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.WithContext(ctx).Where("parent_id IS NULL").Order("post_count DESC, name ASC").Find(&tags).Error
	return tags, err
}
