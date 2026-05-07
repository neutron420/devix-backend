package bookmark

import (
	"context"
	"devix-backend/internal/models"
	"devix-backend/internal/pkg/pagination"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Bookmark(ctx context.Context, userID, postID uuid.UUID) error {
	bookmark := &models.Bookmark{
		ID:     uuid.New(),
		UserID: userID,
		PostID: postID,
	}
	return r.db.WithContext(ctx).Create(bookmark).Error
}

func (r *Repository) Unbookmark(ctx context.Context, userID, postID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Delete(&models.Bookmark{}).Error
}

func (r *Repository) IsBookmarked(ctx context.Context, userID, postID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Bookmark{}).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) List(ctx context.Context, userID uuid.UUID, cursor string, limit int) ([]models.Post, bool, error) {
	var posts []models.Post
	
	db := r.db.WithContext(ctx).
		Table("posts").
		Select("posts.*").
		Joins("JOIN bookmarks ON bookmarks.post_id = posts.id").
		Where("bookmarks.user_id = ?", userID)

	if cursor != "" {
		if decoded, err := pagination.DecodeCursor(cursor); err == nil && decoded != nil {
			db = db.Where("(posts.created_at, posts.id) < (?, ?)", decoded.CreatedAt, decoded.ID)
		}
	}

	err := db.Order("bookmarks.created_at desc").
		Limit(limit + 1).
		Find(&posts).Error
		
	if err != nil {
		return nil, false, err
	}

	hasMore := len(posts) > limit
	if hasMore {
		posts = posts[:limit]
	}

	return posts, hasMore, nil
}
