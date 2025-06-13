package common

import "fmt"

// Domain errors that can be shared across services
var (
	ErrNotFound        = NewDomainError("NOT_FOUND", "Resource not found")
	ErrInvalidInput    = NewDomainError("INVALID_INPUT", "Invalid input provided")
	ErrDatabaseError   = NewDomainError("DATABASE_ERROR", "Database operation failed")
	ErrExternalService = NewDomainError("EXTERNAL_SERVICE_ERROR", "External service error")
	ErrUnauthorized    = NewDomainError("UNAUTHORIZED", "Unauthorized access")
	ErrProcessingError = NewDomainError("PROCESSING_ERROR", "Processing failed")
)

// DomainError represents a business domain error
type DomainError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e DomainError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewDomainError creates a new domain error
func NewDomainError(code, message string) DomainError {
	return DomainError{
		Code:    code,
		Message: message,
	}
}

// WithDetails adds details to a domain error
func (e DomainError) WithDetails(details string) DomainError {
	e.Details = details
	return e
}
