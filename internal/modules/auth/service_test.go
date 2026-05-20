package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"devix-backend/internal/validator"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestAuthEndpointsValidateRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	validator.Setup()

	router := gin.New()
	api := router.Group("/api/v1")
	RegisterRoutes(api, NewHandler(nil))

	signupBody := `{"username":"no","email":"bad","password":"short"}`
	signupReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", strings.NewReader(signupBody))
	signupReq.Header.Set("Content-Type", "application/json")
	signupRec := httptest.NewRecorder()
	router.ServeHTTP(signupRec, signupReq)
	require.Equal(t, http.StatusBadRequest, signupRec.Code, signupRec.Body.String())

	loginBody := `{"email":"bad"}`
	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	router.ServeHTTP(loginRec, loginReq)
	require.Equal(t, http.StatusBadRequest, loginRec.Code, loginRec.Body.String())

	refreshReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", strings.NewReader(`{}`))
	refreshReq.Header.Set("Content-Type", "application/json")
	refreshRec := httptest.NewRecorder()
	router.ServeHTTP(refreshRec, refreshReq)
	require.Equal(t, http.StatusBadRequest, refreshRec.Code, refreshRec.Body.String())

	logoutReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", strings.NewReader(`{}`))
	logoutReq.Header.Set("Content-Type", "application/json")
	logoutRec := httptest.NewRecorder()
	router.ServeHTTP(logoutRec, logoutReq)
	require.Equal(t, http.StatusNoContent, logoutRec.Code, logoutRec.Body.String())
}
