package models

import (
	"time"
)

// SimpleKeyword represents a simple keyword-response pair with matching options
type SimpleKeyword struct {
	ID        int       `json:"id" db:"id"`
	Keyword   string    `json:"keyword" db:"keyword"`
	Response  string    `json:"response" db:"response"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	MatchType string    `json:"match_type" db:"match_type"` // "exact", "contains", "starts_with", "ends_with"
	Priority  int       `json:"priority" db:"priority"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}