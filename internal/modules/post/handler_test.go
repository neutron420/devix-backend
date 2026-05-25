package post

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"devix-backend/internal/models"
	"devix-backend/internal/modules/comment"
	"devix-backend/internal/testutils"
	"devix-backend/internal/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPostAndCommentEndpointsValidateAuthAndIDs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	validator.Setup()

	jwt := testutils.NewTestJWT()
	user := &models.User{ID: uuid.New(), Username: "post-user", Role: "user"}
	token := testutils.AccessTokenFor(t, jwt, user)

	router := gin.New()
	api := router.Group("/api/v1")
	RegisterRoutes(api, NewHandler(nil, nil, nil), jwt)
	comment.RegisterRoutes(api, comment.NewHandler(nil), jwt)

	createBody := `{"title":"How to ship Devix","content":"This is enough useful content.","post_type":"question","status":"published","tags":["go"]}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/posts", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)
	require.Equal(t, http.StatusUnauthorized, createRec.Code, createRec.Body.String())

	commentBody := `{"content":"A useful answer."}`
	commentReq := httptest.NewRequest(http.MethodPost, "/api/v1/posts/not-a-uuid/comments", strings.NewReader(commentBody))
	commentReq.Header.Set("Content-Type", "application/json")
	commentRec := httptest.NewRecorder()
	router.ServeHTTP(commentRec, commentReq)
	require.Equal(t, http.StatusUnauthorized, commentRec.Code, commentRec.Body.String())


	commentsReq := httptest.NewRequest(http.MethodGet, "/api/v1/posts/not-a-uuid/comments", nil)
	commentsRec := httptest.NewRecorder()
	router.ServeHTTP(commentsRec, commentsReq)
	require.Equal(t, http.StatusBadRequest, commentsRec.Code, commentsRec.Body.String())

	deleteCommentReq := httptest.NewRequest(http.MethodDelete, "/api/v1/comments/not-a-uuid", nil)
	deleteCommentReq.Header.Set("Authorization", "Bearer "+token)
	deleteCommentRec := httptest.NewRecorder()
	router.ServeHTTP(deleteCommentRec, deleteCommentReq)
	require.Equal(t, http.StatusBadRequest, deleteCommentRec.Code, deleteCommentRec.Body.String())
}
