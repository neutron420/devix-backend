package org

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"devix-backend/internal/models"
	"devix-backend/internal/testutils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func setupOrgTestRouter(t *testing.T) (*gin.Engine, string) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	jwt := testutils.NewTestJWT()
	owner := &models.User{ID: uuid.New(), Username: "owner", Role: "user"}

	r := gin.New()
	api := r.Group("/api/v1")
	RegisterRoutes(api, NewHandler(nil), jwt)

	return r, testutils.AccessTokenFor(t, jwt, owner)
}

func TestOrganizationEndpointsValidateAuthAndIDs(t *testing.T) {
	router, token := setupOrgTestRouter(t)

	createBody := `{"name":"Devix Labs","bio":"Builder org"}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/organizations", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)
	require.Equal(t, http.StatusUnauthorized, createRec.Code, createRec.Body.String())

	updateBody := `{"name":"Devix Core","bio":"Updated"}`
	updateReq := httptest.NewRequest(http.MethodPut, "/api/v1/organizations/not-a-uuid", strings.NewReader(updateBody))
	updateReq.Header.Set("Authorization", "Bearer "+token)
	updateReq.Header.Set("Content-Type", "application/json")
	updateRec := httptest.NewRecorder()
	router.ServeHTTP(updateRec, updateReq)
	require.Equal(t, http.StatusBadRequest, updateRec.Code, updateRec.Body.String())

	membersReq := httptest.NewRequest(http.MethodGet, "/api/v1/organizations/not-a-uuid/members", nil)
	membersRec := httptest.NewRecorder()
	router.ServeHTTP(membersRec, membersReq)
	require.Equal(t, http.StatusBadRequest, membersRec.Code, membersRec.Body.String())

	removeReq := httptest.NewRequest(http.MethodDelete, "/api/v1/organizations/not-a-uuid/members/"+uuid.NewString(), nil)
	removeReq.Header.Set("Authorization", "Bearer "+token)
	removeRec := httptest.NewRecorder()
	router.ServeHTTP(removeRec, removeReq)
	require.Equal(t, http.StatusBadRequest, removeRec.Code, removeRec.Body.String())
}
