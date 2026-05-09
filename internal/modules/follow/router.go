package follow

import (
	"devix-backend/internal/middleware"
	jwtpkg "devix-backend/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler, jwtManager *jwtpkg.Manager) {
	users := r.Group("/users")
	{
		users.GET("/:username/followers", h.GetFollowers)
		users.GET("/:username/following", h.GetFollowing)
		
		auth := users.Group("")
		auth.Use(middleware.Auth(jwtManager))
		{
			auth.POST("/:username/follow", h.Follow)
			auth.DELETE("/:username/follow", h.Unfollow)
		}
	}
}
