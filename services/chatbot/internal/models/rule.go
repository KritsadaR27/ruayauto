package models

import "time"

type Rule struct {
	ID                int       `json:"id" db:"id"`
	Rule              string    `json:"rule" db:"rule"`
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
}
