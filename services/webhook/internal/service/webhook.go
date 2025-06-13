package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"webhook/integrations/facebook"
	"webhook/integrations/line"
	"webhook/internal/model"
)

// WebhookService handles webhook business logic
type WebhookService struct {
	chatbotURL    string
	facebook      *facebook.Handler
	line          *line.Handler
	httpClient    *http.Client
}

// NewWebhookService creates a new webhook service
func NewWebhookService(
	chatbotURL string,
	facebook *facebook.Handler,
	line *line.Handler,
) *WebhookService {
	return &WebhookService{
		chatbotURL: chatbotURL,
		facebook:   facebook,
		line:       line,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// ProcessMessage processes a unified message
func (s *WebhookService) ProcessMessage(msg *model.UnifiedMessage) error {
	log.Printf("Processing message from %s platform, user: %s, content: %s", 
		msg.Platform, msg.UserID, msg.Content)

	// Send to chatbot service
	response, err := s.sendToChatbot(msg)
	if err != nil {
		return fmt.Errorf("failed to get chatbot response: %w", err)
	}

	// Send response back to the platform if needed
	if response.ShouldReply {
		if err := s.sendResponse(msg, response); err != nil {
			return fmt.Errorf("failed to send response: %w", err)
		}
	}

	return nil
}

// sendToChatbot sends the message to the chatbot service
func (s *WebhookService) sendToChatbot(msg *model.UnifiedMessage) (*model.ChatbotResponse, error) {
	// Create chatbot request
	request := &model.ChatbotRequest{
		Platform:  msg.Platform,
		UserID:    msg.UserID,
		PageID:    msg.PageID,
		Content:   msg.Content,
		CommentID: msg.CommentID,
		PostID:    msg.PostID,
	}

	// Marshal request
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Send HTTP request to chatbot service
	resp, err := s.httpClient.Post(
		s.chatbotURL+"/api/process",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("chatbot service returned status: %d", resp.StatusCode)
	}

	// Parse response
	var chatbotResponse model.ChatbotResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatbotResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &chatbotResponse, nil
}

// sendResponse sends the response back to the appropriate platform
func (s *WebhookService) sendResponse(originalMsg *model.UnifiedMessage, response *model.ChatbotResponse) error {
	switch originalMsg.Platform {
	case model.PlatformFacebook:
		return s.facebook.SendResponse(response, originalMsg)
	case model.PlatformLine:
		return s.line.SendResponse(response, originalMsg)
	case model.PlatformInstagram:
		// TODO: Implement Instagram response
		log.Printf("Instagram response not implemented yet")
		return nil
	case model.PlatformTikTok:
		// TODO: Implement TikTok response
		log.Printf("TikTok response not implemented yet")
		return nil
	case model.PlatformTwitter:
		// TODO: Implement Twitter response
		log.Printf("Twitter response not implemented yet")
		return nil
	default:
		return fmt.Errorf("unsupported platform: %s", originalMsg.Platform)
	}
}
