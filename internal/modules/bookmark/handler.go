package bookmark

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

func (h *Handler) ToggleBookmark(c *gin.Context) {
	userID := c.MustGet("userId").(uuid.UUID)
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperrors.BadRequest("Invalid post ID"))
		return
	}

	bookmarked, err := h.service.ToggleBookmark(c.Request.Context(), userID, postID)
	if err != nil {
		response.Error(c, apperrors.Internal(err))
		return
	}

	msg := "Post bookmarked"
	if !bookmarked {
		msg = "Bookmark removed"
	}

	response.OK(c, gin.H{"bookmarked": bookmarked, "message": msg})
}

func (h *Handler) ListBookmarks(c *gin.Context) {
	userID := c.MustGet("userId").(uuid.UUID)
	cursor := c.Query("cursor")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	res, err := h.service.ListBookmarks(c.Request.Context(), userID, cursor, limit)
	if err != nil {
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.OK(c, res)
}
