package models

import (
	"time"
)

// Keyword represents a keyword and its response configuration
type Keyword struct {
	ID        int       `json:"id" db:"id"`
	Keyword   string    `json:"keyword" db:"keyword"`
	Response  string    `json:"response" db:"response"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	Priority  int       `json:"priority" db:"priority"`
	MatchType string    `json:"match_type" db:"match_type"` // exact, contains, starts_with, ends_with
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Conversation represents a conversation thread
type Conversation struct {
	ID          int       `json:"id" db:"id"`
	PageID      string    `json:"page_id" db:"page_id"`
	UserID      string    `json:"user_id" db:"user_id"`
	CommentID   string    `json:"comment_id,omitempty" db:"comment_id"`
	PostID      string    `json:"post_id,omitempty" db:"post_id"`
	Status      string    `json:"status" db:"status"` // active, closed, archived
	LastMessage time.Time `json:"last_message" db:"last_message"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// AutoReplyRequest represents a request for auto-reply processing
type AutoReplyRequest struct {
	PageID    string `json:"page_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CommentID string `json:"comment_id,omitempty"`
	PostID    string `json:"post_id,omitempty"`
}

// AutoReplyResponse represents the response from auto-reply processing
type AutoReplyResponse struct {
	ShouldReply    bool   `json:"should_reply"`
	Response       string `json:"response,omitempty"`
	MatchedKeyword string `json:"matched_keyword,omitempty"`
	MatchType      string `json:"match_type,omitempty"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}
