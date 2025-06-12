package models

import (
	"time"
)

// FacebookMessage represents an incoming Facebook message
type FacebookMessage struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_time"`
}

// FacebookUser represents a Facebook user
type FacebookUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// FacebookComment represents a Facebook comment
type FacebookComment struct {
	ID        string       `json:"id"`
	Message   string       `json:"message"`
	CreatedAt time.Time    `json:"created_time"`
	From      FacebookUser `json:"from"`
	PostID    string       `json:"post_id"`
}

// FacebookWebhookPayload represents the complete webhook payload from Facebook
type FacebookWebhookPayload struct {
	Object string             `json:"object"`
	Entry  []FacebookWebhookEntry `json:"entry"`
}

// FacebookWebhookEntry represents a single entry in the webhook payload
type FacebookWebhookEntry struct {
	ID      string                   `json:"id"`
	Time    int64                    `json:"time"`
	Changes []FacebookWebhookChange  `json:"changes"`
}

// FacebookWebhookChange represents a change notification
type FacebookWebhookChange struct {
	Field string                     `json:"field"`
	Value FacebookWebhookChangeValue `json:"value"`
}

// FacebookWebhookChangeValue represents the value of a change
type FacebookWebhookChangeValue struct {
	From      FacebookUser `json:"from"`
	Message   string       `json:"message"`
	CommentID string       `json:"comment_id"`
	PostID    string       `json:"post_id"`
	CreatedAt int64        `json:"created_time"`
}

// InternalMessage represents the standardized internal message format
// This is what gets sent to ChatBotCore
type InternalMessage struct {
	Content        string            `json:"content"`
	SenderID       string            `json:"sender_id"`
	ConversationID string            `json:"conversation_id"`
	Timestamp      time.Time         `json:"timestamp"`
	MessageType    string            `json:"message_type"` // "text", "image", etc.
	Metadata       map[string]string `json:"metadata,omitempty"` // Platform-specific metadata
}

// InternalResponse represents the response from ChatBotCore
type InternalResponse struct {
	Success        bool   `json:"success"`
	Reply          string `json:"reply,omitempty"`
	MatchedKeyword string `json:"matched_keyword,omitempty"`
	ShouldReply    bool   `json:"should_reply"`
	Error          string `json:"error,omitempty"`
}

// FacebookReplyRequest represents a request to send a reply via Facebook API
type FacebookReplyRequest struct {
	CommentID string `json:"comment_id"`
	Message   string `json:"message"`
	PageToken string `json:"page_token"`
}
