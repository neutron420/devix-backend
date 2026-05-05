package response

import (
	"net/http"

	apperrors "devix-backend/internal/errors"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Success bool      `json:"success"`
	Error   ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string                      `json:"code"`
	Message string                      `json:"message"`
	Details []apperrors.ValidationError `json:"details,omitempty"`
}

type Meta struct {
	Cursor  string `json:"cursor,omitempty"`
	HasMore bool   `json:"has_more"`
	Total   *int64 `json:"total,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

func OKWithMeta(c *gin.Context, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

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
