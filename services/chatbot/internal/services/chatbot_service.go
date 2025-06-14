package services

import (
	"context"
	"log"

	"chatbot/internal/models"
	"chatbot/internal/repository"
)

type ChatbotService struct {
	keywordRepo *repository.KeywordRepository
	messageRepo *repository.MessageRepository
}

func NewChatbotService(keywordRepo *repository.KeywordRepository, messageRepo *repository.MessageRepository) *ChatbotService {
	return &ChatbotService{
		keywordRepo: keywordRepo,
		messageRepo: messageRepo,
	}
}

func (s *ChatbotService) ProcessMessage(ctx context.Context, req *models.AutoReplyRequest) (*models.AutoReplyResponse, error) {
	log.Printf("🤖 Processing: '%s' from user: %s", req.Content, req.UserID)
	
	// Store the incoming message for future AI analysis
	message := &models.Message{
		Platform:    req.Platform,
		UserID:      req.UserID,
		PageID:      req.PageID,
		Content:     req.Content,
		CommentID:   req.CommentID,
		PostID:      req.PostID,
		MessageType: "comment", // assuming comment for now
		Processed:   false,
	}
	
	// Save message to database
	if err := s.messageRepo.Create(ctx, message); err != nil {
		log.Printf("Failed to save message: %v", err)
		// Continue processing even if save fails
	}
	
	keywords, err := s.keywordRepo.FindMatching(ctx, req.Content)
	if err != nil {
		return &models.AutoReplyResponse{ShouldReply: false}, err
	}
	
	if len(keywords) == 0 {
		return &models.AutoReplyResponse{ShouldReply: false}, nil
	}
	
	keyword := keywords[0] // Use highest priority
	response := keyword.Response
	
	// Mark message as processed
	if message.ID != 0 {
		s.messageRepo.MarkProcessed(ctx, message.ID, response)
	}
	
	return &models.AutoReplyResponse{
		ShouldReply:    true,
		Response:       response,
		MatchedKeyword: keyword.Keyword,
		MatchType:      keyword.MatchType,
	}, nil
}
