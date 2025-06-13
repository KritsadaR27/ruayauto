package models

import (
	"time"
)

// Conversation represents a conversation between a user and the bot
type Conversation struct {
	ID               int       `json:"id" db:"id"`
	PageID           string    `json:"page_id" db:"page_id"`
	FacebookUserID   string    `json:"facebook_user_id" db:"facebook_user_id"`
	FacebookUserName *string   `json:"facebook_user_name" db:"facebook_user_name"`
	SourceType       string    `json:"source_type" db:"source_type"` // comment, message
	PostID           *string   `json:"post_id" db:"post_id"`
	StartedAt        time.Time `json:"started_at" db:"started_at"`
	LastMessageAt    time.Time `json:"last_message_at" db:"last_message_at"`
	Status           string    `json:"status" db:"status"` // active, paused, closed
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// Message represents a single message in a conversation
type Message struct {
	ID             int        `json:"id" db:"id"`
	ConversationID int        `json:"conversation_id" db:"conversation_id"`
	MessageID      *string    `json:"message_id" db:"message_id"` // Facebook message ID
	SenderType     string     `json:"sender_type" db:"sender_type"` // user, bot
	Content        *string    `json:"content" db:"content"`
	MessageType    string     `json:"message_type" db:"message_type"` // text, image, sticker, video, audio, file
	Metadata       string     `json:"metadata" db:"metadata"` // JSON string
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	ProcessedAt    *time.Time `json:"processed_at" db:"processed_at"`
	ResponseTimeMs *int       `json:"response_time_ms" db:"response_time_ms"`
}

// Constants for conversation and message types
const (
	// Source types
	SourceTypeComment = "comment"
	SourceTypeMessage = "message"

	// Conversation status
	StatusActive = "active"
	StatusPaused = "paused"
	StatusClosed = "closed"

	// Sender types
	SenderTypeUser = "user"
	SenderTypeBot  = "bot"

	// Message types
	MessageTypeText   = "text"
	MessageTypeImage  = "image"
	MessageTypeSticker = "sticker"
	MessageTypeVideo  = "video"
	MessageTypeAudio  = "audio"
	MessageTypeFile   = "file"
)
