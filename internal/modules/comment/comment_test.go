package comment

import (
	"context"
	"testing"

	"devix-backend/internal/models"
	"devix-backend/internal/testutils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupCommentService(t *testing.T) (*Service, *models.User, *models.Post) {
	t.Helper()
	db := testutils.SetupTestDB(t)
	repo := NewRepository(db)
	logger := testutils.NewTestLogger()
	service := NewService(repo, nil, nil, logger)

	// Seed user
	user := &models.User{
		ID:       uuid.New(),
		Username: "commentuser",
		Email:    "commentuser@devix.app",
		Role:     "user",
		IsActive: true,
		PasswordHash: "notused",
	}
	require.NoError(t, db.Create(user).Error)

	// Seed post
	post := &models.Post{
		ID:       uuid.New(),
		AuthorID: user.ID,
		Title:    "Test Post for Comments",
		Slug:     "test-post-comments",
		Content:  "Some content for testing comments.",
		PostType: "concept",
		Status:   "published",
	}
	require.NoError(t, db.Create(post).Error)

	return service, user, post
}

func TestCommentService(t *testing.T) {
	service, user, post := setupCommentService(t)

	var firstCommentID string

	t.Run("create a root comment", func(t *testing.T) {
		req := &CreateCommentRequest{
			Content: "This is a great post!",
		}

		res, err := service.Create(context.Background(), post.ID, user.ID, req)
		require.NoError(t, err)
		assert.Equal(t, post.ID.String(), res.PostID)
		assert.Equal(t, "This is a great post!", res.Content)
		assert.Equal(t, 0, res.Depth)
		assert.Nil(t, res.ParentID)
		firstCommentID = res.ID
	})

	t.Run("create a reply to existing comment", func(t *testing.T) {
		req := &CreateCommentRequest{
			Content:  "I agree with you!",
			ParentID: &firstCommentID,
		}

		res, err := service.Create(context.Background(), post.ID, user.ID, req)
		require.NoError(t, err)
		assert.Equal(t, 1, res.Depth)
		assert.NotNil(t, res.ParentID)
		assert.Equal(t, firstCommentID, *res.ParentID)
	})

	t.Run("create a deeply nested reply (depth 2)", func(t *testing.T) {
		// First get the reply to nest under
		comments, err := service.GetByPostID(context.Background(), post.ID)
		require.NoError(t, err)
		require.NotEmpty(t, comments)
		// The first root comment should have replies
		require.NotEmpty(t, comments[0].Replies)
		replyID := comments[0].Replies[0].ID

		req := &CreateCommentRequest{
			Content:  "Deeply nested reply!",
			ParentID: &replyID,
		}

		res, err := service.Create(context.Background(), post.ID, user.ID, req)
		require.NoError(t, err)
		assert.Equal(t, 2, res.Depth)
	})

	t.Run("reject reply exceeding max depth", func(t *testing.T) {
		// Get the depth-2 comment
		comments, err := service.GetByPostID(context.Background(), post.ID)
		require.NoError(t, err)
		require.NotEmpty(t, comments[0].Replies)
		require.NotEmpty(t, comments[0].Replies[0].Replies)
		deepCommentID := comments[0].Replies[0].Replies[0].ID

		// Create a depth-3 reply
		replyReq := &CreateCommentRequest{
			Content:  "Depth 3 reply",
			ParentID: &deepCommentID,
		}
		depth3, err := service.Create(context.Background(), post.ID, user.ID, replyReq)
		require.NoError(t, err)
		assert.Equal(t, 3, depth3.Depth)

		// Depth-4 should fail (maxCommentDepth is 3)
		req := &CreateCommentRequest{
			Content:  "Too deep!",
			ParentID: &depth3.ID,
		}
		_, err = service.Create(context.Background(), post.ID, user.ID, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Maximum reply depth reached")
	})

	t.Run("reply to non-existent parent fails", func(t *testing.T) {
		fakeParent := uuid.New().String()
		req := &CreateCommentRequest{
			Content:  "Replying to nothing",
			ParentID: &fakeParent,
		}

		_, err := service.Create(context.Background(), post.ID, user.ID, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Parent comment")
	})

	t.Run("list comments returns threaded tree", func(t *testing.T) {
		comments, err := service.GetByPostID(context.Background(), post.ID)
		require.NoError(t, err)

		// Should have 1 root comment
		assert.Equal(t, 1, len(comments))
		assert.Equal(t, 0, comments[0].Depth)
		// Root should have replies
		assert.NotEmpty(t, comments[0].Replies)
	})

	t.Run("list comments for empty post returns empty array", func(t *testing.T) {
		emptyPostID := uuid.New()
		comments, err := service.GetByPostID(context.Background(), emptyPostID)
		require.NoError(t, err)
		assert.Equal(t, 0, len(comments))
	})

	t.Run("update own comment", func(t *testing.T) {
		commentID, _ := uuid.Parse(firstCommentID)
		req := &UpdateCommentRequest{
			Content: "Updated: this is even better!",
		}

		err := service.Update(context.Background(), commentID, user.ID, req)
		assert.NoError(t, err)

		// Verify content changed
		comments, err := service.GetByPostID(context.Background(), post.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated: this is even better!", comments[0].Content)
	})

	t.Run("update non-existent comment", func(t *testing.T) {
		req := &UpdateCommentRequest{
			Content: "Nope",
		}
		err := service.Update(context.Background(), uuid.New(), user.ID, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Comment")
	})

	t.Run("update another user's comment fails", func(t *testing.T) {
		commentID, _ := uuid.Parse(firstCommentID)
		otherUser := uuid.New()
		req := &UpdateCommentRequest{
			Content: "I'm hacking your comment",
		}

		err := service.Update(context.Background(), commentID, otherUser, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only edit your own")
	})

	t.Run("delete own comment (soft delete)", func(t *testing.T) {
		// Create a new comment to delete
		createReq := &CreateCommentRequest{
			Content: "This will be deleted",
		}
		created, err := service.Create(context.Background(), post.ID, user.ID, createReq)
		require.NoError(t, err)

		deleteID, _ := uuid.Parse(created.ID)
		err = service.Delete(context.Background(), deleteID, user.ID, "user")
		assert.NoError(t, err)
	})

	t.Run("delete non-existent comment", func(t *testing.T) {
		err := service.Delete(context.Background(), uuid.New(), user.ID, "user")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Comment")
	})

	t.Run("delete another user's comment fails", func(t *testing.T) {
		commentID, _ := uuid.Parse(firstCommentID)
		err := service.Delete(context.Background(), commentID, uuid.New(), "user")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only delete your own")
	})

	t.Run("admin can delete any comment", func(t *testing.T) {
		createReq := &CreateCommentRequest{
			Content: "Admin will delete this",
		}
		created, err := service.Create(context.Background(), post.ID, user.ID, createReq)
		require.NoError(t, err)

		deleteID, _ := uuid.Parse(created.ID)
		adminID := uuid.New()
		err = service.Delete(context.Background(), deleteID, adminID, "admin")
		assert.NoError(t, err)
	})

	t.Run("moderator can delete any comment", func(t *testing.T) {
		createReq := &CreateCommentRequest{
			Content: "Moderator will delete this",
		}
		created, err := service.Create(context.Background(), post.ID, user.ID, createReq)
		require.NoError(t, err)

		deleteID, _ := uuid.Parse(created.ID)
		modID := uuid.New()
		err = service.Delete(context.Background(), deleteID, modID, "moderator")
		assert.NoError(t, err)
	})

	t.Run("invalid parent ID format", func(t *testing.T) {
		badParent := "not-a-uuid"
		req := &CreateCommentRequest{
			Content:  "Bad parent",
			ParentID: &badParent,
		}

		_, err := service.Create(context.Background(), post.ID, user.ID, req)
		assert.Error(t, err)
	})
}
