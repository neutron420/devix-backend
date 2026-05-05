package apperrors

import (
	"fmt"
	"net/http"
)

const (
	CodeNotFound        = "NOT_FOUND"
	CodeUnauthorized    = "UNAUTHORIZED"
	CodeForbidden       = "FORBIDDEN"
	CodeValidation      = "VALIDATION_ERROR"
	CodeConflict        = "CONFLICT"
	CodeInternal        = "INTERNAL_ERROR"
	CodeBadRequest      = "BAD_REQUEST"
	CodeTooManyReqs     = "TOO_MANY_REQUESTS"
	CodePayloadTooLarge = "PAYLOAD_TOO_LARGE"
)

type AppError struct {
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	StatusCode int               `json:"-"`
	Details    []ValidationError `json:"details,omitempty"`
	Err        error             `json:"-"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NotFound(resource string) *AppError {
	return &AppError{
		Code:       CodeNotFound,
		Message:    fmt.Sprintf("%s not found", resource),
		StatusCode: http.StatusNotFound,
	}
}

func Unauthorized(msg string) *AppError {
	return &AppError{
		Code:       CodeUnauthorized,
		Message:    msg,
		StatusCode: http.StatusUnauthorized,
	}
}

func Forbidden(msg string) *AppError {
	return &AppError{
		Code:       CodeForbidden,
		Message:    msg,
		StatusCode: http.StatusForbidden,
	}
}

func BadRequest(msg string) *AppError {
	return &AppError{
		Code:       CodeBadRequest,
		Message:    msg,
		StatusCode: http.StatusBadRequest,
	}
}

func Validation(details []ValidationError) *AppError {
	return &AppError{
		Code:       CodeValidation,
		Message:    "Validation failed",
		StatusCode: http.StatusUnprocessableEntity,
		Details:    details,
	}
}

func Conflict(msg string) *AppError {
	return &AppError{
		Code:       CodeConflict,
		Message:    msg,
		StatusCode: http.StatusConflict,
	}
}

func Internal(err error) *AppError {
	return &AppError{
		Code:       CodeInternal,
		Message:    "An internal error occurred",
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

func TooManyRequests() *AppError {
	return &AppError{
		Code:       CodeTooManyReqs,
		Message:    "Too many requests, please try again later",
		StatusCode: http.StatusTooManyRequests,
	}
}

func PayloadTooLarge(msg string) *AppError {
	return &AppError{
		Code:       CodePayloadTooLarge,
		Message:    msg,
		StatusCode: http.StatusRequestEntityTooLarge,
	}
}

func Wrap(err error, code string, msg string, status int) *AppError {
	return &AppError{
		Code:       code,
		Message:    msg,
		StatusCode: status,
		Err:        err,
	}
}
