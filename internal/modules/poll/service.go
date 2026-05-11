package poll

import (
	"context"
	"time"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	repo *Repository
	log  zerolog.Logger
}

func NewService(repo *Repository, log zerolog.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log.With().Str("module", "poll").Logger(),
	}
}

func (s *Service) CreatePoll(ctx context.Context, req *CreatePollRequest) (*PollResponse, error) {
	postID, err := uuid.Parse(req.PostID)
	if err != nil {
		return nil, apperrors.BadRequest("Invalid post ID")
	}

	if req.ExpiresAt.Before(time.Now()) {
		return nil, apperrors.BadRequest("Expiration date must be in the future")
	}

	poll := &models.Poll{
		ID:        uuid.New(),
		PostID:    postID,
		Question:  req.Question,
		ExpiresAt: req.ExpiresAt,
		CreatedAt: time.Now(),
	}

	options := make([]models.PollOption, 0, len(req.Options))
	for _, text := range req.Options {
		options = append(options, models.PollOption{
			ID:   uuid.New(),
			Text: text,
		})
	}

	if err := s.repo.Create(ctx, poll, options); err != nil {
		return nil, apperrors.Internal(err)
	}

	return s.toResponse(poll, options, uuid.Nil), nil
}

func (s *Service) GetPoll(ctx context.Context, postID uuid.UUID, userID uuid.UUID) (*PollResponse, error) {
	poll, options, err := s.repo.GetByPostID(ctx, postID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if poll == nil {
		return nil, nil
	}

	var votedID uuid.UUID
	if userID != uuid.Nil {
		votedID, _ = s.repo.GetUserVote(ctx, poll.ID, userID)
	}

	return s.toResponse(poll, options, votedID), nil
}

func (s *Service) Vote(ctx context.Context, pollID, optionID, userID uuid.UUID) error {
	if err := s.repo.Vote(ctx, pollID, optionID, userID); err != nil {
		return apperrors.Conflict("You have already voted or poll is closed")
	}
	return nil
}

func (s *Service) toResponse(poll *models.Poll, options []models.PollOption, votedID uuid.UUID) *PollResponse {
	total := 0
	opts := make([]OptionResp, 0, len(options))
	for _, o := range options {
		total += o.Votes
		opts = append(opts, OptionResp{
			ID:    o.ID.String(),
			Text:  o.Text,
			Votes: o.Votes,
		})
	}

	return &PollResponse{
		ID:         poll.ID.String(),
		Question:   poll.Question,
		Options:    opts,
		ExpiresAt:  poll.ExpiresAt,
		TotalVotes: total,
		HasVoted:   votedID != uuid.Nil,
		VotedID:    votedID.String(),
	}
}
