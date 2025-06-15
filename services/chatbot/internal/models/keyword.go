package models

import "time"

type Keyword struct {
	ID        int       `json:"id" db:"id"`
	Keyword   string    `json:"keyword" db:"keyword"`
	Response  string    `json:"response" db:"response"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	Priority  int       `json:"priority" db:"priority"`
	MatchType string    `json:"match_type" db:"match_type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	// Enhanced fields from migration 004
	RuleName       string `json:"rule_name" db:"rule_name"`
	HideAfterReply bool   `json:"hide_after_reply" db:"hide_after_reply"`
	SendToInbox    bool   `json:"send_to_inbox" db:"send_to_inbox"`
	InboxMessage   string `json:"inbox_message" db:"inbox_message"`
	InboxImage     string `json:"inbox_image" db:"inbox_image"`
	CreatedBy      *int   `json:"created_by" db:"created_by"`
}

type AutoReplyRequest struct {
	Platform  string `json:"platform"`
	PageID    string `json:"page_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CommentID string `json:"comment_id,omitempty"`
	PostID    string `json:"post_id,omitempty"`
}

type AutoReplyResponse struct {
	ShouldReply      bool    `json:"should_reply"`
	Response         string  `json:"response,omitempty"`
	HasMedia         bool    `json:"has_media,omitempty"`
	MediaDescription *string `json:"media_description,omitempty"`
	MatchedKeyword   string  `json:"matched_keyword,omitempty"`
	MatchType        string  `json:"match_type,omitempty"`
}

// New models for enhanced rules system

type KeywordResponse struct {
	ID           int       `json:"id" db:"id"`
	KeywordID    int       `json:"keyword_id" db:"keyword_id"`
	ResponseText string    `json:"response_text" db:"response_text"`
	ResponseType string    `json:"response_type" db:"response_type"`
	MediaURL     *string   `json:"media_url,omitempty" db:"media_url"`
	Weight       int       `json:"weight" db:"weight"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type KeywordPage struct {
	ID        int       `json:"id" db:"id"`
	KeywordID int       `json:"keyword_id" db:"keyword_id"`
	PageID    int       `json:"page_id" db:"page_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type FallbackResponse struct {
	ID           int       `json:"id" db:"id"`
	PageID       *int      `json:"page_id" db:"page_id"`
	ResponseText string    `json:"response_text" db:"response_text"`
	ResponseType string    `json:"response_type" db:"response_type"`
	MediaURL     *string   `json:"media_url,omitempty" db:"media_url"`
	Weight       int       `json:"weight" db:"weight"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type RateLimit struct {
	ID             int       `json:"id" db:"id"`
	PageID         int       `json:"page_id" db:"page_id"`
	UserID         string    `json:"user_id" db:"user_id"`
	KeywordID      *int      `json:"keyword_id" db:"keyword_id"`
	LastResponseAt time.Time `json:"last_response_at" db:"last_response_at"`
	ResponseCount  int       `json:"response_count" db:"response_count"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type ScheduledMessage struct {
	ID          int       `json:"id" db:"id"`
	PageID      int       `json:"page_id" db:"page_id"`
	UserID      string    `json:"user_id" db:"user_id"`
	MessageText string    `json:"message_text" db:"message_text"`
	MediaURL    *string   `json:"media_url,omitempty" db:"media_url"`
	ScheduledAt time.Time `json:"scheduled_at" db:"scheduled_at"`
	Status      string    `json:"status" db:"status"`
	Attempts    int       `json:"attempts" db:"attempts"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type IncomingMessage struct {
	ID               int       `json:"id" db:"id"`
	PageID           int       `json:"page_id" db:"page_id"`
	UserID           string    `json:"user_id" db:"user_id"`
	MessageText      string    `json:"message_text" db:"message_text"`
	CommentID        string    `json:"comment_id" db:"comment_id"`
	PostID           string    `json:"post_id" db:"post_id"`
	MatchedKeywordID *int      `json:"matched_keyword_id" db:"matched_keyword_id"`
	ResponseSent     bool      `json:"response_sent" db:"response_sent"`
	HiddenComment    bool      `json:"hidden_comment" db:"hidden_comment"`
	SentToInbox      bool      `json:"sent_to_inbox" db:"sent_to_inbox"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

type RuleAnalytics struct {
	ID              int        `json:"id" db:"id"`
	KeywordID       int        `json:"keyword_id" db:"keyword_id"`
	PageID          int        `json:"page_id" db:"page_id"`
	TriggerCount    int        `json:"trigger_count" db:"trigger_count"`
	ResponseCount   int        `json:"response_count" db:"response_count"`
	LastTriggeredAt *time.Time `json:"last_triggered_at" db:"last_triggered_at"`
	Date            time.Time  `json:"date" db:"date"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
}

// Enhanced response for multiple responses support
type EnhancedAutoReplyResponse struct {
	ShouldReply       bool               `json:"should_reply"`
	Response          string             `json:"response,omitempty"`
	Responses         []KeywordResponse  `json:"responses,omitempty"`
	MatchedKeyword    string             `json:"matched_keyword,omitempty"`
	MatchType         string             `json:"match_type,omitempty"`
	HideAfterReply    bool               `json:"hide_after_reply"`
	SendToInbox       bool               `json:"send_to_inbox"`
	InboxMessage      string             `json:"inbox_message,omitempty"`
	InboxImage        string             `json:"inbox_image,omitempty"`
	ScheduledMessages []ScheduledMessage `json:"scheduled_messages,omitempty"`
}

// View model for keywords with responses
type KeywordWithResponses struct {
	ID             int               `json:"id"`
	Keyword        string            `json:"keyword"`
	DefaultMessage string            `json:"default_message"`
	RuleName       string            `json:"rule_name"`
	HideAfterReply bool              `json:"hide_after_reply"`
	SendToInbox    bool              `json:"send_to_inbox"`
	InboxMessage   string            `json:"inbox_message"`
	InboxImage     string            `json:"inbox_image"`
	IsActive       bool              `json:"is_active"`
	Priority       int               `json:"priority"`
	CreatedBy      *int              `json:"created_by"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	Responses      []KeywordResponse `json:"responses"`
	PageIDs        []int             `json:"page_ids"`
}
