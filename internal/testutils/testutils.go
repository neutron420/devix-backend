package testutils

import (
	"io"
	"strings"
	"testing"
	"time"

	"devix-backend/internal/models"
	jwtpkg "devix-backend/internal/pkg/jwt"

	"github.com/rs/zerolog"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
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

func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test sqlite database: %v", err)
	}

	// Remove PostgreSQL-specific DEFAULT gen_random_uuid() for SQLite compatibility during migrations
	db.Callback().Raw().Before("gorm:raw").Register("sqlite_uuid_fix", func(d *gorm.DB) {
		sql := d.Statement.SQL.String()
		if strings.Contains(sql, "DEFAULT gen_random_uuid()") {
			newSQL := strings.ReplaceAll(sql, "DEFAULT gen_random_uuid()", "")
			d.Statement.SQL.Reset()
			d.Statement.SQL.WriteString(newSQL)
		}
	})

	err = db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Media{},
		&models.Comment{},
		&models.Tag{},
		&models.Vote{},
		&models.RefreshToken{},
		&models.Notification{},
		&models.Bookmark{},
		&models.Follow{},
		&models.AuditLog{},
		&models.Report{},
		&models.ActivityLog{},
		&models.Conversation{},
		&models.Message{},
		&models.Organization{},
		&models.OrgMember{},
		&models.Poll{},
		&models.PollOption{},
		&models.PollVote{},
		&models.AnalyticsEvent{},
	)
	if err != nil {
		t.Fatalf("failed to auto-migrate test database: %v", err)
	}

	return db
}
