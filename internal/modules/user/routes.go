package user

import (
	"devix-backend/internal/middleware"
	jwtpkg "devix-backend/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, jwtManager *jwtpkg.Manager) {
	users := rg.Group("/users")
	{

		me := users.Group("/me", middleware.Auth(jwtManager))
		{
			me.GET("", handler.GetMe)
			me.PUT("", handler.UpdateMe)
			me.PUT("/avatar", handler.UpdateAvatar)
		}

		users.GET("/:username", handler.GetPublicProfile)
	}
}
