package follow

import (
	"errors"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/middleware"
	"devix-backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Follow(c *gin.Context) {
	followerID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required"))
		return
	}
	targetUsername := c.Param("username")

	if err := h.service.Follow(c.Request.Context(), followerID, targetUsername); err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.OK(c, gin.H{"message": "Successfully followed user"})
}

func (h *Handler) Unfollow(c *gin.Context) {
	followerID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required"))
		return
	}
	targetUsername := c.Param("username")

	if err := h.service.Unfollow(c.Request.Context(), followerID, targetUsername); err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.OK(c, gin.H{"message": "Successfully unfollowed user"})
}

func (h *Handler) GetFollowers(c *gin.Context) {
	username := c.Param("username")

	followers, err := h.service.GetFollowers(c.Request.Context(), username)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.OK(c, followers)
}

func (h *Handler) GetFollowing(c *gin.Context) {
	username := c.Param("username")

	following, err := h.service.GetFollowing(c.Request.Context(), username)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.OK(c, following)
}
