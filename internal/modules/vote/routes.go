package vote

import (
	"devix-backend/internal/middleware"
	jwtpkg "devix-backend/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, jwtManager *jwtpkg.Manager) {
	auth := middleware.Auth(jwtManager)
	rg.POST("/posts/:id/vote", auth, handler.VoteOnPost)
	rg.DELETE("/posts/:id/vote", auth, handler.RemovePostVote)
	rg.POST("/comments/:id/vote", auth, handler.VoteOnComment)
}
