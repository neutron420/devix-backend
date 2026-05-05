package user

import (
	"errors"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/middleware"
	"devix-backend/internal/modules/media"
	"devix-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service      *Service
	mediaService *media.Service
}

func NewHandler(service *Service, mediaService *media.Service) *Handler {
	return &Handler{
		service:      service,
		mediaService: mediaService,
	}
}

func (h *Handler) GetMe(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	profile, err := h.service.GetMyProfile(c.Request.Context(), userID)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.OK(c, profile)
}

func (h *Handler) UpdateMe(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("Invalid request body: "+err.Error()))
		return
	}

	profile, err := h.service.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.OK(c, profile)
}

func (h *Handler) UpdateAvatar(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		response.Error(c, apperrors.BadRequest("Avatar file is required"))
		return
	}
	defer file.Close()

	url, err := h.mediaService.UploadAvatar(c.Request.Context(), file, header)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}

	profile, err := h.service.UpdateAvatar(c.Request.Context(), userID, url)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.OK(c, profile)
}

func (h *Handler) GetPublicProfile(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		response.Error(c, apperrors.BadRequest("Username is required"))
		return
	}

	profile, err := h.service.GetPublicProfile(c.Request.Context(), username)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.OK(c, profile)
}
