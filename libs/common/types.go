package common

import "time"

// APIResponse represents standardized API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// APIError represents API error details
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Meta represents response metadata
type Meta struct {
	RequestID string    `json:"request_id,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Page      int       `json:"page,omitempty"`
	PerPage   int       `json:"per_page,omitempty"`
	Total     int       `json:"total,omitempty"`
}

// ServiceInfo represents service information
type ServiceInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Port    int    `json:"port"`
	Health  string `json:"health"`
}
