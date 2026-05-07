package notification

import (
	"devix-backend/internal/middleware"
	jwtpkg "devix-backend/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *Handler, jwtManager *jwtpkg.Manager) {
	notif := r.Group("/notifications")
	notif.Use(middleware.Auth(jwtManager))
	{
		notif.GET("", h.GetNotifications)
		notif.PATCH("/:id/read", h.MarkAsRead)
		notif.PATCH("/read-all", h.MarkAllAsRead)
	}
}
