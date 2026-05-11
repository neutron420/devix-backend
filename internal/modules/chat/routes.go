package chat

import (
	"devix-backend/internal/middleware"
	jwtpkg "devix-backend/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, jwtManager *jwtpkg.Manager) {
	chat := rg.Group("/chat", middleware.Auth(jwtManager))
	{
		chat.POST("/messages", handler.SendMessage)
		chat.GET("/conversations", handler.GetConversations)
		chat.GET("/conversations/:id/messages", handler.GetMessages)
		chat.PATCH("/conversations/:id/read", handler.MarkAsRead)
		chat.POST("/typing/:id", handler.SendTyping)
	}
}
