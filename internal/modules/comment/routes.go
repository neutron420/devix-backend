package comment

import (
	"devix-backend/internal/middleware"
	jwtpkg "devix-backend/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, jwtManager *jwtpkg.Manager) {
	rg.GET("/posts/:slug/comments", handler.GetByPostID)
	rg.POST("/posts/:slug/comments", middleware.Auth(jwtManager), handler.Create)
	comments := rg.Group("/comments", middleware.Auth(jwtManager))
	{
		comments.PUT("/:id", handler.Update)
		comments.DELETE("/:id", handler.Delete)
	}
}
