package services

import (
	"context"
	"log"

	"chatbot-core/models"
	"chatbot-core/repository"
)

// ChatBotService handles the main chatbot logic
type ChatBotService struct {
	keywordRepo *repository.KeywordRepository
}

// NewChatBotService creates a new chatbot service
func NewChatBotService(keywordRepo *repository.KeywordRepository) *ChatBotService {
	return &ChatBotService{
		keywordRepo: keywordRepo,
	}
}

// ProcessMessage processes an incoming message and determines the response
func (s *ChatBotService) ProcessMessage(ctx context.Context, req *models.AutoReplyRequest) (*models.AutoReplyResponse, error) {
	log.Printf("🤖 Processing message: '%s' from user: %s on page: %s", req.Content, req.UserID, req.PageID)

	// Find matching keywords
	matchingKeywords, err := s.keywordRepo.FindMatching(ctx, req.Content)
	if err != nil {
		log.Printf("❌ Error finding matching keywords: %v", err)
		return &models.AutoReplyResponse{ShouldReply: false}, err
	}

	if len(matchingKeywords) == 0 {
		log.Println("💭 No matching keywords found")
		return &models.AutoReplyResponse{ShouldReply: false}, nil
	}

	// Use the highest priority keyword (first in the sorted list)
	selectedKeyword := matchingKeywords[0]

	log.Printf("✅ Found matching keyword: '%s' -> '%s'", selectedKeyword.Keyword, selectedKeyword.Response)

	return &models.AutoReplyResponse{
		ShouldReply:    true,
		Response:       selectedKeyword.Response,
		MatchedKeyword: selectedKeyword.Keyword,
		MatchType:      selectedKeyword.MatchType,
	}, nil
}
