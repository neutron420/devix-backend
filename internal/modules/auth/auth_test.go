package auth

import (
	"context"
	"testing"

	"devix-backend/internal/config"
	"devix-backend/internal/pkg/email"
	"devix-backend/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAuthService(t *testing.T) (*Service, *Repository) {
	t.Helper()
	db := testutils.SetupTestDB(t)
	repo := NewRepository(db)
	jwt := testutils.NewTestJWT()
	mailer := email.NewMailer(config.EmailConfig{})
	policy := config.PasswordPolicyConfig{
		MinLength:     8,
		RequireUpper:  false,
		RequireLower:  false,
		RequireNumber: false,
		RequireSymbol: false,
	}
	logger := testutils.NewTestLogger()
	service := NewService(repo, jwt, mailer, nil, policy, logger)
	return service, repo
}

func TestAuthService(t *testing.T) {
	service, _ := setupAuthService(t)

	t.Run("successful signup", func(t *testing.T) {
		req := &SignupRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		res, err := service.Signup(context.Background(), req)
		require.NoError(t, err)
		assert.Equal(t, "testuser", res.User.Username)
		assert.Equal(t, "test@example.com", res.User.Email)
		assert.Equal(t, "user", res.User.Role)
		assert.NotEmpty(t, res.Tokens.AccessToken)
		assert.NotEmpty(t, res.Tokens.RefreshToken)
		assert.NotEmpty(t, res.Tokens.ExpiresAt)
	})

	t.Run("duplicate email signup", func(t *testing.T) {
		req := &SignupRequest{
			Username: "anotheruser",
			Email:    "test@example.com",
			Password: "password123",
		}

		_, err := service.Signup(context.Background(), req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Email already taken")
	})

	t.Run("duplicate username signup", func(t *testing.T) {
		req := &SignupRequest{
			Username: "testuser",
			Email:    "another@example.com",
			Password: "password123",
		}

		_, err := service.Signup(context.Background(), req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Username already taken")
	})

	t.Run("signup with weak password", func(t *testing.T) {
		svc, _ := setupAuthService(t)
		svc.passwordPolicy = config.PasswordPolicyConfig{
			MinLength:     10,
			RequireUpper:  true,
			RequireLower:  true,
			RequireNumber: true,
			RequireSymbol: true,
		}
		req := &SignupRequest{
			Username: "strongpwuser",
			Email:    "strong@example.com",
			Password: "short",
		}

		_, err := svc.Signup(context.Background(), req)
		assert.Error(t, err)
	})

	t.Run("signup normalizes email to lowercase", func(t *testing.T) {
		req := &SignupRequest{
			Username: "caseuser",
			Email:    "  Case@Example.COM  ",
			Password: "password123",
		}

		res, err := service.Signup(context.Background(), req)
		require.NoError(t, err)
		assert.Equal(t, "case@example.com", res.User.Email)
	})

	t.Run("successful login", func(t *testing.T) {
		req := &LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		res, err := service.Login(context.Background(), req, "127.0.0.1")
		require.NoError(t, err)
		assert.Equal(t, "testuser", res.User.Username)
		assert.NotEmpty(t, res.Tokens.AccessToken)
		assert.NotEmpty(t, res.Tokens.RefreshToken)
	})

	t.Run("login with wrong password", func(t *testing.T) {
		req := &LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		_, err := service.Login(context.Background(), req, "127.0.0.1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid credentials")
	})

	t.Run("login with non-existent email", func(t *testing.T) {
		req := &LoginRequest{
			Email:    "nobody@example.com",
			Password: "password123",
		}

		_, err := service.Login(context.Background(), req, "127.0.0.1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid credentials")
	})

	t.Run("login with deactivated account", func(t *testing.T) {
		// Create and deactivate a user
		signupReq := &SignupRequest{
			Username: "deactivated",
			Email:    "deactivated@example.com",
			Password: "password123",
		}
		_, err := service.Signup(context.Background(), signupReq)
		require.NoError(t, err)

		user, err := service.repo.GetUserByEmail(context.Background(), "deactivated@example.com")
		require.NoError(t, err)
		user.IsActive = false
		service.repo.db.Save(user)

		loginReq := &LoginRequest{
			Email:    "deactivated@example.com",
			Password: "password123",
		}
		_, err = service.Login(context.Background(), loginReq, "127.0.0.1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "deactivated")
	})

	t.Run("token refresh", func(t *testing.T) {
		loginReq := &LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		loginRes, err := service.Login(context.Background(), loginReq, "127.0.0.1")
		require.NoError(t, err)

		refreshRes, err := service.RefreshTokens(context.Background(), loginRes.Tokens.RefreshToken)
		require.NoError(t, err)
		assert.NotEmpty(t, refreshRes.Tokens.AccessToken)
		assert.NotEmpty(t, refreshRes.Tokens.RefreshToken)
		// Old refresh token should be rotated (different from original)
		assert.NotEqual(t, loginRes.Tokens.RefreshToken, refreshRes.Tokens.RefreshToken)
	})

	t.Run("refresh with invalid token", func(t *testing.T) {
		_, err := service.RefreshTokens(context.Background(), "invalid-token-string")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid token")
	})

	t.Run("refresh with already-used token (rotation enforcement)", func(t *testing.T) {
		loginReq := &LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		loginRes, err := service.Login(context.Background(), loginReq, "127.0.0.1")
		require.NoError(t, err)

		// First refresh succeeds
		_, err = service.RefreshTokens(context.Background(), loginRes.Tokens.RefreshToken)
		require.NoError(t, err)

		// Second refresh with the same old token fails (already revoked)
		_, err = service.RefreshTokens(context.Background(), loginRes.Tokens.RefreshToken)
		assert.Error(t, err)
	})

	t.Run("logout", func(t *testing.T) {
		loginReq := &LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		loginRes, err := service.Login(context.Background(), loginReq, "127.0.0.1")
		require.NoError(t, err)

		err = service.Logout(context.Background(), loginRes.Tokens.RefreshToken)
		assert.NoError(t, err)

		// After logout, refresh should fail
		_, err = service.RefreshTokens(context.Background(), loginRes.Tokens.RefreshToken)
		assert.Error(t, err)
	})

	t.Run("verify email with valid token", func(t *testing.T) {
		// Get the user's verification token
		user, err := service.repo.GetUserByEmail(context.Background(), "test@example.com")
		require.NoError(t, err)
		require.NotNil(t, user.VerificationToken)

		err = service.VerifyEmail(context.Background(), *user.VerificationToken)
		assert.NoError(t, err)

		// Verify user is now verified
		updated, err := service.repo.GetUserByEmail(context.Background(), "test@example.com")
		require.NoError(t, err)
		assert.True(t, updated.IsVerified)
	})

	t.Run("verify email with invalid token", func(t *testing.T) {
		err := service.VerifyEmail(context.Background(), "bogus-token-that-does-not-exist")
		assert.Error(t, err)
	})

	t.Run("forgot password for existing user", func(t *testing.T) {
		err := service.ForgotPassword(context.Background(), "test@example.com")
		assert.NoError(t, err)

		// Should have set a reset token
		user, err := service.repo.GetUserByEmail(context.Background(), "test@example.com")
		require.NoError(t, err)
		assert.NotNil(t, user.ResetToken)
		assert.NotNil(t, user.ResetExpiresAt)
	})

	t.Run("forgot password for non-existent email does not error", func(t *testing.T) {
		// Should not error (to avoid revealing whether a user exists)
		err := service.ForgotPassword(context.Background(), "nobody@nonexistent.com")
		assert.NoError(t, err)
	})

	t.Run("reset password with valid token", func(t *testing.T) {
		user, err := service.repo.GetUserByEmail(context.Background(), "test@example.com")
		require.NoError(t, err)
		require.NotNil(t, user.ResetToken)

		err = service.ResetPassword(context.Background(), *user.ResetToken, "newpassword123")
		assert.NoError(t, err)

		// Can login with new password
		loginReq := &LoginRequest{
			Email:    "test@example.com",
			Password: "newpassword123",
		}
		_, err = service.Login(context.Background(), loginReq, "127.0.0.1")
		assert.NoError(t, err)
	})

	t.Run("reset password with invalid token", func(t *testing.T) {
		err := service.ResetPassword(context.Background(), "invalid-reset-token", "newpassword123")
		assert.Error(t, err)
	})
}
