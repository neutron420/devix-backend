package vote

import (
	"context"

	apperrors "devix-backend/internal/errors"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	repo *Repository
	log  zerolog.Logger
}

func NewService(repo *Repository, log zerolog.Logger) *Service {
	return &Service{repo: repo, log: log.With().Str("module", "vote").Logger()}
}

func (s *Service) VoteOnPost(ctx context.Context, userID, postID uuid.UUID, voteType int) (int, error) {
	if err := s.repo.UpsertPostVote(ctx, userID, postID, voteType); err != nil {
		s.log.Error().Err(err).Msg("failed to upsert post vote")
		return 0, apperrors.Internal(err)
	}
	count, err := s.repo.RecalcPostVoteCount(ctx, postID)
	if err != nil {
		return 0, apperrors.Internal(err)
	}
	return count, nil
}

func (s *Service) VoteOnComment(ctx context.Context, userID, commentID uuid.UUID, voteType int) (int, error) {
	if err := s.repo.UpsertCommentVote(ctx, userID, commentID, voteType); err != nil {
		s.log.Error().Err(err).Msg("failed to upsert comment vote")
		return 0, apperrors.Internal(err)
	}
	count, err := s.repo.RecalcCommentVoteCount(ctx, commentID)
	if err != nil {
		return 0, apperrors.Internal(err)
	}
	return count, nil
}

func (s *Service) RemovePostVote(ctx context.Context, userID, postID uuid.UUID) (int, error) {
	if err := s.repo.DeletePostVote(ctx, userID, postID); err != nil {
		return 0, apperrors.Internal(err)
	}
	count, err := s.repo.RecalcPostVoteCount(ctx, postID)
	if err != nil {
		return 0, apperrors.Internal(err)
	}
	return count, nil
}

func (s *Service) RemoveCommentVote(ctx context.Context, userID, commentID uuid.UUID) (int, error) {
	if err := s.repo.DeleteCommentVote(ctx, userID, commentID); err != nil {
		return 0, apperrors.Internal(err)
	}
	count, err := s.repo.RecalcCommentVoteCount(ctx, commentID)
	if err != nil {
		return 0, apperrors.Internal(err)
	}
	return count, nil
}
