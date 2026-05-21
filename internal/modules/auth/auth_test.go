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

func TestAuthService(t *testing.T) {
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
		assert.NotEmpty(t, res.Tokens.AccessToken)
		assert.NotEmpty(t, res.Tokens.RefreshToken)
	})

	t.Run("duplicate signup", func(t *testing.T) {
		req := &SignupRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		_, err := service.Signup(context.Background(), req)
		assert.Error(t, err)
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
	})

	t.Run("login invalid credentials", func(t *testing.T) {
		req := &LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		_, err := service.Login(context.Background(), req, "127.0.0.1")
		assert.Error(t, err)
	})
}
