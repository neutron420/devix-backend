package user

import (
	"context"
	"time"

	apperrors "devix-backend/internal/errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// Service handles user business logic.
type Service struct {
	repo *Repository
	log  zerolog.Logger
}

// NewService creates a new user service.
func NewService(repo *Repository, log zerolog.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log.With().Str("module", "user").Logger(),
	}
}

// GetMyProfile returns the authenticated user's full profile.
func (s *Service) GetMyProfile(ctx context.Context, userID uuid.UUID) (*ProfileResponse, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get user")
		return nil, apperrors.Internal(err)
	}
	if user == nil {
		return nil, apperrors.NotFound("User")
	}

	return &ProfileResponse{
		ID:          user.ID.String(),
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Bio:         user.Bio,
		AvatarURL:   user.AvatarURL,
		Role:        user.Role,
		IsVerified:  user.IsVerified,
		PostCount:   user.PostCount,
		Reputation:  user.Reputation,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
	}, nil
}

// GetPublicProfile returns a public user profile by username.
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

// UpdateProfile updates the authenticated user's profile.
func (s *Service) UpdateProfile(ctx context.Context, userID uuid.UUID, req *UpdateProfileRequest) (*ProfileResponse, error) {
	// Check username uniqueness if being changed
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

	if err := s.repo.UpdateProfile(ctx, userID, req.DisplayName, req.Bio, req.Username); err != nil {
		s.log.Error().Err(err).Msg("failed to update profile")
		return nil, apperrors.Internal(err)
	}

	return s.GetMyProfile(ctx, userID)
}

// UpdateAvatar updates the user's profile picture URL.
func (s *Service) UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarURL string) (*ProfileResponse, error) {
	if err := s.repo.UpdateAvatar(ctx, userID, avatarURL); err != nil {
		s.log.Error().Err(err).Msg("failed to update avatar")
		return nil, apperrors.Internal(err)
	}

	return s.GetMyProfile(ctx, userID)
}
