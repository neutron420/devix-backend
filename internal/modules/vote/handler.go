package vote

import (
	"errors"

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

func (h *Handler) VoteOnPost(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required"))
		return
	}
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperrors.BadRequest("Invalid post ID"))
		return
	}
	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("vote_type must be 1 (upvote) or -1 (downvote)"))
		return
	}
	newCount, err := h.service.VoteOnPost(c.Request.Context(), userID, postID, req.VoteType)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}
	response.OK(c, gin.H{"vote_type": req.VoteType, "new_count": newCount})
}

func (h *Handler) VoteOnComment(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required"))
		return
	}
	commentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperrors.BadRequest("Invalid comment ID"))
		return
	}
	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("vote_type must be 1 (upvote) or -1 (downvote)"))
		return
	}
	newCount, err := h.service.VoteOnComment(c.Request.Context(), userID, commentID, req.VoteType)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}
	response.OK(c, gin.H{"vote_type": req.VoteType, "new_count": newCount})
}

func (h *Handler) RemoveCommentVote(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required"))
		return
	}
	commentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperrors.BadRequest("Invalid comment ID"))
		return
	}
	newCount, err := h.service.RemoveCommentVote(c.Request.Context(), userID, commentID)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}
	response.OK(c, gin.H{"new_count": newCount})
}

func (h *Handler) RemovePostVote(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required"))
		return
	}
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperrors.BadRequest("Invalid post ID"))
		return
	}
	newCount, err := h.service.RemovePostVote(c.Request.Context(), userID, postID)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}
	response.OK(c, gin.H{"new_count": newCount})
}
