package user

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

func (s *Service) GetByUsername(ctx context.Context, username string) (*ProfileResponse, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if user == nil {
		return nil, nil
	}
	return s.toResponse(user), nil
}

func NewService(repo *Repository, log zerolog.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log.With().Str("module", "user").Logger(),
	}
}

func (s *Service) GetMyProfile(ctx context.Context, userID uuid.UUID) (*ProfileResponse, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get user")
		return nil, apperrors.Internal(err)
	}
	if user == nil {
		return nil, apperrors.NotFound("User")
	}

	return s.toResponse(user), nil
}

func (s *Service) GetPublicProfile(ctx context.Context, username string) (*PublicProfileResponse, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get user")
		return nil, apperrors.Internal(err)
	}
	if user == nil {
		return nil, apperrors.NotFound("User")
	}

	return &PublicProfileResponse{
		ID:          user.ID.String(),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Bio:         user.Bio,
		AvatarURL:   user.AvatarURL,
		PostCount:   user.PostCount,
		Reputation:  user.Reputation,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID uuid.UUID, req *UpdateProfileRequest) (*ProfileResponse, error) {

	if req.Username != nil {
		exists, err := s.repo.UsernameExists(ctx, *req.Username, userID)
		if err != nil {
			s.log.Error().Err(err).Msg("failed to check username")
			return nil, apperrors.Internal(err)
		}
		if exists {
			return nil, apperrors.Conflict("This username is already taken")
		}
	}

	if err := s.repo.UpdateProfile(ctx, userID, req); err != nil {
		s.log.Error().Err(err).Msg("failed to update profile")
		return nil, apperrors.Internal(err)
	}

	return s.GetMyProfile(ctx, userID)
}

func (s *Service) DeleteAccount(ctx context.Context, userID uuid.UUID) error {
	if err := s.repo.Delete(ctx, userID); err != nil {
		s.log.Error().Err(err).Msg("failed to delete user")
		return apperrors.Internal(err)
	}
	return nil
}

func (s *Service) toResponse(user *models.User) *ProfileResponse {
	return &ProfileResponse{
		ID:          user.ID.String(),
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Bio:         user.Bio,
		AvatarURL:   user.AvatarURL,
		WebsiteURL:  user.WebsiteURL,
		GitHubURL:   user.GitHubURL,
		TwitterURL:  user.TwitterURL,
		Location:    user.Location,
		Preferences: user.Preferences,
		Role:        user.Role,
		IsVerified:  user.IsVerified,
		PostCount:   user.PostCount,
		Reputation:  user.Reputation,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
	}
}

func (s *Service) UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarURL string) (*ProfileResponse, error) {
	if err := s.repo.UpdateAvatar(ctx, userID, avatarURL); err != nil {
		s.log.Error().Err(err).Msg("failed to update avatar")
		return nil, apperrors.Internal(err)
	}

	return s.GetMyProfile(ctx, userID)
}
