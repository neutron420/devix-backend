package post

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"devix-backend/internal/config"
	"devix-backend/internal/models"
	"devix-backend/internal/modules/media"
	tagmod "devix-backend/internal/modules/tag"
	"devix-backend/internal/modules/search"
	"devix-backend/internal/pkg/cache"
	"devix-backend/internal/testutils"
	"devix-backend/internal/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPostRouter(t *testing.T) (*gin.Engine, string, *models.User) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	validator.Setup()

	db := testutils.SetupTestDB(t)
	logger := testutils.NewTestLogger()

	repo := NewRepository(db)
	mediaRepo := media.NewRepository(db)
	tagRepo := tagmod.NewRepository(db)

	mediaSvc := media.NewService(mediaRepo, nil, config.MediaConfig{}, logger)
	tagSvc := tagmod.NewService(tagRepo, cache.New(nil), logger)
	searchSvc := search.NewService(nil, logger)
	postSvc := NewService(repo, mediaSvc, tagSvc, searchSvc, cache.New(nil), nil, logger)

	handler := NewHandler(postSvc, mediaSvc, nil)
	jwtManager := testutils.NewTestJWT()

	user := &models.User{
		ID:       uuid.New(),
		Username: "author",
		Email:    "author@devix.app",
		Role:     "user",
		IsActive: true,
		PasswordHash: "notused",
	}
	require.NoError(t, db.Create(user).Error)
	token := testutils.AccessTokenFor(t, jwtManager, user)

	// Seed published posts
	posts := []*models.Post{
		{ID: uuid.New(), Title: "First Post Title", Slug: "first-post-title", Content: "First post content for testing.", PostType: "concept", Status: "published", AuthorID: user.ID},
		{ID: uuid.New(), Title: "Second Post Title", Slug: "second-post-title", Content: "Second post content for testing.", PostType: "question", Status: "published", AuthorID: user.ID},
		{ID: uuid.New(), Title: "Draft Post Title", Slug: "draft-post-title", Content: "Draft post content for testing.", PostType: "build-log", Status: "draft", AuthorID: user.ID},
	}
	for _, p := range posts {
		require.NoError(t, db.Create(p).Error)
	}

	router := gin.New()
	rg := router.Group("/api/v1")
	RegisterRoutes(rg, handler, jwtManager)

	return router, token, user
}

func TestPostRoutes(t *testing.T) {
	router, token, _ := setupPostRouter(t)

	t.Run("GET /posts/:slug - Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/posts/first-post-title", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "First Post Title")
		assert.Contains(t, w.Body.String(), `"success":true`)
	})

	t.Run("GET /posts/:slug - Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/posts/non-existent-slug-123", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), `"success":false`)
	})

	t.Run("POST /posts - Create Success", func(t *testing.T) {
		body := `{"title":"Another Post Title","content":"This is the content for another post.","post_type":"concept","status":"published"}`
		req, _ := http.NewRequest("POST", "/api/v1/posts", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "another-post-title")
		assert.Contains(t, w.Body.String(), `"success":true`)
	})

	t.Run("POST /posts - Requires Auth", func(t *testing.T) {
		body := `{"title":"No Auth Post","content":"This should fail without auth.","post_type":"question"}`
		req, _ := http.NewRequest("POST", "/api/v1/posts", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("POST /posts - Validation: Missing Title", func(t *testing.T) {
		body := `{"content":"Content without a title.","post_type":"concept"}`
		req, _ := http.NewRequest("POST", "/api/v1/posts", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST /posts - Validation: Title Too Short", func(t *testing.T) {
		body := `{"title":"Hi","content":"This content is enough.","post_type":"concept"}`
		req, _ := http.NewRequest("POST", "/api/v1/posts", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST /posts - Validation: Content Too Short", func(t *testing.T) {
		body := `{"title":"Valid Title Here","content":"Short","post_type":"concept"}`
		req, _ := http.NewRequest("POST", "/api/v1/posts", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST /posts - Validation: Invalid Post Type", func(t *testing.T) {
		body := `{"title":"Valid Title Here","content":"This content is enough.","post_type":"invalid"}`
		req, _ := http.NewRequest("POST", "/api/v1/posts", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST /posts - Create as Draft", func(t *testing.T) {
		body := `{"title":"Draft Creation Test","content":"This is a draft post for testing.","post_type":"build-log","status":"draft"}`
		req, _ := http.NewRequest("POST", "/api/v1/posts", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), `"status":"draft"`)
	})

	t.Run("GET /posts - List Published", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/posts?limit=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"success":true`)
	})

	t.Run("GET /posts/drafts - Requires Auth", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/posts/drafts", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("GET /posts/drafts - With Auth", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/posts/drafts", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"success":true`)
	})

	t.Run("PUT /posts/:slug - Update Post", func(t *testing.T) {
		body := `{"title":"Updated First Post Title"}`
		req, _ := http.NewRequest("PUT", "/api/v1/posts/first-post-title", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Updated First Post Title")
	})

	t.Run("PUT /posts/:slug - Requires Auth", func(t *testing.T) {
		body := `{"title":"Hacked Title"}`
		req, _ := http.NewRequest("PUT", "/api/v1/posts/first-post-title", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("DELETE /posts/:slug - Requires Auth", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/v1/posts/second-post-title", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("DELETE /posts/:slug - Delete Post", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/v1/posts/second-post-title", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should be 204 or 200 depending on handler
		assert.True(t, w.Code == http.StatusNoContent || w.Code == http.StatusOK)
	})

	t.Run("GET /feed - Feed Alias Works", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/feed?limit=5", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"success":true`)
	})

	t.Run("GET /feed/following - Requires Auth", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/feed/following", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("GET /feed/explore - Public", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/feed/explore?limit=5", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GET /search - Search Alias Works", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/search?q=test&limit=5", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("PATCH /posts/:slug/autosave - Requires Auth", func(t *testing.T) {
		body := `{"title":"Autosaved Title"}`
		req, _ := http.NewRequest("PATCH", "/api/v1/posts/draft-post-title/autosave", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("POST /posts with tags", func(t *testing.T) {
		body := `{"title":"Post With Tags Test","content":"This post has tags attached to it.","post_type":"question","status":"published","tags":["golang","testing"]}`
		req, _ := http.NewRequest("POST", "/api/v1/posts", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "post-with-tags-test")
	})
}
