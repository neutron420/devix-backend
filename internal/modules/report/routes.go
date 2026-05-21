package report

import (
	"devix-backend/internal/middleware"
	jwtpkg "devix-backend/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, jwtManager *jwtpkg.Manager) {
	reports := rg.Group("/reports")
	reports.Use(middleware.Auth(jwtManager))
	{
		reports.POST("", handler.Create)
		reports.GET("/pending", middleware.RequireRole("admin", "moderator"), handler.ListPending)
		reports.GET("", middleware.RequireRole("admin", "moderator"), handler.ListAll)
		reports.PATCH("/:id", middleware.RequireRole("admin", "moderator"), handler.Review)
	}
}
