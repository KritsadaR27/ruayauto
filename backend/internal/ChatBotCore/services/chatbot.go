package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ruaymanagement/backend/internal/ChatBotCore/config"
	"ruaymanagement/backend/internal/ChatBotCore/models"
	"ruaymanagement/backend/internal/ChatBotCore/repository"
)

// ChatBotService handles core chatbot logic
type ChatBotService struct {
	repos  *repository.Repositories
	config *config.Config
}

// NewChatBotService creates a new chatbot service
func NewChatBotService(repos *repository.Repositories, config *config.Config) *ChatBotService {
	return &ChatBotService{
		repos:  repos,
		config: config,
	}
}

// ProcessMessage processes an incoming message and generates a response
func (s *ChatBotService) ProcessMessage(ctx context.Context, req *models.ProcessMessageRequest) (*models.ProcessMessageResponse, error) {
	startTime := time.Now()

	// 1. Find or create conversation
	conversation, err := s.findOrCreateConversation(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create conversation: %w", err)
	}

	// 2. Store incoming user message
	userMessage, err := s.storeUserMessage(ctx, conversation.ID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to store user message: %w", err)
	}

	// 3. Process message and find response (placeholder - will implement keyword matching later)
	response := s.processMessageLogic(ctx, req.MessageText)

	// 4. Store bot response message (if we're replying)
	if response.ShouldReply {
		_, err = s.storeBotMessage(ctx, conversation.ID, response)
		if err != nil {
			return nil, fmt.Errorf("failed to store bot message: %w", err)
		}
	}

	// 5. Update conversation last message time
	if err := s.repos.Conversation.UpdateLastMessageTime(ctx, conversation.ID); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to update last message time: %v\n", err)
	}

	// Calculate processing time
	processingTime := int(time.Since(startTime).Milliseconds())

	// Build response
	response.ProcessingTimeMs = processingTime
	response.ConversationID = conversation.ID
	response.MessageID = userMessage.ID

	return response, nil
}

// findOrCreateConversation finds existing conversation or creates a new one
func (s *ChatBotService) findOrCreateConversation(ctx context.Context, req *models.ProcessMessageRequest) (*models.Conversation, error) {
	// Try to find existing conversation
	conv, err := s.repos.Conversation.GetByFacebookUserID(ctx, req.FacebookUserID)
	if err == nil {
		// Update user name if it changed
		if conv.UserName != req.UserName {
			conv.UserName = req.UserName
			if updateErr := s.repos.Conversation.Update(ctx, conv); updateErr != nil {
				fmt.Printf("Warning: failed to update user name: %v\n", updateErr)
			}
		}
		return conv, nil
	}

	// Create new conversation
	newConv := &models.Conversation{
		FacebookUserID: req.FacebookUserID,
		UserName:       req.UserName,
		Status:         "active",
	}

	if err := s.repos.Conversation.Create(ctx, newConv); err != nil {
		return nil, err
	}

	return newConv, nil
}

// storeUserMessage stores the incoming user message
func (s *ChatBotService) storeUserMessage(ctx context.Context, conversationID int, req *models.ProcessMessageRequest) (*models.Message, error) {
	message := &models.Message{
		ConversationID:    conversationID,
		FacebookMessageID: req.FacebookMessageID,
		SenderType:        "user",
		MessageType:       req.MessageType,
		MessageText:       &req.MessageText,
	}

	if err := s.repos.Message.Create(ctx, message); err != nil {
		return nil, err
	}

	return message, nil
}

// storeBotMessage stores the bot's response message
func (s *ChatBotService) storeBotMessage(ctx context.Context, conversationID int, response *models.ProcessMessageResponse) (*models.Message, error) {
	message := &models.Message{
		ConversationID: conversationID,
		SenderType:     "bot",
		MessageType:    response.ResponseType,
		MessageText:    &response.ResponseContent,
		MatchedKeyword: response.MatchedKeyword,
	}

	if err := s.repos.Message.Create(ctx, message); err != nil {
		return nil, err
	}

	return message, nil
}

// processMessageLogic handles the core message processing logic
func (s *ChatBotService) processMessageLogic(ctx context.Context, messageText string) *models.ProcessMessageResponse {
	response := &models.ProcessMessageResponse{
		ShouldReply:     false,
		ResponseType:    "text",
		ResponseContent: s.getDefaultResponse(messageText),
	}

	// Find matching keywords
	matchedKeywords, err := s.repos.Keyword.FindMatchingKeywords(ctx, messageText)
	if err != nil {
		fmt.Printf("Warning: failed to find matching keywords: %v\n", err)
		return response
	}

	// If we have matched keywords, use the highest priority one
	if len(matchedKeywords) > 0 {
		bestMatch := matchedKeywords[0] // Already sorted by priority in repository
		response.ShouldReply = true
		response.MatchedKeyword = &bestMatch.Keyword
		response.ResponseContent = bestMatch.ResponseContent
		response.ResponseType = bestMatch.ResponseType
	} else {
		// Fallback to simple keyword matching for common cases
		if contains(messageText, []string{"สวัสดี", "หวัดดี", "hello", "hi"}) {
			response.ShouldReply = true
			response.ResponseContent = "สวัสดีค่ะ! ยินดีต้อนรับครับ 😊"
			response.MatchedKeyword = stringPtr("สวัสดี")
		} else if contains(messageText, []string{"ขอบคุณ", "thank you", "thanks"}) {
			response.ShouldReply = true
			response.ResponseContent = "ยินดีค่ะ! มีอะไรให้ช่วยอีกไหมคะ"
			response.MatchedKeyword = stringPtr("ขอบคุณ")
		} else if contains(messageText, []string{"ราคา", "เท่าไหร่", "price"}) {
			response.ShouldReply = true
			response.ResponseContent = "สำหรับข้อมูลราคา กรุณาแจ้งสินค้าที่สนใจค่ะ"
			response.MatchedKeyword = stringPtr("ราคา")
		}
	}

	return response
}

// getDefaultResponse returns a default response
func (s *ChatBotService) getDefaultResponse(messageText string) string {
	return "ขอบคุณที่ติดต่อมาค่ะ! เรากำลังดำเนินการตอบกลับ"
}

// GetConversations retrieves conversations with pagination
func (s *ChatBotService) GetConversations(ctx context.Context, offset, limit int) ([]models.Conversation, int, error) {
	conversations, err := s.repos.Conversation.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repos.Conversation.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return conversations, total, nil
}

// GetConversation retrieves a specific conversation
func (s *ChatBotService) GetConversation(ctx context.Context, id int) (*models.Conversation, error) {
	return s.repos.Conversation.GetByID(ctx, id)
}

// UpdateConversationStatus updates conversation status
func (s *ChatBotService) UpdateConversationStatus(ctx context.Context, id int, status string) error {
	return s.repos.Conversation.SetStatus(ctx, id, status)
}

// Helper functions
func contains(text string, keywords []string) bool {
	textLower := strings.ToLower(text)
	for _, keyword := range keywords {
		if strings.Contains(textLower, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func stringPtr(s string) *string {
	return &s
}
