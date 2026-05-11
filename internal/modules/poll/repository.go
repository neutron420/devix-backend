package poll

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

func (r *Repository) Create(ctx context.Context, poll *models.Poll, options []models.PollOption) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(poll).Error; err != nil {
			return err
		}
		for i := range options {
			options[i].PollID = poll.ID
		}
		return tx.Create(&options).Error
	})
}

func (r *Repository) GetByPostID(ctx context.Context, postID uuid.UUID) (*models.Poll, []models.PollOption, error) {
	var poll models.Poll
	if err := r.db.WithContext(ctx).First(&poll, "post_id = ?", postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	var options []models.PollOption
	err := r.db.WithContext(ctx).Where("poll_id = ?", poll.ID).Find(&options).Error
	return &poll, options, err
}

func (r *Repository) Vote(ctx context.Context, pollID, optionID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Create Vote record
		vote := &models.PollVote{
			PollID:   pollID,
			OptionID: optionID,
			UserID:   userID,
		}
		if err := tx.Create(vote).Error; err != nil {
			return err // Will fail if already voted (composite PK)
		}

		// 2. Increment Option count
		return tx.Model(&models.PollOption{}).
			Where("id = ? AND poll_id = ?", optionID, pollID).
			Update("votes", gorm.Expr("votes + ?", 1)).Error
	})
}

func (r *Repository) GetUserVote(ctx context.Context, pollID, userID uuid.UUID) (uuid.UUID, error) {
	var vote models.PollVote
	err := r.db.WithContext(ctx).Where("poll_id = ? AND user_id = ?", pollID, userID).First(&vote).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return uuid.Nil, nil
	}
	return vote.OptionID, err
}
