package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// Keyword represents a keyword-response mapping
type Keyword struct {
	ID        int       `json:"id" db:"id"`
	Keyword   string    `json:"keyword" db:"keyword"`
	Response  string    `json:"response" db:"response"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// KeywordPair represents a pair of keywords and responses for bulk operations
type KeywordPair struct {
	Keywords  []string `json:"keywords"`
	Responses []string `json:"responses"`
}

// MessageAnalytics represents a message interaction record
type MessageAnalytics struct {
	ID                  int       `json:"id" db:"id"`
	FacebookMessageID   *string   `json:"facebook_message_id" db:"facebook_message_id"`
	SenderID            string    `json:"sender_id" db:"sender_id"`
	RecipientID         string    `json:"recipient_id" db:"recipient_id"`
	MessageText         *string   `json:"message_text" db:"message_text"`
	MatchedKeyword      *string   `json:"matched_keyword" db:"matched_keyword"`
	ResponseSent        *string   `json:"response_sent" db:"response_sent"`
	ResponseTimeMs      *int      `json:"response_time_ms" db:"response_time_ms"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}

// WebhookLog represents a webhook request/response log
type WebhookLog struct {
	ID               int             `json:"id" db:"id"`
	WebhookType      string          `json:"webhook_type" db:"webhook_type"`
	RequestBody      json.RawMessage `json:"request_body" db:"request_body"`
	ResponseStatus   *int            `json:"response_status" db:"response_status"`
	ResponseBody     *string         `json:"response_body" db:"response_body"`
	ProcessingTimeMs *int            `json:"processing_time_ms" db:"processing_time_ms"`
	ErrorMessage     *string         `json:"error_message" db:"error_message"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
}

// CommentData represents processed comment data
type CommentData struct {
	CommentID string `json:"comment_id"`
	Message   string `json:"message"`
	SenderID  string `json:"sender_id"`
	PostID    string `json:"post_id"`
}

// AutoReplyConfig holds configuration for auto-reply behavior
type AutoReplyConfig struct {
	EnableDefault      bool     `json:"enableDefault"`
	DefaultResponses   []string `json:"defaultResponses"`
	CaseSensitive      bool     `json:"caseSensitive"`
	EnableFuzzyMatch   bool     `json:"enableFuzzyMatch"`
	MinMatchThreshold  float64  `json:"minMatchThreshold"`
	NoTag              bool     `json:"noTag"`
	NoSticker          bool     `json:"noSticker"`
}

// MatchResult represents the result of keyword matching
type MatchResult struct {
	Matched          bool    `json:"matched"`
	MatchedKeyword   string  `json:"matchedKeyword,omitempty"`
	Response         string  `json:"response,omitempty"`
	ConfidenceScore  float64 `json:"confidenceScore,omitempty"`
	MatchType        string  `json:"matchType,omitempty"` // "exact", "partial", "fuzzy", "default"
}

// APIResponse represents standardized API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// APIError represents API error structure
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Meta represents metadata for API responses
type Meta struct {
	Timestamp time.Time `json:"timestamp"`
	RequestID string    `json:"request_id,omitempty"`
	Version   string    `json:"version,omitempty"`
}

// FacebookError represents Facebook API error structure
type FacebookError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (e FacebookError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Type, e.Code, e.Message)
}

// HealthStatus represents service health status
type HealthStatus struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Services  map[string]interface{} `json:"services"`
}
