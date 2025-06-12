package services

import (
	"ruaymanagement/backend/external/FacebookConnect/models"
	"time"
)

// MessageTransformer handles transformation between Facebook and internal formats
type MessageTransformer struct{}

// NewMessageTransformer creates a new message transformer
func NewMessageTransformer() *MessageTransformer {
	return &MessageTransformer{}
}

// FacebookToInternal converts Facebook webhook data to internal message format
func (t *MessageTransformer) FacebookToInternal(change models.FacebookWebhookChange) *models.InternalMessage {
	// Create conversation ID from Facebook user ID and post ID
	conversationID := change.Value.From.ID
	if change.Value.PostID != "" {
		conversationID = change.Value.From.ID + "_" + change.Value.PostID
	}

	// Convert timestamp
	timestamp := time.Now()
	if change.Value.CreatedAt > 0 {
		timestamp = time.Unix(change.Value.CreatedAt, 0)
	}

	return &models.InternalMessage{
		Content:        change.Value.Message,
		SenderID:       change.Value.From.ID,
		ConversationID: conversationID,
		Timestamp:      timestamp,
		MessageType:    "text",
		Metadata: map[string]string{
			"platform":    "facebook",
			"sender_name": change.Value.From.Name,
			"comment_id":  change.Value.CommentID,
			"post_id":     change.Value.PostID,
		},
	}
}

// ValidateMessage checks if a Facebook message should be processed
func (t *MessageTransformer) ValidateMessage(change models.FacebookWebhookChange) bool {
	// Skip empty messages
	if change.Value.Message == "" {
		return false
	}

	// Skip messages without comment ID (not a comment)
	if change.Value.CommentID == "" {
		return false
	}

	// Skip messages with mentions (@ tags) - business rule
	if containsMention(change.Value.Message) {
		return false
	}

	return true
}

// containsMention checks if message contains @ mentions
func containsMention(message string) bool {
	// Simple check for @ symbol
	// In production, this could be more sophisticated
	for _, char := range message {
		if char == '@' {
			return true
		}
	}
	return false
}
