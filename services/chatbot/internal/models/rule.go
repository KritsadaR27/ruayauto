package models

import "time"

type Rule struct {
	ID                int       `json:"id" db:"id"`
	Response          string    `json:"response" db:"response"`
	IsActive          bool      `json:"is_active" db:"is_active"`
	Priority          int       `json:"priority" db:"priority"`
	MatchType         string    `json:"match_type" db:"match_type"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
	RuleName          string    `json:"rule_name" db:"rule_name"`
	HideAfterReply    bool      `json:"hide_after_reply" db:"hide_after_reply"`
	SendToInbox       bool      `json:"send_to_inbox" db:"send_to_inbox"`
	InboxMessage      string    `json:"inbox_message" db:"inbox_message"`
	InboxImage        string    `json:"inbox_image" db:"inbox_image"`
	CreatedBy         *int      `json:"created_by" db:"created_by"`
	Keywords          []string  `json:"keywords" db:"-"` // not a DB column, loaded from rule_keywords
}

type RuleResponse struct {
	ID           int       `json:"id" db:"id"`
	RuleID       int       `json:"rule_id" db:"rule_id"`
	ResponseText string    `json:"response_text" db:"response_text"`
	ResponseType string    `json:"response_type" db:"response_type"`
	MediaURL     *string   `json:"media_url,omitempty" db:"media_url"`
	Weight       int       `json:"weight" db:"weight"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
