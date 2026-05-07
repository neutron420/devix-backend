package notification

import (
	"strconv"

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

func (h *Handler) GetNotifications(c *gin.Context) {
	userID := c.MustGet("userId").(uuid.UUID)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	res, err := h.service.GetUserNotifications(c.Request.Context(), userID, page, limit)
	if err != nil {
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.OK(c, res)
}

func (h *Handler) MarkAsRead(c *gin.Context) {
	userID := c.MustGet("userId").(uuid.UUID)
	notificationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperrors.BadRequest("Invalid notification ID"))
		return
	}

	if err := h.service.MarkAsRead(c.Request.Context(), userID, notificationID); err != nil {
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.OK(c, gin.H{"message": "Notification marked as read"})
}

func (h *Handler) MarkAllAsRead(c *gin.Context) {
	userID := c.MustGet("userId").(uuid.UUID)

	if err := h.service.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.OK(c, gin.H{"message": "All notifications marked as read"})
}
