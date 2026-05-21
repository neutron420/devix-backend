package vote

import (
	"devix-backend/internal/middleware"
	jwtpkg "devix-backend/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, jwtManager *jwtpkg.Manager) {
	auth := middleware.Auth(jwtManager)
	rg.POST("/posts/:slug/vote", auth, handler.VoteOnPost)
	rg.DELETE("/posts/:slug/vote", auth, handler.RemovePostVote)
	rg.POST("/comments/:id/vote", auth, handler.VoteOnComment)
	rg.DELETE("/comments/:id/vote", auth, handler.RemoveCommentVote)
}
