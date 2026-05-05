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

		return true
	},
}

type Handler struct {
	Hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{Hub: hub}
}

func (h *Handler) ServeWS(c *gin.Context) {

	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required for WebSocket"))
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
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
