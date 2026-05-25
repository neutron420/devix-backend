package notification

import (
	"context"
	"testing"

	"devix-backend/internal/models"
	"devix-backend/internal/testutils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupNotificationService(t *testing.T) (*Service, *models.User, *models.User) {
	t.Helper()
	db := testutils.SetupTestDB(t)
	repo := NewRepository(db)
	logger := testutils.NewTestLogger()
	service := NewService(repo, nil, logger)

	// Seed recipient user
	recipient := &models.User{
		ID:           uuid.New(),
		Username:     "recipient",
		Email:        "recipient@devix.app",
		Role:         "user",
		IsActive:     true,
		PasswordHash: "notused",
	}
	require.NoError(t, db.Create(recipient).Error)

	// Seed actor user
	actor := &models.User{
		ID:           uuid.New(),
		Username:     "actor",
		Email:        "actor@devix.app",
		Role:         "user",
		IsActive:     true,
		PasswordHash: "notused",
	}
	require.NoError(t, db.Create(actor).Error)

	return service, recipient, actor
}

func TestNotificationService(t *testing.T) {
	service, recipient, actor := setupNotificationService(t)

	t.Run("create notification", func(t *testing.T) {
		targetID := uuid.New()
		err := service.CreateNotification(context.Background(), recipient.ID, actor.ID, targetID, "commented")
		assert.NoError(t, err)
	})

	t.Run("self-notification is skipped", func(t *testing.T) {
		targetID := uuid.New()
		err := service.CreateNotification(context.Background(), actor.ID, actor.ID, targetID, "commented")
		assert.NoError(t, err)

		// Verify no notification was created for self
		res, err := service.GetUserNotifications(context.Background(), actor.ID, 1, 20)
		require.NoError(t, err)
		assert.Equal(t, 0, len(res.Notifications))
	})

	t.Run("create multiple notifications", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			targetID := uuid.New()
			err := service.CreateNotification(context.Background(), recipient.ID, actor.ID, targetID, "voted")
			require.NoError(t, err)
		}

		res, err := service.GetUserNotifications(context.Background(), recipient.ID, 1, 20)
		require.NoError(t, err)
		// 1 from first test + 5 = 6
		assert.Equal(t, 6, len(res.Notifications))
		assert.Equal(t, int64(6), res.UnreadCount)
		assert.Equal(t, int64(6), res.Total)
	})

	t.Run("list notifications with pagination", func(t *testing.T) {
		// Page 1, limit 3
		res, err := service.GetUserNotifications(context.Background(), recipient.ID, 1, 3)
		require.NoError(t, err)
		assert.Equal(t, 3, len(res.Notifications))
		assert.Equal(t, 1, res.Page)
		assert.Equal(t, 3, res.Limit)
		assert.True(t, res.HasMore)
		assert.Equal(t, int64(6), res.Total)

		// Page 2, limit 3
		res2, err := service.GetUserNotifications(context.Background(), recipient.ID, 2, 3)
		require.NoError(t, err)
		assert.Equal(t, 3, len(res2.Notifications))
		assert.Equal(t, 2, res2.Page)
		assert.False(t, res2.HasMore)

		// Page 3, limit 3 — should be empty
		res3, err := service.GetUserNotifications(context.Background(), recipient.ID, 3, 3)
		require.NoError(t, err)
		assert.Equal(t, 0, len(res3.Notifications))
	})

	t.Run("default page and limit clamping", func(t *testing.T) {
		// Page < 1 defaults to 1
		res, err := service.GetUserNotifications(context.Background(), recipient.ID, 0, 20)
		require.NoError(t, err)
		assert.Equal(t, 1, res.Page)

		// Limit > 50 clamps to 20
		res, err = service.GetUserNotifications(context.Background(), recipient.ID, 1, 999)
		require.NoError(t, err)
		assert.Equal(t, 20, res.Limit)

		// Limit < 1 defaults to 20
		res, err = service.GetUserNotifications(context.Background(), recipient.ID, 1, -5)
		require.NoError(t, err)
		assert.Equal(t, 20, res.Limit)
	})

	t.Run("all notifications start as unread", func(t *testing.T) {
		res, err := service.GetUserNotifications(context.Background(), recipient.ID, 1, 20)
		require.NoError(t, err)
		for _, n := range res.Notifications {
			assert.False(t, n.IsRead)
		}
	})

	t.Run("mark single notification as read", func(t *testing.T) {
		res, err := service.GetUserNotifications(context.Background(), recipient.ID, 1, 1)
		require.NoError(t, err)
		require.NotEmpty(t, res.Notifications)

		notifID, err := uuid.Parse(res.Notifications[0].ID)
		require.NoError(t, err)

		err = service.MarkAsRead(context.Background(), recipient.ID, notifID)
		assert.NoError(t, err)

		// Verify unread count decreased
		res2, err := service.GetUserNotifications(context.Background(), recipient.ID, 1, 20)
		require.NoError(t, err)
		assert.Equal(t, int64(5), res2.UnreadCount)
	})

	t.Run("mark all as read", func(t *testing.T) {
		err := service.MarkAllAsRead(context.Background(), recipient.ID)
		assert.NoError(t, err)

		res, err := service.GetUserNotifications(context.Background(), recipient.ID, 1, 20)
		require.NoError(t, err)
		assert.Equal(t, int64(0), res.UnreadCount)
		for _, n := range res.Notifications {
			assert.True(t, n.IsRead)
		}
	})

	t.Run("empty notifications for user with none", func(t *testing.T) {
		randomUser := uuid.New()
		res, err := service.GetUserNotifications(context.Background(), randomUser, 1, 20)
		require.NoError(t, err)
		assert.Equal(t, 0, len(res.Notifications))
		assert.Equal(t, int64(0), res.UnreadCount)
		assert.Equal(t, int64(0), res.Total)
		assert.False(t, res.HasMore)
	})

	t.Run("notification response includes correct action", func(t *testing.T) {
		targetID := uuid.New()
		_ = service.CreateNotification(context.Background(), recipient.ID, actor.ID, targetID, "followed")

		res, err := service.GetUserNotifications(context.Background(), recipient.ID, 1, 1)
		require.NoError(t, err)
		require.NotEmpty(t, res.Notifications)
		assert.Equal(t, "followed", res.Notifications[0].Action)
		assert.Equal(t, targetID.String(), res.Notifications[0].TargetID)
	})
}
