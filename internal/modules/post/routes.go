package post

import (
	"devix-backend/internal/middleware"
	jwtpkg "devix-backend/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, jwtManager *jwtpkg.Manager) {
	posts := rg.Group("/posts")
	posts.Use(middleware.OptionalAuth(jwtManager))
	{
		posts.GET("", handler.List)
		posts.GET("/:id", handler.GetBySlug)
		posts.GET("/drafts", middleware.Auth(jwtManager), handler.ListDrafts)
		posts.POST("", middleware.Auth(jwtManager), handler.Create)
		posts.PUT("/:id", middleware.Auth(jwtManager), handler.Update)
		posts.PATCH("/:id/autosave", middleware.Auth(jwtManager), handler.Autosave)
		posts.DELETE("/:id", middleware.Auth(jwtManager), handler.Delete)
		posts.POST("/:id/media", middleware.Auth(jwtManager), handler.UploadMedia)
	}
	rg.GET("/feed", middleware.OptionalAuth(jwtManager), handler.List)
	rg.GET("/feed/following", middleware.Auth(jwtManager), handler.ListFollowing)
	rg.GET("/feed/explore", middleware.OptionalAuth(jwtManager), handler.ListExplore)
}
