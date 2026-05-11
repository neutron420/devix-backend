package auth

import (
	"errors"

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

func (h *Handler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("Invalid request body: "+err.Error()))
		return
	}

	result, err := h.service.Signup(c.Request.Context(), &req)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.Created(c, result)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("Invalid request body"))
		return
	}

	result, err := h.service.Login(c.Request.Context(), &req)
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

func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("Refresh token is required"))
		return
	}

	result, err := h.service.RefreshTokens(c.Request.Context(), req.RefreshToken)
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

func (h *Handler) Logout(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {

		response.NoContent(c)
		return
	}

	if err := h.service.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		h.service.log.Warn().Err(err).Msg("failed to logout")
	}
	response.NoContent(c)
}

func (h *Handler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("Token is required"))
		return
	}

	if err := h.service.VerifyEmail(c.Request.Context(), req.Token); err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.OK(c, gin.H{"message": "Email verified successfully"})
}

func (h *Handler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("Email is required"))
		return
	}

	if err := h.service.ForgotPassword(c.Request.Context(), req.Email); err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.OK(c, gin.H{"message": "If an account exists with that email, a reset link has been sent"})
}

func (h *Handler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("Invalid request body"))
		return
	}

	if err := h.service.ResetPassword(c.Request.Context(), req.Token, req.Password); err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}

	response.OK(c, gin.H{"message": "Password reset successfully"})
}
