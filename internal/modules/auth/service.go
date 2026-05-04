package auth

import (
	"context"
	"strings"
	"time"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/models"
	"devix-backend/internal/pkg/hash"
	jwtpkg "devix-backend/internal/pkg/jwt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// Service handles auth business logic.
type Service struct {
	repo       *Repository
	jwt        *jwtpkg.Manager
	log        zerolog.Logger
}

// NewService creates a new auth service.
func NewService(repo *Repository, jwt *jwtpkg.Manager, log zerolog.Logger) *Service {
	return &Service{
		repo: repo,
		jwt:  jwt,
		log:  log.With().Str("module", "auth").Logger(),
	}
}

// Signup registers a new user.
func (s *Service) Signup(ctx context.Context, req *SignupRequest) (*AuthResponse, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)

	existing, _ := s.repo.GetUserByEmail(ctx, req.Email)
	if existing != nil {
		return nil, apperrors.Conflict("Email already taken")
	}

	existing, _ = s.repo.GetUserByUsername(ctx, req.Username)
	if existing != nil {
		return nil, apperrors.Conflict("Username already taken")
	}

	passwordHash, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	user := &models.User{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         "user",
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, apperrors.Internal(err)
	}

	tokenPair, err := s.jwt.GenerateTokenPair(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	tokenHash := HashToken(tokenPair.RefreshToken)
	_ = s.repo.StoreRefreshToken(ctx, user.ID, tokenHash, time.Now().Add(s.jwt.GetRefreshExpiry()))

	return &AuthResponse{
		User: UserResponse{
			ID:       user.ID.String(),
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
		Tokens: TokenResponse{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			ExpiresAt:    tokenPair.ExpiresAt.Format(time.RFC3339),
		},
	}, nil
}

// Login authenticates a user.
func (s *Service) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, apperrors.Unauthorized("Invalid credentials")
	}

	if valid, _ := hash.VerifyPassword(req.Password, user.PasswordHash); !valid {
		return nil, apperrors.Unauthorized("Invalid credentials")
	}

	tokenPair, err := s.jwt.GenerateTokenPair(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	tokenHash := HashToken(tokenPair.RefreshToken)
	_ = s.repo.StoreRefreshToken(ctx, user.ID, tokenHash, time.Now().Add(s.jwt.GetRefreshExpiry()))
	_ = s.repo.UpdateLastLogin(ctx, user.ID)

	return &AuthResponse{
		User: UserResponse{
			ID:       user.ID.String(),
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
		Tokens: TokenResponse{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			ExpiresAt:    tokenPair.ExpiresAt.Format(time.RFC3339),
		},
	}, nil
}

// RefreshTokens rotates tokens.
func (s *Service) RefreshTokens(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	tokenHash := HashToken(refreshToken)
	stored, err := s.repo.GetRefreshToken(ctx, tokenHash)
	if err != nil || stored == nil {
		return nil, apperrors.Unauthorized("Invalid token")
	}

	_ = s.repo.RevokeRefreshToken(ctx, tokenHash)

	user, err := s.repo.GetUserByID(ctx, stored.UserID)
	if err != nil || user == nil {
		return nil, apperrors.Unauthorized("User not found")
	}

	tokenPair, err := s.jwt.GenerateTokenPair(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	newTokenHash := HashToken(tokenPair.RefreshToken)
	_ = s.repo.StoreRefreshToken(ctx, user.ID, newTokenHash, time.Now().Add(s.jwt.GetRefreshExpiry()))

	return &AuthResponse{
		User: UserResponse{
			ID:       user.ID.String(),
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
		Tokens: TokenResponse{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			ExpiresAt:    tokenPair.ExpiresAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := HashToken(refreshToken)
	return s.repo.RevokeRefreshToken(ctx, tokenHash)
}
