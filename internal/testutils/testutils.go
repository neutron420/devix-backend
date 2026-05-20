package testutils

import (
	"io"
	"testing"
	"time"

	"devix-backend/internal/models"
	jwtpkg "devix-backend/internal/pkg/jwt"

	"github.com/rs/zerolog"
)

func NewTestJWT() *jwtpkg.Manager {
	return jwtpkg.NewManager(
		"test-access-secret-123456789012345678901234",
		"test-refresh-secret-123456789012345678901234",
		15*time.Minute,
		24*time.Hour,
	)
}

func NewTestLogger() zerolog.Logger {
	return zerolog.New(io.Discard)
}

func AccessTokenFor(t *testing.T, jwt *jwtpkg.Manager, user *models.User) string {
	t.Helper()
	tokens, err := jwt.GenerateTokenPair(user.ID, user.Username, user.Role)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	return tokens.AccessToken
}
