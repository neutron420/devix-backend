package websocket

import (
	"net/http"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/middleware"
	"devix-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, you should check the origin against your allowed origins
		return true 
	},
}

type Handler struct {
	Hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{Hub: hub}
}

// ServeWS handles websocket requests from the peer.
func (h *Handler) ServeWS(c *gin.Context) {
	// Authenticate using the same middleware logic
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required for WebSocket"))
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// Upgrader handles the response on error
		return
	}

	client := &Client{
		Hub:    h.Hub,
		Conn:   conn,
		send:   make(chan []byte, 256),
		UserID: userID,
	}
	client.Hub.register <- client

	// Start pumps in new goroutines
	go client.WritePump()
	go client.ReadPump()
}
