package errors

import (
	"fmt"
	"net/http"
)

// AppError represents application error with HTTP status
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Common errors
var (
	ErrNotFound     = &AppError{Code: "NOT_FOUND", Message: "Resource not found", StatusCode: http.StatusNotFound}
	ErrBadRequest   = &AppError{Code: "BAD_REQUEST", Message: "Invalid request", StatusCode: http.StatusBadRequest}
	ErrUnauthorized = &AppError{Code: "UNAUTHORIZED", Message: "Unauthorized", StatusCode: http.StatusUnauthorized}
	ErrInternal     = &AppError{Code: "INTERNAL_ERROR", Message: "Internal server error", StatusCode: http.StatusInternalServerError}
)

// New creates a new AppError
func New(code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// Wrap wraps an error with AppError
func Wrap(err error, code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}
