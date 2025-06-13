package line

import (
	"encoding/json"
	"fmt"
	"log"

	"webhook/internal/model"
)

// Handler handles LINE webhook events
type Handler struct {
	channelSecret string
	channelToken  string
	client        *Client
}

// NewHandler creates a new LINE webhook handler
func NewHandler(channelSecret, channelToken string) *Handler {
	return &Handler{
		channelSecret: channelSecret,
		channelToken:  channelToken,
		client:        NewClient(channelToken),
	}
}

// HandleWebhook processes LINE webhook data and converts to UnifiedMessage
func (h *Handler) HandleWebhook(payload []byte, signature string) ([]*model.UnifiedMessage, error) {
	// Verify webhook signature
	if err := VerifyWebhook(payload, signature, h.channelSecret); err != nil {
		return nil, fmt.Errorf("webhook verification failed: %w", err)
	}
	
	// Parse LINE webhook payload
	var linePayload LineWebhookPayload
	if err := json.Unmarshal(payload, &linePayload); err != nil {
		return nil, fmt.Errorf("failed to parse LINE payload: %w", err)
	}
	
	var messages []*model.UnifiedMessage
	
	// Process each event
	for _, event := range linePayload.Events {
		if shouldProcessEvent(event) {
			unified := MapToUnifiedMessage(event)
			messages = append(messages, unified)
		}
	}
	
	return messages, nil
}

// SendResponse sends a response back to LINE based on the chatbot response
func (h *Handler) SendResponse(response *model.ChatbotResponse, originalMsg *model.UnifiedMessage) error {
	if !response.ShouldReply {
		return nil
	}
	
	var err error
	
	// Send reply using reply token if available
	if originalMsg.ReplyToken != "" {
		_, err = h.client.ReplyMessage(originalMsg.ReplyToken, response.Response)
		if err != nil {
			log.Printf("Failed to send LINE reply: %v", err)
			return fmt.Errorf("failed to send reply: %w", err)
		}
		log.Printf("Sent LINE reply to user %s", originalMsg.UserID)
	} else {
		// Fall back to push message
		_, err = h.client.PushMessage(originalMsg.UserID, response.Response)
		if err != nil {
			log.Printf("Failed to send LINE push message: %v", err)
			return fmt.Errorf("failed to send push message: %w", err)
		}
		log.Printf("Sent LINE push message to user %s", originalMsg.UserID)
	}
	
	return nil
}

// shouldProcessEvent determines if a LINE event should be processed
func shouldProcessEvent(event LineEvent) bool {
	// Only process message events with text or sticker content
	if event.Type == "message" && event.Message != nil {
		return event.Message.Type == "text" || event.Message.Type == "sticker"
	}
	return false
}
