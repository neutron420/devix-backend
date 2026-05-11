package analytics

import (
	"devix-backend/internal/middleware"
	jwtpkg "devix-backend/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, jwtManager *jwtpkg.Manager) {
	analytics := rg.Group("/analytics", middleware.Auth(jwtManager))
	{
		analytics.GET("/posts/:id", handler.GetPostStats)
	}
}
