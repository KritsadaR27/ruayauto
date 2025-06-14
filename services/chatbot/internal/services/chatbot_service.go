package services

import (
	"context"
	"log"
	"math/rand"
	"strings"
	"time"

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
	
	// Parse page ID for rate limiting (convert string to int)
	pageID := 1 // Default page ID - should be parsed from req.PageID
	
	// Store the incoming message for analytics
	incomingMsg := &models.IncomingMessage{
		PageID:      pageID,
		UserID:      req.UserID,
		MessageText: req.Content,
		CommentID:   req.CommentID,
		PostID:      req.PostID,
	}
	
	// Log incoming message for analytics
	if err := s.keywordRepo.LogIncomingMessage(ctx, incomingMsg); err != nil {
		log.Printf("Failed to log incoming message: %v", err)
		// Continue processing even if logging fails
	}
	
	// Store the message in the legacy messages table too
	message := &models.Message{
		Platform:    req.Platform,
		UserID:      req.UserID,
		PageID:      req.PageID,
		Content:     req.Content,
		CommentID:   req.CommentID,
		PostID:      req.PostID,
		MessageType: "comment",
		Processed:   false,
	}
	
	if err := s.messageRepo.Create(ctx, message); err != nil {
		log.Printf("Failed to save message: %v", err)
	}
	
	// Find matching keywords
	keywords, err := s.keywordRepo.FindMatching(ctx, req.Content)
	if err != nil {
		return &models.AutoReplyResponse{ShouldReply: false}, err
	}
	
	// If no keywords matched, try fallback responses
	if len(keywords) == 0 {
		return s.handleFallbackResponse(ctx, pageID, req)
	}
	
	// Use highest priority keyword
	keyword := keywords[0]
	
	// Check rate limiting
	rateLimited, err := s.keywordRepo.CheckRateLimit(ctx, pageID, req.UserID, keyword.ID, 5) // 5 minutes cooldown
	if err != nil {
		log.Printf("Rate limit check failed: %v", err)
	}
	
	if rateLimited {
		log.Printf("Rate limited for user %s, keyword %s", req.UserID, keyword.Keyword)
		// Update analytics for trigger but no response
		s.keywordRepo.UpdateRuleAnalytics(ctx, keyword.ID, pageID, true, false)
		return &models.AutoReplyResponse{ShouldReply: false}, nil
	}
	
	// Get all responses for this keyword
	responses, err := s.keywordRepo.GetKeywordResponses(ctx, keyword.ID)
	if err != nil {
		log.Printf("Failed to get keyword responses: %v", err)
		// Fall back to legacy comma-separated responses
		return s.handleLegacyResponse(ctx, keyword, message, pageID)
	}
	
	var selectedResponse string
	
	// If we have multiple responses, randomly select one
	if len(responses) > 0 {
		selectedResponse = s.selectRandomResponse(responses)
	} else {
		// Fall back to legacy comma-separated responses in the Response field
		selectedResponse = s.selectLegacyResponse(keyword.Response)
	}
	
	// Update rate limit
	rateLimit := &models.RateLimit{
		PageID:        pageID,
		UserID:        req.UserID,
		KeywordID:     &keyword.ID,
		ResponseCount: 1,
	}
	
	if err := s.keywordRepo.UpdateRateLimit(ctx, rateLimit); err != nil {
		log.Printf("Failed to update rate limit: %v", err)
	}
	
	// Update analytics
	s.keywordRepo.UpdateRuleAnalytics(ctx, keyword.ID, pageID, true, true)
	
	// Update incoming message with matched keyword
	incomingMsg.MatchedKeywordID = &keyword.ID
	incomingMsg.ResponseSent = true
	
	// Update the incoming message record in the database
	if err := s.keywordRepo.UpdateIncomingMessage(ctx, incomingMsg); err != nil {
		log.Printf("Failed to update incoming message: %v", err)
	}
	
	// Mark legacy message as processed
	if message.ID != 0 {
		s.messageRepo.MarkProcessed(ctx, message.ID, selectedResponse)
	}
	
	return &models.AutoReplyResponse{
		ShouldReply:    true,
		Response:       selectedResponse,
		MatchedKeyword: keyword.Keyword,
		MatchType:      keyword.MatchType,
	}, nil
}

// Helper methods for enhanced functionality

// handleFallbackResponse handles cases where no keywords match
func (s *ChatbotService) handleFallbackResponse(ctx context.Context, pageID int, req *models.AutoReplyRequest) (*models.AutoReplyResponse, error) {
	// Try to get fallback responses for this page
	fallbackResponses, err := s.keywordRepo.GetFallbackResponses(ctx, &pageID)
	if err != nil {
		log.Printf("Failed to get fallback responses: %v", err)
		return &models.AutoReplyResponse{ShouldReply: false}, nil
	}
	
	if len(fallbackResponses) == 0 {
		return &models.AutoReplyResponse{ShouldReply: false}, nil
	}
	
	// Select random fallback response based on weight
	selectedResponse := s.selectRandomFallbackResponse(fallbackResponses)
	
	return &models.AutoReplyResponse{
		ShouldReply: true,
		Response:    selectedResponse,
	}, nil
}

// handleLegacyResponse handles backward compatibility with comma-separated responses
func (s *ChatbotService) handleLegacyResponse(ctx context.Context, keyword *models.Keyword, message *models.Message, pageID int) (*models.AutoReplyResponse, error) {
	response := s.selectLegacyResponse(keyword.Response)
	
	// Update analytics
	s.keywordRepo.UpdateRuleAnalytics(ctx, keyword.ID, pageID, true, true)
	
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

// selectRandomResponse selects a random response from keyword responses with weight consideration
func (s *ChatbotService) selectRandomResponse(responses []*models.KeywordResponse) string {
	if len(responses) == 0 {
		return ""
	}
	
	// Calculate total weight
	totalWeight := 0
	for _, resp := range responses {
		totalWeight += resp.Weight
	}
	
	// If no weights, select randomly
	if totalWeight == 0 {
		rand.Seed(time.Now().UnixNano())
		return responses[rand.Intn(len(responses))].ResponseText
	}
	
	// Select based on weight
	rand.Seed(time.Now().UnixNano())
	randomWeight := rand.Intn(totalWeight)
	currentWeight := 0
	
	for _, resp := range responses {
		currentWeight += resp.Weight
		if randomWeight < currentWeight {
			return resp.ResponseText
		}
	}
	
	// Fallback to first response
	return responses[0].ResponseText
}

// selectRandomFallbackResponse selects a random fallback response with weight consideration
func (s *ChatbotService) selectRandomFallbackResponse(responses []*models.FallbackResponse) string {
	if len(responses) == 0 {
		return ""
	}
	
	// Calculate total weight
	totalWeight := 0
	for _, resp := range responses {
		totalWeight += resp.Weight
	}
	
	// If no weights, select randomly
	if totalWeight == 0 {
		rand.Seed(time.Now().UnixNano())
		return responses[rand.Intn(len(responses))].ResponseText
	}
	
	// Select based on weight
	rand.Seed(time.Now().UnixNano())
	randomWeight := rand.Intn(totalWeight)
	currentWeight := 0
	
	for _, resp := range responses {
		currentWeight += resp.Weight
		if randomWeight < currentWeight {
			return resp.ResponseText
		}
	}
	
	// Fallback to first response
	return responses[0].ResponseText
}

// selectLegacyResponse handles legacy comma-separated responses
func (s *ChatbotService) selectLegacyResponse(response string) string {
	// Handle multiple responses separated by ", "
	if strings.Contains(response, ", ") {
		responses := strings.Split(response, ", ")
		// Remove empty responses
		var validResponses []string
		for _, r := range responses {
			if trimmed := strings.TrimSpace(r); trimmed != "" {
				validResponses = append(validResponses, trimmed)
			}
		}
		
		// Randomly select one response
		if len(validResponses) > 0 {
			rand.Seed(time.Now().UnixNano())
			return validResponses[rand.Intn(len(validResponses))]
		}
	}
	
	return response
}
