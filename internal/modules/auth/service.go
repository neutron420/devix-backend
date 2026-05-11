package auth

import (
	"context"
	"strings"
	"time"

	"devix-backend/internal/config"
	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/models"
	"devix-backend/internal/pkg/email"
	"devix-backend/internal/pkg/hash"
	jwtpkg "devix-backend/internal/pkg/jwt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	repo           *Repository
	jwt            *jwtpkg.Manager
	mailer         *email.Mailer
	lockout        *LockoutService
	passwordPolicy config.PasswordPolicyConfig
	log            zerolog.Logger
}

func NewService(repo *Repository, jwt *jwtpkg.Manager, mailer *email.Mailer, lockout *LockoutService, policy config.PasswordPolicyConfig, log zerolog.Logger) *Service {
	return &Service{
		repo:           repo,
		jwt:            jwt,
		mailer:         mailer,
		lockout:        lockout,
		passwordPolicy: policy,
		log:            log.With().Str("module", "auth").Logger(),
	}
}

func (s *Service) Signup(ctx context.Context, req *SignupRequest) (*AuthResponse, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)

	if err := ValidatePassword(s.passwordPolicy, req.Password); err != nil {
		return nil, apperrors.BadRequest(err.Error())
	}

	existing, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if existing != nil {
		return nil, apperrors.Conflict("Email already taken")
	}

	existing, err = s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if existing != nil {
		return nil, apperrors.Conflict("Username already taken")
	}

	passwordHash, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	verificationToken := uuid.New().String()
	user := &models.User{
		ID:                uuid.New(),
		Username:          req.Username,
		Email:             req.Email,
		PasswordHash:      passwordHash,
		Role:              "user",
		VerificationToken: &verificationToken,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, apperrors.Internal(err)
	}

	// Send verification email asynchronously
	go func() {
		if err := s.mailer.SendVerificationEmail(user.Email, user.Username, verificationToken); err != nil {
			s.log.Warn().Err(err).Str("email", user.Email).Msg("failed to send verification email")
		}
	}()

	tokenPair, err := s.jwt.GenerateTokenPair(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	tokenHash := HashToken(tokenPair.RefreshToken)
	if err := s.repo.StoreRefreshToken(ctx, user.ID, tokenHash, time.Now().Add(s.jwt.GetRefreshExpiry())); err != nil {
		s.log.Warn().Err(err).Str("user_id", user.ID.String()).Msg("failed to store refresh token")
	}

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

func (s *Service) Login(ctx context.Context, req *LoginRequest, clientIP string) (*AuthResponse, error) {
	emailKey := strings.ToLower(strings.TrimSpace(req.Email))
	lockoutEmailKey := "email:" + emailKey
	lockoutIPKey := ""
	if strings.TrimSpace(clientIP) != "" {
		lockoutIPKey = "ip:" + strings.TrimSpace(clientIP)
	}

	if s.lockout != nil {
		if retryAfter, locked := s.lockout.Check(ctx, s.lockout.NormalizeKey(lockoutEmailKey)); locked {
			return nil, &LockoutError{RetryAfter: retryAfter}
		}
		if lockoutIPKey != "" {
			if retryAfter, locked := s.lockout.Check(ctx, s.lockout.NormalizeKey(lockoutIPKey)); locked {
				return nil, &LockoutError{RetryAfter: retryAfter}
			}
		}
	}

	user, err := s.repo.GetUserByEmail(ctx, emailKey)
	if err != nil || user == nil {
		if s.lockout != nil {
			s.lockout.RegisterFailure(ctx, s.lockout.NormalizeKey(lockoutEmailKey))
			if lockoutIPKey != "" {
				s.lockout.RegisterFailure(ctx, s.lockout.NormalizeKey(lockoutIPKey))
			}
		}
		return nil, apperrors.Unauthorized("Invalid credentials")
	}

	if !user.IsActive {
		return nil, apperrors.Forbidden("Account is deactivated. Please contact support or reactivate.")
	}

	valid, err := hash.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil || !valid {
		if s.lockout != nil {
			s.lockout.RegisterFailure(ctx, s.lockout.NormalizeKey(lockoutEmailKey))
			if lockoutIPKey != "" {
				s.lockout.RegisterFailure(ctx, s.lockout.NormalizeKey(lockoutIPKey))
			}
		}
		return nil, apperrors.Unauthorized("Invalid credentials")
	}
	if s.lockout != nil {
		s.lockout.Clear(ctx, s.lockout.NormalizeKey(lockoutEmailKey))
		if lockoutIPKey != "" {
			s.lockout.Clear(ctx, s.lockout.NormalizeKey(lockoutIPKey))
		}
	}

	tokenPair, err := s.jwt.GenerateTokenPair(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	tokenHash := HashToken(tokenPair.RefreshToken)
	if err := s.repo.StoreRefreshToken(ctx, user.ID, tokenHash, time.Now().Add(s.jwt.GetRefreshExpiry())); err != nil {
		s.log.Warn().Err(err).Str("user_id", user.ID.String()).Msg("failed to store refresh token")
	}
	if err := s.repo.UpdateLastLogin(ctx, user.ID); err != nil {
		s.log.Warn().Err(err).Str("user_id", user.ID.String()).Msg("failed to update last login")
	}

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

func (s *Service) RefreshTokens(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	tokenHash := HashToken(refreshToken)
	stored, err := s.repo.GetRefreshToken(ctx, tokenHash)
	if err != nil || stored == nil {
		return nil, apperrors.Unauthorized("Invalid token")
	}

	if err := s.repo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		s.log.Warn().Err(err).Str("token_hash", tokenHash).Msg("failed to revoke old token")
	}

	user, err := s.repo.GetUserByID(ctx, stored.UserID)
	if err != nil || user == nil {
		return nil, apperrors.Unauthorized("User not found")
	}

	if !user.IsActive {
		return nil, apperrors.Forbidden("Account is deactivated")
	}

	tokenPair, err := s.jwt.GenerateTokenPair(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	newTokenHash := HashToken(tokenPair.RefreshToken)
	if err := s.repo.StoreRefreshToken(ctx, user.ID, newTokenHash, time.Now().Add(s.jwt.GetRefreshExpiry())); err != nil {
		s.log.Warn().Err(err).Str("user_id", user.ID.String()).Msg("failed to store new refresh token")
	}

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

func (s *Service) VerifyEmail(ctx context.Context, token string) error {
	user, err := s.repo.GetUserByVerificationToken(ctx, token)
	if err != nil {
		return apperrors.Internal(err)
	}
	if user == nil {
		return apperrors.BadRequest("Invalid or expired verification token")
	}

	return s.repo.SetUserVerified(ctx, user.ID)
}

func (s *Service) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return apperrors.Internal(err)
	}
	if user == nil {
		// Don't reveal if user exists for security
		return nil
	}

	resetToken := uuid.New().String()
	expiresAt := time.Now().Add(1 * time.Hour)

	if err := s.repo.SetUserResetToken(ctx, user.ID, resetToken, expiresAt); err != nil {
		return apperrors.Internal(err)
	}

	go func() {
		if err := s.mailer.SendPasswordResetEmail(user.Email, resetToken); err != nil {
			s.log.Warn().Err(err).Str("email", user.Email).Msg("failed to send password reset email")
		}
	}()

	return nil
}

func (s *Service) ResetPassword(ctx context.Context, token, newPassword string) error {
	user, err := s.repo.GetUserByResetToken(ctx, token)
	if err != nil {
		return apperrors.Internal(err)
	}
	if user == nil {
		return apperrors.BadRequest("Invalid or expired reset token")
	}

	if err := ValidatePassword(s.passwordPolicy, newPassword); err != nil {
		return apperrors.BadRequest(err.Error())
	}

	passwordHash, err := hash.HashPassword(newPassword)
	if err != nil {
		return apperrors.Internal(err)
	}

	return s.repo.UpdatePassword(ctx, user.ID, passwordHash)
}
