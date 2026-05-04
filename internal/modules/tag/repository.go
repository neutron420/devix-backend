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
