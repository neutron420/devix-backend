package post

import (
	"devix-backend/internal/middleware"
	jwtpkg "devix-backend/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, jwtManager *jwtpkg.Manager) {
	posts := rg.Group("/posts")
	{
		posts.GET("", handler.List)
		posts.GET("/:id", handler.GetBySlug)
		posts.POST("", middleware.Auth(jwtManager), handler.Create)
		posts.PUT("/:id", middleware.Auth(jwtManager), handler.Update)
		posts.DELETE("/:id", middleware.Auth(jwtManager), handler.Delete)
		posts.POST("/:id/media", middleware.Auth(jwtManager), handler.UploadMedia)
	}
	rg.GET("/feed", handler.List)
}
