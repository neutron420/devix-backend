package org

import (
	"errors"
	"strconv"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/middleware"
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

func (h *Handler) CreateOrg(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	var req CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("Invalid request body"))
		return
	}

	org, err := h.service.CreateOrg(c.Request.Context(), userID, &req)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.Created(c, org)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	res, err := h.service.List(c.Request.Context(), page, limit)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.OK(c, res)
}

func (h *Handler) GetMembers(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperrors.BadRequest("Invalid organization ID"))
		return
	}

	members, err := h.service.GetMembers(c.Request.Context(), orgID)
	if err != nil {
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.OK(c, members)
}

func (h *Handler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.Error(c, apperrors.BadRequest("Slug is required"))
		return
	}

	org, err := h.service.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.OK(c, org)
}

func (h *Handler) UpdateOrg(c *gin.Context) {
	actorID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperrors.BadRequest("Invalid organization ID"))
		return
	}

	var req UpdateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("Invalid request body"))
		return
	}

	org, err := h.service.UpdateOrg(c.Request.Context(), actorID, orgID, &req)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.OK(c, org)
}

func (h *Handler) DeleteOrg(c *gin.Context) {
	actorID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperrors.BadRequest("Invalid organization ID"))
		return
	}

	if err := h.service.DeleteOrg(c.Request.Context(), actorID, orgID); err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.NoContent(c)
}

func (h *Handler) GetByID(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperrors.BadRequest("Invalid organization ID"))
		return
	}

	org, err := h.service.GetByID(c.Request.Context(), orgID)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.OK(c, org)
}

func (h *Handler) AddMember(c *gin.Context) {
	actorID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperrors.BadRequest("Invalid organization ID"))
		return
	}

	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("Invalid request body"))
		return
	}

	if err := h.service.AddMember(c.Request.Context(), actorID, orgID, &req); err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.NoContent(c)
}

func (h *Handler) RemoveMember(c *gin.Context) {
	actorID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required"))
		return
	}

	orgID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperrors.BadRequest("Invalid organization ID"))
		return
	}

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		response.Error(c, apperrors.BadRequest("Invalid user ID"))
		return
	}

	if err := h.service.RemoveMember(c.Request.Context(), actorID, orgID, userID); err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.NoContent(c)
}
