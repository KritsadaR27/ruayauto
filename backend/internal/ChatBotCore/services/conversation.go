package service

import (
	"context"
	"fmt"
	"time"

	"ruaymanagement/backend/internal/ChatBotCore/models"
	"ruaymanagement/backend/internal/ChatBotCore/repository"
)

// ConversationService provides business logic for conversation operations
type ConversationService struct {
	repo repository.ConversationRepository
}

// NewConversationService creates a new conversation service
func NewConversationService(repo repository.ConversationRepository) *ConversationService {
	return &ConversationService{
		repo: repo,
	}
}

// CreateConversation creates a new conversation
func (s *ConversationService) CreateConversation(ctx context.Context, pageID, facebookUserID string, sourceType string, postID *string) (*models.Conversation, error) {
	conv := &models.Conversation{
		PageID:           pageID,
		FacebookUserID:   facebookUserID,
		SourceType:       sourceType,
		PostID:           postID,
		Status:           models.StatusActive,
		StartedAt:        time.Now(),
		LastMessageAt:    time.Now(),
	}

	err := s.repo.Create(ctx, conv)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	return conv, nil
}

// GetOrCreateConversation gets an existing conversation or creates a new one
func (s *ConversationService) GetOrCreateConversation(ctx context.Context, pageID, facebookUserID string, sourceType string, postID *string) (*models.Conversation, error) {
	// Try to get existing conversation
	conv, err := s.repo.GetByFacebookUserID(ctx, facebookUserID)
	if err == nil {
		return conv, nil
	}

	// If not found, create new conversation
	return s.CreateConversation(ctx, pageID, facebookUserID, sourceType, postID)
}

// GetConversation retrieves a conversation by ID
func (s *ConversationService) GetConversation(ctx context.Context, id int) (*models.Conversation, error) {
	return s.repo.GetByID(ctx, id)
}

// UpdateConversation updates an existing conversation
func (s *ConversationService) UpdateConversation(ctx context.Context, conv *models.Conversation) error {
	return s.repo.Update(ctx, conv)
}

// ListConversations retrieves conversations with pagination
func (s *ConversationService) ListConversations(ctx context.Context, offset, limit int) ([]models.Conversation, error) {
	return s.repo.List(ctx, offset, limit)
}

// ListConversationsByStatus retrieves conversations by status with pagination
func (s *ConversationService) ListConversationsByStatus(ctx context.Context, status string, offset, limit int) ([]models.Conversation, error) {
	return s.repo.ListByStatus(ctx, status, offset, limit)
}

// CountConversations returns total number of conversations
func (s *ConversationService) CountConversations(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

// CountConversationsByStatus returns number of conversations by status
func (s *ConversationService) CountConversationsByStatus(ctx context.Context, status string) (int, error) {
	return s.repo.CountByStatus(ctx, status)
}

// UpdateLastMessageTime updates the last message time for a conversation
func (s *ConversationService) UpdateLastMessageTime(ctx context.Context, id int) error {
	return s.repo.UpdateLastMessageTime(ctx, id)
}

// SetConversationStatus updates the status of a conversation
func (s *ConversationService) SetConversationStatus(ctx context.Context, id int, status string) error {
	return s.repo.SetStatus(ctx, id, status)
}

// DeleteConversation deletes a conversation
func (s *ConversationService) DeleteConversation(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
