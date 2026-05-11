package poll

import (
	"devix-backend/internal/middleware"
	jwtpkg "devix-backend/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, jwtManager *jwtpkg.Manager) {
	polls := rg.Group("/polls")
	{
		polls.POST("", middleware.Auth(jwtManager), handler.CreatePoll)
		polls.POST("/:id/vote", middleware.Auth(jwtManager), handler.Vote)
	}
}
