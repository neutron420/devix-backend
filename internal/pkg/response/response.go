package response

import (
	"net/http"

	apperrors "devix-backend/internal/errors"

	"github.com/gin-gonic/gin"
)

// SuccessResponse is the standard success response envelope.
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// ErrorResponse is the standard error response envelope.
type ErrorResponse struct {
	Success bool             `json:"success"`
	Error   ErrorBody        `json:"error"`
}

// ErrorBody holds the error details.
type ErrorBody struct {
	Code    string                     `json:"code"`
	Message string                     `json:"message"`
	Details []apperrors.ValidationError `json:"details,omitempty"`
}

// Meta holds pagination metadata.
type Meta struct {
	Cursor  string `json:"cursor,omitempty"`
	HasMore bool   `json:"has_more"`
	Total   *int64 `json:"total,omitempty"`
}

// OK sends a 200 response with data.
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

// OKWithMeta sends a 200 response with data and pagination metadata.
func OKWithMeta(c *gin.Context, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// Created sends a 201 response with data.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

// NoContent sends a 204 response.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error sends an error response based on an AppError.
func Error(c *gin.Context, err *apperrors.AppError) {
	c.JSON(err.StatusCode, ErrorResponse{
		Success: false,
		Error: ErrorBody{
			Code:    err.Code,
			Message: err.Message,
			Details: err.Details,
		},
	})
}

// Abort sends an error response and aborts the middleware chain.
func Abort(c *gin.Context, err *apperrors.AppError) {
	c.AbortWithStatusJSON(err.StatusCode, ErrorResponse{
		Success: false,
		Error: ErrorBody{
			Code:    err.Code,
			Message: err.Message,
			Details: err.Details,
		},
	})
}
