package service

import (
	"context"
	"fmt"
	"time"

	"ruaymanagement/backend/internal/ChatBotCore/models"
	"ruaymanagement/backend/internal/ChatBotCore/repository"
)

// MessageService provides business logic for message operations
type MessageService struct {
	repo repository.MessageRepository
}

// NewMessageService creates a new message service
func NewMessageService(repo repository.MessageRepository) *MessageService {
	return &MessageService{
		repo: repo,
	}
}

// CreateMessage creates a new message
func (s *MessageService) CreateMessage(ctx context.Context, conversationID int, messageID *string, senderType, messageType string, content *string, metadata string) (*models.Message, error) {
	msg := &models.Message{
		ConversationID: conversationID,
		MessageID:      messageID,
		SenderType:     senderType,
		MessageType:    messageType,
		Content:        content,
		Metadata:       metadata,
		CreatedAt:      time.Now(),
	}

	err := s.repo.Create(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	return msg, nil
}

// CreateUserMessage creates a message from a user
func (s *MessageService) CreateUserMessage(ctx context.Context, conversationID int, messageID *string, messageType string, content *string, metadata string) (*models.Message, error) {
	return s.CreateMessage(ctx, conversationID, messageID, models.SenderTypeUser, messageType, content, metadata)
}

// CreateBotMessage creates a message from the bot
func (s *MessageService) CreateBotMessage(ctx context.Context, conversationID int, messageType string, content *string, metadata string, responseTimeMs *int) (*models.Message, error) {
	msg := &models.Message{
		ConversationID: conversationID,
		SenderType:     models.SenderTypeBot,
		MessageType:    messageType,
		Content:        content,
		Metadata:       metadata,
		CreatedAt:      time.Now(),
		ProcessedAt:    &time.Time{},
		ResponseTimeMs: responseTimeMs,
	}
	*msg.ProcessedAt = time.Now()

	err := s.repo.Create(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot message: %w", err)
	}

	return msg, nil
}

// GetMessage retrieves a message by ID
func (s *MessageService) GetMessage(ctx context.Context, id int) (*models.Message, error) {
	return s.repo.GetByID(ctx, id)
}

// GetMessagesByConversation retrieves messages for a conversation with pagination
func (s *MessageService) GetMessagesByConversation(ctx context.Context, conversationID int, offset, limit int) ([]models.Message, error) {
	return s.repo.GetByConversationID(ctx, conversationID, offset, limit)
}

// GetMessagesBySenderType retrieves messages by sender type for a conversation with pagination
func (s *MessageService) GetMessagesBySenderType(ctx context.Context, conversationID int, senderType string, offset, limit int) ([]models.Message, error) {
	return s.repo.GetBySenderType(ctx, conversationID, senderType, offset, limit)
}

// CountMessagesByConversation returns the count of messages in a conversation
func (s *MessageService) CountMessagesByConversation(ctx context.Context, conversationID int) (int, error) {
	return s.repo.CountByConversationID(ctx, conversationID)
}

// DeleteMessage deletes a message
func (s *MessageService) DeleteMessage(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
