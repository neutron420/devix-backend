package report

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"devix-backend/internal/models"
	"devix-backend/internal/testutils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReportRoutesAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := testutils.SetupTestDB(t)
	repo := NewRepository(db)
	logger := testutils.NewTestLogger()
	service := NewService(repo, logger)
	handler := NewHandler(service)
	jwtManager := testutils.NewTestJWT()

	// Seed users with different roles
	adminUser := &models.User{
		ID:       uuid.New(),
		Username: "admin",
		Email:    "admin@devix.app",
		Role:     "admin",
		IsActive: true,
	}
	regularUser := &models.User{
		ID:       uuid.New(),
		Username: "regular",
		Email:    "regular@devix.app",
		Role:     "user",
		IsActive: true,
	}
	require.NoError(t, db.Create(adminUser).Error)
	require.NoError(t, db.Create(regularUser).Error)

	adminToken := testutils.AccessTokenFor(t, jwtManager, adminUser)
	regularToken := testutils.AccessTokenFor(t, jwtManager, regularUser)

	router := gin.New()
	rg := router.Group("/api")
	RegisterRoutes(rg, handler, jwtManager)

	t.Run("GET /reports/pending - unauthenticated", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/reports/pending", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("GET /reports/pending - authenticated regular user (Forbidden)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/reports/pending", nil)
		req.Header.Set("Authorization", "Bearer "+regularToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("GET /reports/pending - authenticated admin (OK)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/reports/pending", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("POST /reports - authenticated regular user (OK)", func(t *testing.T) {
		body := `{"target_type":"post", "target_id":"` + uuid.NewString() + `", "reason":"spam", "description":"it is spam"}`
		req, _ := http.NewRequest("POST", "/api/reports", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+regularToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}
