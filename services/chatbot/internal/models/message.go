package models

import "time"

type Message struct {
	ID          int       `json:"id" db:"id"`
	Platform    string    `json:"platform" db:"platform"`
	UserID      string    `json:"user_id" db:"user_id"`
	PageID      string    `json:"page_id" db:"page_id"`
	Content     string    `json:"content" db:"content"`
	CommentID   string    `json:"comment_id,omitempty" db:"comment_id"`
	PostID      string    `json:"post_id,omitempty" db:"post_id"`
	MessageType string    `json:"message_type" db:"message_type"` // comment, inbox, post
	Processed   bool      `json:"processed" db:"processed"`
	Response    string    `json:"response,omitempty" db:"response"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
