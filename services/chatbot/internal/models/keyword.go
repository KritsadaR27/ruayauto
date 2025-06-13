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
}

type AutoReplyRequest struct {
	PageID    string `json:"page_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CommentID string `json:"comment_id,omitempty"`
	PostID    string `json:"post_id,omitempty"`
}

type AutoReplyResponse struct {
	ShouldReply    bool   `json:"should_reply"`
	Response       string `json:"response,omitempty"`
	MatchedKeyword string `json:"matched_keyword,omitempty"`
	MatchType      string `json:"match_type,omitempty"`
}
