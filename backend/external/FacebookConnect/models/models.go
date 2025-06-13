package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// Conversation represents a conversation record
type Conversation struct {
	ID               int                    `json:"id" db:"id"`
	PageID           string                 `json:"page_id" db:"page_id"`
	FacebookUserID   string                 `json:"facebook_user_id" db:"facebook_user_id"`
	FacebookUserName *string                `json:"facebook_user_name" db:"facebook_user_name"`
	SourceType       string                 `json:"source_type" db:"source_type"` // comment, message
	PostID           *string                `json:"post_id" db:"post_id"`
	StartedAt        time.Time              `json:"started_at" db:"started_at"`
	LastMessageAt    time.Time              `json:"last_message_at" db:"last_message_at"`
	Status           string                 `json:"status" db:"status"` // active, paused, closed
	CreatedAt        time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at" db:"updated_at"`
}

// Message represents a message record
type Message struct {
	ID             int                    `json:"id" db:"id"`
	ConversationID int                    `json:"conversation_id" db:"conversation_id"`
	MessageID      *string                `json:"message_id" db:"message_id"`
	SenderType     string                 `json:"sender_type" db:"sender_type"` // user, bot
	Content        *string                `json:"content" db:"content"`
	MessageType    string                 `json:"message_type" db:"message_type"` // text, image, sticker, video, audio, file
	Metadata       map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt      time.Time              `json:"created_at" db:"created_at"`
	ProcessedAt    *time.Time             `json:"processed_at" db:"processed_at"`
	ResponseTimeMs *int                   `json:"response_time_ms" db:"response_time_ms"`
}

// Keyword represents a keyword record
type Keyword struct {
	ID        int       `json:"id" db:"id"`
	Keyword   string    `json:"keyword" db:"keyword"`
	Response  string    `json:"response" db:"response"`
	MatchType string    `json:"match_type" db:"match_type"` // exact, contains
	Priority  int       `json:"priority" db:"priority"`
	Active    bool      `json:"active" db:"active"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ResponseTemplate represents a response template record
type ResponseTemplate struct {
	ID           int                    `json:"id" db:"id"`
	Name         string                 `json:"name" db:"name"`
	ResponseType string                 `json:"response_type" db:"response_type"` // text, image, button, carousel, quick_reply
	Content      string                 `json:"content" db:"content"`
	Metadata     map[string]interface{} `json:"metadata" db:"metadata"`
	Active       bool                   `json:"active" db:"active"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" db:"updated_at"`
}

// KeywordResponse represents the keyword-response mapping
type KeywordResponse struct {
	KeywordID          int       `json:"keyword_id" db:"keyword_id"`
	ResponseTemplateID int       `json:"response_template_id" db:"response_template_id"`
	OrderIndex         int       `json:"order_index" db:"order_index"`
	Active             bool      `json:"active" db:"active"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
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
