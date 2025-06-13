package service

import (
	"context"
	"fmt"
	"time"

	"ruaymanagement/backend/internal/ChatBotCore/models"
	"ruaymanagement/backend/internal/ChatBotCore/repository"
)

// ChatBotService provides the main business logic for the chatbot
type ChatBotService struct {
	conversationService *ConversationService
	messageService      *MessageService
	keywordService      *KeywordService
}

// NewChatBotService creates a new chatbot service
func NewChatBotService(
	conversationRepo repository.ConversationRepository,
	messageRepo repository.MessageRepository,
	keywordRepo repository.KeywordRepository,
) *ChatBotService {
	return &ChatBotService{
		conversationService: NewConversationService(conversationRepo),
		messageService:      NewMessageService(messageRepo),
		keywordService:      NewKeywordService(keywordRepo),
	}
}

// ProcessMessageRequest represents a message processing request
type ProcessMessageRequest struct {
	PageID           string  `json:"page_id"`
	FacebookUserID   string  `json:"facebook_user_id"`
	FacebookUserName *string `json:"facebook_user_name,omitempty"`
	MessageID        *string `json:"message_id,omitempty"`
	MessageText      string  `json:"message_text"`
	MessageType      string  `json:"message_type"`
	SourceType       string  `json:"source_type"`
	PostID           *string `json:"post_id,omitempty"`
	Metadata         string  `json:"metadata,omitempty"`
}

// ProcessMessageResponse represents the response from message processing
type ProcessMessageResponse struct {
	ConversationID   int     `json:"conversation_id"`
	UserMessageID    int     `json:"user_message_id"`
	BotMessageID     *int    `json:"bot_message_id,omitempty"`
	ResponseText     *string `json:"response_text,omitempty"`
	MatchedKeyword   *string `json:"matched_keyword,omitempty"`
	ProcessingTimeMs int     `json:"processing_time_ms"`
}

// ProcessMessage processes an incoming message and generates a response
func (s *ChatBotService) ProcessMessage(ctx context.Context, req ProcessMessageRequest) (*ProcessMessageResponse, error) {
	startTime := time.Now()

	// Validate input
	if req.FacebookUserID == "" {
		return nil, fmt.Errorf("facebook_user_id is required")
	}
	if req.MessageText == "" {
		return nil, fmt.Errorf("message_text is required")
	}

	// Get or create conversation
	conversation, err := s.conversationService.GetOrCreateConversation(
		ctx, req.PageID, req.FacebookUserID, req.SourceType, req.PostID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create conversation: %w", err)
	}

	// Update conversation user name if provided
	if req.FacebookUserName != nil && conversation.FacebookUserName != req.FacebookUserName {
		conversation.FacebookUserName = req.FacebookUserName
		if err := s.conversationService.UpdateConversation(ctx, conversation); err != nil {
			// Log but don't fail
			fmt.Printf("Warning: failed to update conversation user name: %v\n", err)
		}
	}

	// Create user message
	userMessage, err := s.messageService.CreateUserMessage(
		ctx, conversation.ID, req.MessageID, req.MessageType, &req.MessageText, req.Metadata,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user message: %w", err)
	}

	// Find matching keyword
	matchedKeyword, err := s.keywordService.FindMatchingKeyword(ctx, req.MessageText)
	if err != nil {
		return nil, fmt.Errorf("failed to find matching keyword: %w", err)
	}

	response := &ProcessMessageResponse{
		ConversationID:   conversation.ID,
		UserMessageID:    userMessage.ID,
		ProcessingTimeMs: int(time.Since(startTime).Milliseconds()),
	}

	// Generate response if keyword matched
	if matchedKeyword != nil {
		responseTimeMs := int(time.Since(startTime).Milliseconds())

		// Create bot response message
		botMessage, err := s.messageService.CreateBotMessage(
			ctx, conversation.ID, models.MessageTypeText, &matchedKeyword.Response, "", &responseTimeMs,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create bot message: %w", err)
		}

		response.BotMessageID = &botMessage.ID
		response.ResponseText = &matchedKeyword.Response
		response.MatchedKeyword = &matchedKeyword.Keyword
	}

	// Update conversation last message time
	if err := s.conversationService.UpdateLastMessageTime(ctx, conversation.ID); err != nil {
		// Log but don't fail
		fmt.Printf("Warning: failed to update last message time: %v\n", err)
	}

	return response, nil
}

// GetConversations retrieves conversations with pagination
func (s *ChatBotService) GetConversations(ctx context.Context, offset, limit int) ([]models.Conversation, int, error) {
	conversations, err := s.conversationService.ListConversations(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.conversationService.CountConversations(ctx)
	if err != nil {
		return nil, 0, err
	}

	return conversations, total, nil
}

// GetConversation retrieves a specific conversation
func (s *ChatBotService) GetConversation(ctx context.Context, id int) (*models.Conversation, error) {
	return s.conversationService.GetConversation(ctx, id)
}

// SetConversationStatus updates a conversation status
func (s *ChatBotService) SetConversationStatus(ctx context.Context, id int, status string) error {
	return s.conversationService.SetConversationStatus(ctx, id, status)
}

// GetMessages retrieves messages for a conversation
func (s *ChatBotService) GetMessages(ctx context.Context, conversationID int, offset, limit int) ([]models.Message, int, error) {
	messages, err := s.messageService.GetMessagesByConversation(ctx, conversationID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.messageService.CountMessagesByConversation(ctx, conversationID)
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// GetKeywords retrieves all keywords
func (s *ChatBotService) GetKeywords(ctx context.Context) ([]models.SimpleKeyword, error) {
	return s.keywordService.GetAllKeywords(ctx)
}

// CreateKeyword creates a new keyword
func (s *ChatBotService) CreateKeyword(ctx context.Context, keyword, response string) error {
	return s.keywordService.CreateKeyword(ctx, keyword, response)
}

// UpdateKeyword updates a keyword
func (s *ChatBotService) UpdateKeyword(ctx context.Context, id int, keyword, response string) error {
	return s.keywordService.UpdateKeyword(ctx, id, keyword, response)
}

// DeleteKeyword deletes a keyword
func (s *ChatBotService) DeleteKeyword(ctx context.Context, id int) error {
	return s.keywordService.DeleteKeyword(ctx, id)
}

// ProcessMessagePublic processes a message request from external API (models package)
func (s *ChatBotService) ProcessMessagePublic(ctx context.Context, req *models.ProcessMessageRequest) (*models.ProcessMessageResponse, error) {
	// Convert models.ProcessMessageRequest to service.ProcessMessageRequest
	serviceReq := ProcessMessageRequest{
		PageID:           "default", // Default page ID
		FacebookUserID:   req.FacebookUserID,
		FacebookUserName: &req.UserName,
		MessageID:        req.FacebookMessageID,
		MessageText:      req.MessageText,
		MessageType:      req.MessageType,
		SourceType:       "api", // Default source type
		PostID:           nil,   // No post ID from API request
		Metadata:         "",    // Convert map to string if needed
	}

	// Convert MessageData to metadata string if present
	if req.MessageData != nil && len(req.MessageData) > 0 {
		// For now, just use empty string - could serialize to JSON if needed
		serviceReq.Metadata = ""
	}

	// Call internal ProcessMessage
	serviceResp, err := s.ProcessMessage(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	// Convert service.ProcessMessageResponse to models.ProcessMessageResponse
	modelsResp := &models.ProcessMessageResponse{
		ShouldReply:      serviceResp.ResponseText != nil,
		ResponseType:     "text", // Default to text
		ResponseContent:  "",
		ResponseData:     make(map[string]interface{}),
		MatchedKeyword:   serviceResp.MatchedKeyword,
		ProcessingTimeMs: serviceResp.ProcessingTimeMs,
		ConversationID:   serviceResp.ConversationID,
		MessageID:        serviceResp.UserMessageID,
	}

	// Set response content if we have a response
	if serviceResp.ResponseText != nil {
		modelsResp.ResponseContent = *serviceResp.ResponseText
	}

	return modelsResp, nil
}
