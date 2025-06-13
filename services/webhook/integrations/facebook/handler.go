package facebook

import (
	"encoding/json"
	"fmt"
	"log"

	"webhook/internal/model"
)

// Handler handles Facebook webhook events
type Handler struct {
	verifyToken string
	pageToken   string
	appSecret   string
	client      *Client
}

// NewHandler creates a new Facebook webhook handler
func NewHandler(verifyToken, pageToken, appSecret string) *Handler {
	return &Handler{
		verifyToken: verifyToken,
		pageToken:   pageToken,
		appSecret:   appSecret,
		client:      NewClient(pageToken, appSecret),
	}
}

// VerifyWebhook verifies the webhook challenge
func (h *Handler) VerifyWebhook(challengeToken string) error {
	return VerifyToken(challengeToken, h.verifyToken)
}

// HandleWebhook processes Facebook webhook data and converts to UnifiedMessage
func (h *Handler) HandleWebhook(payload []byte, signature string) ([]*model.UnifiedMessage, error) {
	// Verify webhook signature
	if err := VerifyWebhook(payload, signature, h.appSecret); err != nil {
		return nil, fmt.Errorf("webhook verification failed: %w", err)
	}
	
	// Parse Facebook webhook payload
	var fbPayload FacebookWebhookPayload
	if err := json.Unmarshal(payload, &fbPayload); err != nil {
		return nil, fmt.Errorf("failed to parse Facebook payload: %w", err)
	}
	
	var messages []*model.UnifiedMessage
	
	// Process each entry
	for _, entry := range fbPayload.Entry {
		// Process changes (comments, posts)
		for _, change := range entry.Changes {
			if shouldProcessChange(change) {
				unified := MapToUnifiedMessage(entry, change)
				messages = append(messages, unified)
			}
		}
		
		// Process direct messages
		for _, messaging := range entry.Messaging {
			if messaging.Message.Text != "" {
				unified := MapMessagingToUnifiedMessage(entry, messaging)
				messages = append(messages, unified)
			}
		}
	}
	
	return messages, nil
}

// SendResponse sends a response back to Facebook based on the chatbot response
func (h *Handler) SendResponse(response *model.ChatbotResponse, originalMsg *model.UnifiedMessage) error {
	if !response.ShouldReply {
		return nil
	}
	
	var err error
	
	// Send comment reply
	if originalMsg.CommentID != "" {
		_, err = h.client.SendComment(originalMsg.CommentID, response.Response)
		if err != nil {
			log.Printf("Failed to send Facebook comment reply: %v", err)
			return fmt.Errorf("failed to send comment reply: %w", err)
		}
		log.Printf("Sent Facebook comment reply to %s", originalMsg.CommentID)
	} else {
		// Send direct message
		_, err = h.client.SendMessage(originalMsg.UserID, response.Response)
		if err != nil {
			log.Printf("Failed to send Facebook message: %v", err)
			return fmt.Errorf("failed to send message: %w", err)
		}
		log.Printf("Sent Facebook message to user %s", originalMsg.UserID)
	}
	
	return nil
}

// shouldProcessChange determines if a Facebook change should be processed
func shouldProcessChange(change FacebookChange) bool {
	// Only process comment additions for now
	if change.Field == "comments" && change.Value.Verb == "add" {
		return true
	}
	
	// Can add more conditions for posts, reactions, etc.
	return false
}
