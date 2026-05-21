package websocket

import (
	"net/http"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/middleware"
	"devix-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	Hub      *Hub
	upgrader websocket.Upgrader
}

func NewHandler(hub *Hub, allowedOrigins []string) *Handler {
	return &Handler{
		Hub: hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				if origin == "" {
					return true
				}
				for _, allowed := range allowedOrigins {
					if allowed == "*" || origin == allowed {
						return true
					}
				}
				return false
			},
		},
	}
}

func (h *Handler) ServeWS(c *gin.Context) {

	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required for WebSocket"))
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {

		return
	}

	client := &Client{
		Hub:    h.Hub,
		Conn:   conn,
		send:   make(chan []byte, 256),
		UserID: userID,
	}
	client.Hub.register <- client

	go client.WritePump()
	go client.ReadPump()
}
