package services

import (
	"context"
	"webhook/internal/models"
)

// PlatformHandler defines interface for handling different platforms
type PlatformHandler interface {
	// VerifyWebhook verifies incoming webhook signature/token
	VerifyWebhook(signature, body string) bool

	// ParseEvent parses platform-specific webhook event
	ParseEvent(rawData interface{}) (*models.UnifiedMessage, error)

	// SendResponse sends response back to the platform
	SendResponse(ctx context.Context, message *models.UnifiedMessage, response string) error

	// GetPlatform returns the platform type
	GetPlatform() models.Platform
}

// WebhookService manages all platform handlers
type WebhookService struct {
	handlers   map[models.Platform]PlatformHandler
	chatbotURL string
}

// NewWebhookService creates a new webhook service
func NewWebhookService(chatbotURL string) *WebhookService {
	return &WebhookService{
		handlers:   make(map[models.Platform]PlatformHandler),
		chatbotURL: chatbotURL,
	}
}

// RegisterHandler registers a platform handler
func (ws *WebhookService) RegisterHandler(handler PlatformHandler) {
	ws.handlers[handler.GetPlatform()] = handler
}

// GetHandler returns handler for specific platform
func (ws *WebhookService) GetHandler(platform models.Platform) (PlatformHandler, bool) {
	handler, exists := ws.handlers[platform]
	return handler, exists
}

// ProcessWebhook processes incoming webhook event
func (ws *WebhookService) ProcessWebhook(ctx context.Context, platform models.Platform, event *models.WebhookEvent) error {
	handler, exists := ws.GetHandler(platform)
	if !exists {
		return ErrPlatformNotSupported
	}

	// Parse platform-specific event to unified message
	message, err := handler.ParseEvent(event.RawData)
	if err != nil {
		return err
	}

	// Send to chatbot service
	response, err := ws.sendToChatbot(ctx, message)
	if err != nil {
		return err
	}

	// Send response back to platform if needed
	if response.ShouldReply {
		return handler.SendResponse(ctx, message, response.Response)
	}

	return nil
}

// sendToChatbot sends message to chatbot service
func (ws *WebhookService) sendToChatbot(ctx context.Context, message *models.UnifiedMessage) (*models.ChatbotResponse, error) {
	// TODO: Implement HTTP client to chatbot service
	// This will call chatbot service on port 8090
	return &models.ChatbotResponse{}, nil
}
