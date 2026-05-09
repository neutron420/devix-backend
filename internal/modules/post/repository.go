package post

import (
	"context"
	"errors"
	"time"

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

func (r *Repository) Create(ctx context.Context, post *models.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *Repository) GetBySlug(ctx context.Context, slug string) (*models.Post, error) {
	var post models.Post
	err := r.db.WithContext(ctx).
		Preload("Tags").
		Preload("Media").
		Where("slug = ?", slug).
		First(&post).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var author models.User
	if err := r.db.WithContext(ctx).First(&author, "id = ?", post.AuthorID).Error; err == nil {
		post.Author = &models.UserPublicProfile{
			ID:          author.ID,
			Username:    author.Username,
			DisplayName: author.DisplayName,
			AvatarURL:   author.AvatarURL,
		}
	}

	return &post, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.Post, error) {
	var post models.Post
	err := r.db.WithContext(ctx).First(&post, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &post, err
}

func (r *Repository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Post, bool, error) {
	var posts []models.Post
	if len(ids) == 0 {
		return posts, false, nil
	}

	err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Order("created_at DESC").
		Find(&posts).Error

	return posts, false, err
}


func (r *Repository) List(ctx context.Context, query FeedQuery) ([]models.Post, bool, error) {
	limit := pagination.NormalizeLimit(query.Limit)
	var posts []models.Post

	db := r.db.WithContext(ctx).
		Table("posts").
		Select("posts.*").
		Where("posts.status = ?", "published")

	if query.Type != "" {
		db = db.Where("posts.post_type = ?", query.Type)
	}

	if query.AuthorID != "" {
		db = db.Where("posts.author_id = ?", query.AuthorID)
	}

	if len(query.AuthorIDs) > 0 {
		db = db.Where("posts.author_id IN ?", query.AuthorIDs)
	}

	if len(query.ExcludeAuthorIDs) > 0 {
		db = db.Where("posts.author_id NOT IN ?", query.ExcludeAuthorIDs)
	}

	if query.Search != "" {

		db = db.Where("to_tsvector('english', posts.title || ' ' || posts.content) @@ plainto_tsquery('english', ?)", query.Search)
	}

	if query.Tag != "" {
		db = db.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
			Joins("JOIN tags ON tags.id = post_tags.tag_id").
			Where("tags.slug = ?", query.Tag)
	}

	if query.Cursor != "" {
		if decoded, err := pagination.DecodeCursor(query.Cursor); err == nil && decoded != nil {
			db = db.Where("(posts.created_at, posts.id) < (?, ?)", decoded.CreatedAt, decoded.ID)
		}
	}

	if query.Sort == "trending" {
		// Smart Score: (Votes*2 + Comments*3 + Views*0.1 + 1) / (HoursSinceCreation + 2)^1.8
		scoreExpr := "((posts.vote_count * 2.0) + (posts.comment_count * 3.0) + (posts.view_count * 0.1) + 1) / POWER(EXTRACT(EPOCH FROM (NOW() - posts.created_at)) / 3600 + 2, 1.8)"

		if query.RequestUserID != uuid.Nil {
			db = db.Joins("LEFT JOIN follows ON follows.following_id = posts.author_id AND follows.follower_id = ?", query.RequestUserID)
			db = db.Order("CASE WHEN follows.follower_id IS NOT NULL THEN (" + scoreExpr + ") * 1.5 ELSE (" + scoreExpr + ") END DESC")
		} else {
			db = db.Order("(" + scoreExpr + ") DESC")
		}
	} else {
		db = db.Order("posts.created_at DESC, posts.id DESC")
	}

	err := db.Limit(limit + 1).Find(&posts).Error
	if err != nil {
		return nil, false, err
	}

	hasMore := len(posts) > limit
	if hasMore {
		posts = posts[:limit]
	}

	for i := range posts {
		var author models.User
		if err := r.db.WithContext(ctx).First(&author, "id = ?", posts[i].AuthorID).Error; err == nil {
			posts[i].Author = &models.UserPublicProfile{
				ID:          author.ID,
				Username:    author.Username,
				DisplayName: author.DisplayName,
				AvatarURL:   author.AvatarURL,
			}
		}
	}

	return posts, hasMore, nil
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, title, content, externalLinks *string, status *string) error {
	updates := make(map[string]interface{})
	if title != nil {
		updates["title"] = *title
	}
	if content != nil {
		updates["content"] = *content
	}
	if externalLinks != nil {
		updates["external_links"] = *externalLinks
	}
	if status != nil {
		updates["status"] = *status
	}
	updates["updated_at"] = time.Now()

	return r.db.WithContext(ctx).Model(&models.Post{}).Where("id = ?", id).Updates(updates).Error
}

func (r *Repository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Post{}).Where("id = ?", id).Update("status", "deleted").Error
}

func (r *Repository) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Post{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func (r *Repository) SlugExists(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Post{}).Where("slug = ?", slug).Count(&count).Error
	return count > 0, err
}

func (r *Repository) GetPostTags(ctx context.Context, postID uuid.UUID) ([]models.Tag, error) {
	var post models.Post
	err := r.db.WithContext(ctx).Preload("Tags").First(&post, "id = ?", postID).Error
	return post.Tags, err
}

func (r *Repository) SetPostTags(ctx context.Context, postID uuid.UUID, tagIDs []uuid.UUID) error {
	var post models.Post
	if err := r.db.WithContext(ctx).First(&post, "id = ?", postID).Error; err != nil {
		return err
	}

	var tags []models.Tag
	if len(tagIDs) > 0 {
		if err := r.db.WithContext(ctx).Find(&tags, "id IN ?", tagIDs).Error; err != nil {
			return err
		}
	}

	return r.db.WithContext(ctx).Model(&post).Association("Tags").Replace(tags)
}

func (r *Repository) IncrementUserPostCount(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).
		UpdateColumn("post_count", gorm.Expr("post_count + ?", 1)).Error
}

func (r *Repository) DecrementUserPostCount(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).
		UpdateColumn("post_count", gorm.Expr("GREATEST(post_count - ?, 0)", 1)).Error
}
