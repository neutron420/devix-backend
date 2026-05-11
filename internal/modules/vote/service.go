package vote

import (
	"context"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/modules/notification"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	repo         *Repository
	notifService *notification.Service
	userService  interface {
		AdjustReputation(ctx context.Context, userID uuid.UUID, points int) error
	}
	log zerolog.Logger
}

func NewService(repo *Repository, notifService *notification.Service, log zerolog.Logger) *Service {
	return &Service{
		repo:         repo,
		notifService: notifService,
		log:          log.With().Str("module", "vote").Logger(),
	}
}

func (s *Service) SetUserService(us interface {
	AdjustReputation(ctx context.Context, userID uuid.UUID, points int) error
}) {
	s.userService = us
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

	go func() {
		authorID, err := s.repo.GetPostAuthorID(context.Background(), postID)
		if err != nil || authorID == userID {
			return
		}
		if voteType > 0 {
			if err := s.notifService.CreateNotification(context.Background(), authorID, userID, postID, "post_voted"); err != nil {
				s.log.Warn().Err(err).Msg("failed to create post vote notification")
			}
			if s.userService != nil {
				if err := s.userService.AdjustReputation(context.Background(), authorID, 5); err != nil {
					s.log.Warn().Err(err).Msg("failed to adjust author reputation")
				}
			}
		} else if voteType < 0 {
			if s.userService != nil {
				if err := s.userService.AdjustReputation(context.Background(), authorID, -2); err != nil {
					s.log.Warn().Err(err).Msg("failed to adjust author reputation")
				}
			}
		}
	}()

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

	go func() {
		authorID, err := s.repo.GetCommentAuthorID(context.Background(), commentID)
		if err != nil || authorID == userID {
			return
		}
		if voteType > 0 {
			if err := s.notifService.CreateNotification(context.Background(), authorID, userID, commentID, "comment_voted"); err != nil {
				s.log.Warn().Err(err).Msg("failed to create comment vote notification")
			}
			if s.userService != nil {
				if err := s.userService.AdjustReputation(context.Background(), authorID, 2); err != nil {
					s.log.Warn().Err(err).Msg("failed to adjust author reputation")
				}
			}
		} else if voteType < 0 {
			if s.userService != nil {
				if err := s.userService.AdjustReputation(context.Background(), authorID, -1); err != nil {
					s.log.Warn().Err(err).Msg("failed to adjust author reputation")
				}
			}
		}
	}()

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
