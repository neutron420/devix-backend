package comment

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

func (r *Repository) Create(ctx context.Context, c *models.Comment) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.Comment, error) {
	var c models.Comment
	err := r.db.WithContext(ctx).First(&c, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &c, err
}

func (r *Repository) GetByPostID(ctx context.Context, postID uuid.UUID) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.WithContext(ctx).
		Where("post_id = ?", postID).
		Order("created_at asc").
		Find(&comments).Error

	if err != nil {
		return nil, err
	}

	for i := range comments {
		var author models.User
		if err := r.db.WithContext(ctx).First(&author, "id = ?", comments[i].AuthorID).Error; err == nil {
			comments[i].Author = &models.UserPublicProfile{
				ID:          author.ID,
				Username:    author.Username,
				DisplayName: author.DisplayName,
				AvatarURL:   author.AvatarURL,
			}
		}
	}

	return comments, nil
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, content string) error {
	return r.db.WithContext(ctx).Model(&models.Comment{}).Where("id = ?", id).Update("content", content).Error
}

func (r *Repository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Comment{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"content":    "[deleted]",
		}).Error
}

func (r *Repository) IncrementPostCommentCount(ctx context.Context, postID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Post{}).Where("id = ?", postID).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error
}

func (r *Repository) GetParentDepth(ctx context.Context, parentID uuid.UUID) (int, error) {
	var c models.Comment
	err := r.db.WithContext(ctx).Select("depth").First(&c, "id = ?", parentID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return -1, nil
	}
	return c.Depth, err
}
