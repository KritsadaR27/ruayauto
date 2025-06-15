package services

import (
	"context"
	"log"

	"chatbot-core/models"
	"chatbot-core/repository"
)

// ChatBotService handles the main chatbot logic
type ChatBotService struct {
	keywordRepo      *repository.KeywordRepository
	conversationRepo *repository.ConversationRepository
}

// NewChatBotService creates a new chatbot service
func NewChatBotService(keywordRepo *repository.KeywordRepository, conversationRepo *repository.ConversationRepository) *ChatBotService {
	return &ChatBotService{
		keywordRepo:      keywordRepo,
		conversationRepo: conversationRepo,
	}
}

// ProcessMessage processes an incoming message and determines the response
func (s *ChatBotService) ProcessMessage(ctx context.Context, req *models.AutoReplyRequest) (*models.AutoReplyResponse, error) {
	log.Printf("Processing message: %s from user: %s on page: %s", req.Content, req.UserID, req.PageID)

	// Find or create conversation
	conv, err := s.conversationRepo.FindOrCreate(ctx, req.PageID, req.UserID, req.CommentID, req.PostID)
	if err != nil {
		log.Printf("Error finding/creating conversation: %v", err)
		return &models.AutoReplyResponse{ShouldReply: false}, err
	}

	// Add incoming message to conversation
	incomingMessage := &models.Message{
		ConversationID: conv.ID,
		Content:        req.Content,
		Direction:      "inbound",
		MessageType:    "text",
		ExternalID:     req.CommentID,
	}

	err = s.conversationRepo.AddMessage(ctx, incomingMessage)
	if err != nil {
		log.Printf("Error adding incoming message: %v", err)
	}

	// Find matching keywords
	matchingKeywords, err := s.keywordRepo.FindMatching(ctx, req.Content)
	if err != nil {
		log.Printf("Error finding matching keywords: %v", err)
		return &models.AutoReplyResponse{ShouldReply: false}, err
	}

	if len(matchingKeywords) == 0 {
		log.Println("No matching keywords found")
		return &models.AutoReplyResponse{ShouldReply: false}, nil
	}

	// Use the highest priority keyword (first in the sorted list)
	selectedKeyword := matchingKeywords[0]

	log.Printf("Found matching keyword: %s with response: %s", selectedKeyword.Keyword, selectedKeyword.Response)

	// Add outgoing message to conversation
	outgoingMessage := &models.Message{
		ConversationID: conv.ID,
		Content:        selectedKeyword.Response,
		Direction:      "outbound",
		MessageType:    "text",
	}

	err = s.conversationRepo.AddMessage(ctx, outgoingMessage)
	if err != nil {
		log.Printf("Error adding outgoing message: %v", err)
	}

	return &models.AutoReplyResponse{
		ShouldReply:    true,
		Response:       selectedKeyword.Response,
		MatchedKeyword: selectedKeyword.Keyword,
		MatchType:      selectedKeyword.MatchType,
	}, nil
}
