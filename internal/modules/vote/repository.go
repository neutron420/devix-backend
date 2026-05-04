package vote

import (
	"context"
	"errors"

	"devix-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetPostVote(ctx context.Context, userID, postID uuid.UUID) (*models.Vote, error) {
	var v models.Vote
	err := r.db.WithContext(ctx).Where("user_id = ? AND post_id = ?", userID, postID).First(&v).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &v, err
}

func (r *Repository) GetCommentVote(ctx context.Context, userID, commentID uuid.UUID) (*models.Vote, error) {
	var v models.Vote
	err := r.db.WithContext(ctx).Where("user_id = ? AND comment_id = ?", userID, commentID).First(&v).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &v, err
}

func (r *Repository) UpsertPostVote(ctx context.Context, userID, postID uuid.UUID, voteType int) error {
	vote := models.Vote{
		ID:       uuid.New(),
		UserID:   userID,
		PostID:   &postID,
		VoteType: voteType,
	}

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "post_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"vote_type": voteType}),
	}).Create(&vote).Error
}

func (r *Repository) UpsertCommentVote(ctx context.Context, userID, commentID uuid.UUID, voteType int) error {
	vote := models.Vote{
		ID:        uuid.New(),
		UserID:    userID,
		CommentID: &commentID,
		VoteType:  voteType,
	}

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "comment_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"vote_type": voteType}),
	}).Create(&vote).Error
}

func (r *Repository) DeletePostVote(ctx context.Context, userID, postID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND post_id = ?", userID, postID).Delete(&models.Vote{}).Error
}

func (r *Repository) DeleteCommentVote(ctx context.Context, userID, commentID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND comment_id = ?", userID, commentID).Delete(&models.Vote{}).Error
}

func (r *Repository) RecalcPostVoteCount(ctx context.Context, postID uuid.UUID) (int, error) {
	var sum int64
	err := r.db.WithContext(ctx).Model(&models.Vote{}).Where("post_id = ?", postID).Select("COALESCE(SUM(vote_type), 0)").Scan(&sum).Error
	if err != nil {
		return 0, err
	}
	err = r.db.WithContext(ctx).Model(&models.Post{}).Where("id = ?", postID).Update("vote_count", sum).Error
	return int(sum), err
}

func (r *Repository) RecalcCommentVoteCount(ctx context.Context, commentID uuid.UUID) (int, error) {
	var sum int64
	err := r.db.WithContext(ctx).Model(&models.Vote{}).Where("comment_id = ?", commentID).Select("COALESCE(SUM(vote_type), 0)").Scan(&sum).Error
	if err != nil {
		return 0, err
	}
	err = r.db.WithContext(ctx).Model(&models.Comment{}).Where("id = ?", commentID).Update("vote_count", sum).Error
	return int(sum), err
}
