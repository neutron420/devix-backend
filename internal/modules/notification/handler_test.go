package notification

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"devix-backend/internal/models"
	"devix-backend/internal/testutils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNotificationEndpointsValidateAuthAndIDs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwt := testutils.NewTestJWT()
	user := &models.User{ID: uuid.New(), Username: "notify-user", Role: "user"}
	router := gin.New()
	api := router.Group("/api/v1")
	RegisterRoutes(api, NewHandler(nil), jwt)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/notifications?page=2&limit=2", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusUnauthorized, rec.Code, rec.Body.String())

	readReq := httptest.NewRequest(http.MethodPatch, "/api/v1/notifications/not-a-uuid/read", nil)
	readReq.Header.Set("Authorization", "Bearer "+testutils.AccessTokenFor(t, jwt, user))
	readRec := httptest.NewRecorder()
	router.ServeHTTP(readRec, readReq)
	require.Equal(t, http.StatusBadRequest, readRec.Code, readRec.Body.String())
}
