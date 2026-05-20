package user

import (
	"context"
	"strings"
	"time"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/models"
	"devix-backend/internal/pkg/cache"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	repo  *Repository
	cache *cache.Cache
	log   zerolog.Logger
}

func (s *Service) GetByUsername(ctx context.Context, username string) (*ProfileResponse, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if user == nil {
		return nil, nil
	}
	res, err := s.toResponse(ctx, user)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return res, nil
}

func NewService(repo *Repository, cache *cache.Cache, log zerolog.Logger) *Service {
	return &Service{
		repo:  repo,
		cache: cache,
		log:   log.With().Str("module", "user").Logger(),
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

	res, err := s.toResponse(ctx, user)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return res, nil
}

func (s *Service) GetPublicProfile(ctx context.Context, username string) (*PublicProfileResponse, error) {
	cacheKey := "users:profile:" + username
	var cached PublicProfileResponse
	if ok, _ := s.cache.Get(ctx, cacheKey, &cached); ok {
		return &cached, nil
	}

	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get user")
		return nil, apperrors.Internal(err)
	}
	if user == nil {
		return nil, apperrors.NotFound("User")
	}
	followerCount, err := s.repo.CountFollowers(ctx, user.ID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	followingCount, err := s.repo.CountFollowing(ctx, user.ID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	res := &PublicProfileResponse{
		ID:          user.ID.String(),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Bio:         user.Bio,
		AvatarURL:   user.AvatarURL,
		PostCount:   user.PostCount,
		FollowerCount: followerCount,
		FollowingCount: followingCount,
		Reputation:  user.Reputation,
		Level:       CalculateLevel(user.Reputation),
		Badges:      CalculateBadges(user.Reputation),
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
	}

	_ = s.cache.Set(ctx, cacheKey, res, 10*time.Minute)
	return res, nil
}

func (s *Service) GetPublicProfileByID(ctx context.Context, userID uuid.UUID) (*PublicProfileResponse, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, err
	}
	followerCount, err := s.repo.CountFollowers(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	followingCount, err := s.repo.CountFollowing(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return &PublicProfileResponse{
		ID:          user.ID.String(),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		PostCount:   user.PostCount,
		FollowerCount: followerCount,
		FollowingCount: followingCount,
		Reputation:  user.Reputation,
		Level:       CalculateLevel(user.Reputation),
		Badges:      CalculateBadges(user.Reputation),
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

	// Invalidate cache
	user, _ := s.repo.GetByID(ctx, userID)
	if user != nil {
		_ = s.cache.Delete(ctx, "users:profile:"+user.Username)
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

func (s *Service) UpdateStatus(ctx context.Context, userID uuid.UUID, isActive bool) error {
	if err := s.repo.UpdateStatus(ctx, userID, isActive); err != nil {
		s.log.Error().Err(err).Msg("failed to update user status")
		return apperrors.Internal(err)
	}
	return nil
}

func (s *Service) toResponse(ctx context.Context, user *models.User) (*ProfileResponse, error) {
	followerCount, err := s.repo.CountFollowers(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	followingCount, err := s.repo.CountFollowing(ctx, user.ID)
	if err != nil {
		return nil, err
	}
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
		FollowerCount: followerCount,
		FollowingCount: followingCount,
		Reputation:  user.Reputation,
		Level:       CalculateLevel(user.Reputation),
		Badges:      CalculateBadges(user.Reputation),
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *Service) UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarURL string) (*ProfileResponse, error) {
	if err := s.repo.UpdateAvatar(ctx, userID, avatarURL); err != nil {
		s.log.Error().Err(err).Msg("failed to update avatar")
		return nil, apperrors.Internal(err)
	}

	// Invalidate cache
	user, _ := s.repo.GetByID(ctx, userID)
	if user != nil {
		_ = s.cache.Delete(ctx, "users:profile:"+user.Username)
	}

	return s.GetMyProfile(ctx, userID)
}

func (s *Service) InvalidatePublicProfileCache(ctx context.Context, username string) {
	if strings.TrimSpace(username) == "" {
		return
	}
	_ = s.cache.Delete(ctx, "users:profile:"+username)
}

func (s *Service) InvalidateProfileCacheByID(ctx context.Context, userID uuid.UUID) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return
	}
	_ = s.cache.Delete(ctx, "users:profile:"+user.Username)
}

func (s *Service) AdjustReputation(ctx context.Context, userID uuid.UUID, points int) error {
	if err := s.repo.AdjustReputation(ctx, userID, points); err != nil {
		s.log.Error().Err(err).Str("user_id", userID.String()).Int("points", points).Msg("failed to adjust reputation")
		return apperrors.Internal(err)
	}
	_ = s.cache.DeleteByPattern(ctx, "users:profile:*")
	return nil
}
