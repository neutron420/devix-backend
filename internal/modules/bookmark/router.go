package bookmark

import (
	"devix-backend/internal/middleware"
	jwtpkg "devix-backend/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler, jwtManager *jwtpkg.Manager) {
	bookmarks := r.Group("/bookmarks")
	bookmarks.Use(middleware.Auth(jwtManager))
	{
		bookmarks.GET("", h.ListBookmarks)
	}

	// Toggle bookmark is part of posts group
	posts := r.Group("/posts")
	posts.Use(middleware.Auth(jwtManager))
	{
		posts.POST("/:id/bookmark", h.ToggleBookmark)
		posts.DELETE("/:id/bookmark", h.ToggleBookmark) // Both POST and DELETE can toggle
	}
}
