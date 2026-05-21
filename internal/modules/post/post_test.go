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

func TestPostRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	validator.Setup()

	db := testutils.SetupTestDB(t)
	logger := testutils.NewTestLogger()

	// Repositories
	repo := NewRepository(db)
	mediaRepo := media.NewRepository(db)
	tagRepo := tagmod.NewRepository(db)

	// Services
	mediaSvc := media.NewService(mediaRepo, nil, config.MediaConfig{}, logger)
	tagSvc := tagmod.NewService(tagRepo, cache.New(nil), logger)
	searchSvc := search.NewService(nil, logger)
	postSvc := NewService(repo, mediaSvc, tagSvc, searchSvc, cache.New(nil), nil, logger)

	handler := NewHandler(postSvc, mediaSvc, nil)
	jwtManager := testutils.NewTestJWT()

	// Seed user
	user := &models.User{
		ID:       uuid.New(),
		Username: "author",
		Email:    "author@devix.app",
		Role:     "user",
		IsActive: true,
	}
	require.NoError(t, db.Create(user).Error)
	token := testutils.AccessTokenFor(t, jwtManager, user)

	// Seed a post
	post := &models.Post{
		ID:        uuid.New(),
		Title:     "Testing slugs in Devix",
		Slug:      "testing-slugs-in-devix",
		Content:   "This is testing content. It must be sufficiently long.",
		PostType:  "concept",
		Status:    "published",
		AuthorID:  user.ID,
	}
	require.NoError(t, db.Create(post).Error)

	router := gin.New()
	rg := router.Group("/api/v1")
	RegisterRoutes(rg, handler, jwtManager)

	t.Run("GET /posts/:slug - Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/posts/testing-slugs-in-devix", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Testing slugs in Devix")
	})

	t.Run("GET /posts/:slug - Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/posts/non-existent-slug-123", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
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
	})
}
