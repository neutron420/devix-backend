package analytics

import (
	"errors"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetPostStats(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperrors.BadRequest("Invalid post ID"))
		return
	}

	stats, err := h.service.GetPostStats(c.Request.Context(), postID)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.OK(c, stats)
}

func (h *Handler) TrackViewMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		postIDStr := c.Param("id")
		if postIDStr == "" {
			c.Next()
			return
		}

		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			c.Next()
			return
		}

		ua := c.GetHeader("User-Agent")
		ip := c.ClientIP()
		referrer := c.GetHeader("Referer")

		go h.service.TrackView(c.Request.Context(), postID, ua, ip, referrer)
		c.Next()
	}
}
