package services

import (
	"context"
	"log"

	"chatbot/internal/models"
	"chatbot/internal/repository"
)

type ChatbotService struct {
	keywordRepo *repository.KeywordRepository
}

func NewChatbotService(keywordRepo *repository.KeywordRepository) *ChatbotService {
	return &ChatbotService{keywordRepo: keywordRepo}
}

func (s *ChatbotService) ProcessMessage(ctx context.Context, req *models.AutoReplyRequest) (*models.AutoReplyResponse, error) {
	log.Printf("🤖 Processing: '%s' from user: %s", req.Content, req.UserID)
	
	keywords, err := s.keywordRepo.FindMatching(ctx, req.Content)
	if err != nil {
		return &models.AutoReplyResponse{ShouldReply: false}, err
	}
	
	if len(keywords) == 0 {
		return &models.AutoReplyResponse{ShouldReply: false}, nil
	}
	
	keyword := keywords[0] // Use highest priority
	
	return &models.AutoReplyResponse{
		ShouldReply:    true,
		Response:       keyword.Response,
		MatchedKeyword: keyword.Keyword,
		MatchType:      keyword.MatchType,
	}, nil
}
