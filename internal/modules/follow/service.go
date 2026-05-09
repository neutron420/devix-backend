package follow

import (
	"context"
	"time"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/models"
	"devix-backend/internal/modules/notification"
	"devix-backend/internal/modules/user"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	repo         *Repository
	userService  *user.Service
	notifService *notification.Service
	log          zerolog.Logger
}

func NewService(repo *Repository, userService *user.Service, notifService *notification.Service, log zerolog.Logger) *Service {
	return &Service{
		repo:         repo,
		userService:  userService,
		notifService: notifService,
		log:          log.With().Str("module", "follow").Logger(),
	}
}

func (s *Service) Follow(ctx context.Context, followerID uuid.UUID, targetUsername string) error {
	targetUser, err := s.userService.GetByUsername(ctx, targetUsername)
	if err != nil {
		return err
	}
	if targetUser == nil {
		return apperrors.NotFound("User")
	}

	targetUserID, _ := uuid.Parse(targetUser.ID)
	if followerID == targetUserID {
		return apperrors.BadRequest("You cannot follow yourself")
	}

	isFollowing, err := s.repo.IsFollowing(ctx, followerID, targetUserID)
	if err != nil {
		return apperrors.Internal(err)
	}
	if isFollowing {
		return apperrors.BadRequest("Already following this user")
	}

	follow := &models.Follow{
		FollowerID:  followerID,
		FollowingID: targetUserID,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Follow(ctx, follow); err != nil {
		return apperrors.Internal(err)
	}

	// Trigger Notification
	go func() {
		_ = s.notifService.CreateNotification(context.Background(), targetUserID, followerID, followerID, "followed")
	}()

	return nil
}

func (s *Service) Unfollow(ctx context.Context, followerID uuid.UUID, targetUsername string) error {
	targetUser, err := s.userService.GetByUsername(ctx, targetUsername)
	if err != nil {
		return err
	}
	if targetUser == nil {
		return apperrors.NotFound("User")
	}

	targetUserID, _ := uuid.Parse(targetUser.ID)
	return s.repo.Unfollow(ctx, followerID, targetUserID)
}

func (s *Service) GetFollowers(ctx context.Context, username string) ([]UserProfileResponse, error) {
	targetUser, err := s.userService.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	targetUserID, _ := uuid.Parse(targetUser.ID)

	users, err := s.repo.GetFollowers(ctx, targetUserID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	return s.toProfileList(users), nil
}

func (s *Service) GetFollowing(ctx context.Context, username string) ([]UserProfileResponse, error) {
	targetUser, err := s.userService.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	targetUserID, _ := uuid.Parse(targetUser.ID)

	users, err := s.repo.GetFollowing(ctx, targetUserID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	return s.toProfileList(users), nil
}

func (s *Service) GetFollowingIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	return s.repo.GetFollowingIDs(ctx, userID)
}

func (s *Service) toProfileList(users []models.User) []UserProfileResponse {
	resp := make([]UserProfileResponse, 0, len(users))
	for _, u := range users {
		displayName := ""
		if u.DisplayName != nil {
			displayName = *u.DisplayName
		}
		avatarURL := ""
		if u.AvatarURL != nil {
			avatarURL = *u.AvatarURL
		}
		resp = append(resp, UserProfileResponse{
			ID:          u.ID.String(),
			Username:    u.Username,
			DisplayName: displayName,
			AvatarURL:   avatarURL,
		})
	}
	return resp
}
