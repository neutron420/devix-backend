package tag

import (
	"errors"
	"strconv"
	"strings"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAll(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	if q == "" {
		q = strings.TrimSpace(c.Query("query"))
	}
	if q != "" {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
		result, err := h.service.Search(c.Request.Context(), q, limit)
		if err != nil {
			var appErr *apperrors.AppError
			if errors.As(err, &appErr) {
				response.Error(c, appErr)
				return
			}
			response.Error(c, apperrors.Internal(err))
			return
		}
		response.OK(c, result)
		return
	}

	result, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}
	response.OK(c, result)
}

func (h *Handler) GetTrending(c *gin.Context) {
	result, err := h.service.GetTrending(c.Request.Context(), 10)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}
	response.OK(c, result)
}

func (h *Handler) GetTree(c *gin.Context) {
	result, err := h.service.GetTagTree(c.Request.Context())
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}
	response.OK(c, result)
}

func (h *Handler) GetByCategory(c *gin.Context) {
	category := c.Param("category")
	result, err := h.service.GetByCategory(c.Request.Context(), category)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}
	response.OK(c, result)
}
