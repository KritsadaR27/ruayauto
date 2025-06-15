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
	ruleRepo    *repository.RuleRepository
	messageRepo *repository.MessageRepository
}

func NewChatbotService(ruleRepo *repository.RuleRepository, messageRepo *repository.MessageRepository) *ChatbotService {
	return &ChatbotService{
		ruleRepo:    ruleRepo,
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
	if err := s.ruleRepo.LogIncomingMessage(ctx, incomingMsg); err != nil {
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

	// Find matching rules
	rules, err := s.ruleRepo.FindMatching(ctx, req.Content)
	if err != nil {
		return &models.AutoReplyResponse{ShouldReply: false}, err
	}

	// If no rules matched, try fallback responses
	if len(rules) == 0 {
		return s.handleFallbackResponse(ctx, pageID, req)
	}

	// Use highest priority rule
	rule := rules[0]

	// Check rate limiting
	rateLimited, err := s.ruleRepo.CheckRateLimit(ctx, pageID, req.UserID, rule.ID, 5) // 5 minutes cooldown
	if err != nil {
		log.Printf("Rate limit check failed: %v", err)
	}

	if rateLimited {
		log.Printf("Rate limited for user %s, keywords %v", req.UserID, rule.Keywords)
		// Update analytics for trigger but no response
		s.ruleRepo.UpdateRuleAnalytics(ctx, rule.ID, pageID, true, false)
		return &models.AutoReplyResponse{ShouldReply: false}, nil
	}

	// Get all responses for this rule
	responses, err := s.ruleRepo.GetRuleResponses(ctx, rule.ID)
	if err != nil {
		log.Printf("Failed to get rule responses: %v", err)
		// Fall back to legacy comma-separated responses
		return s.handleLegacyResponse(ctx, rule, message, pageID)
	}

	var selectedResponse *models.RuleResponse
	var responseText string

	// If we have multiple responses, randomly select one
	if len(responses) > 0 {
		selectedResponse = s.selectRandomResponse(responses)
		responseText = selectedResponse.ResponseText
	} else {
		// Fall back to legacy comma-separated responses in the Response field
		responseText = s.selectLegacyResponse(rule.Response)
	}

	// Update rate limit
	rateLimit := &models.RateLimit{
		PageID:        pageID,
		UserID:        req.UserID,
		KeywordID:     &rule.ID,
		ResponseCount: 1,
	}

	if err := s.ruleRepo.UpdateRateLimit(ctx, rateLimit); err != nil {
		log.Printf("Failed to update rate limit: %v", err)
	}

	// Update analytics
	s.ruleRepo.UpdateRuleAnalytics(ctx, rule.ID, pageID, true, true)

	// Update incoming message with matched rule
	incomingMsg.MatchedKeywordID = &rule.ID
	incomingMsg.ResponseSent = true

	// Update the incoming message record in the database
	if err := s.ruleRepo.UpdateIncomingMessage(ctx, incomingMsg); err != nil {
		log.Printf("Failed to update incoming message: %v", err)
	}

	// Mark legacy message as processed
	if message.ID != 0 {
		s.messageRepo.MarkProcessed(ctx, message.ID, responseText)
	}

	// Create response with composite support
	response := &models.AutoReplyResponse{
		ShouldReply:    true,
		Response:       responseText,
		MatchedKeyword: strings.Join(rule.Keywords, ", "),
		MatchType:      rule.MatchType,
	}

	// Add media information if available
	if selectedResponse != nil {
		response.HasMedia = selectedResponse.HasMedia
		response.MediaDescription = selectedResponse.MediaDescription
	}

	return response, nil
}

// Helper methods for enhanced functionality

// handleFallbackResponse handles cases where no keywords match
func (s *ChatbotService) handleFallbackResponse(ctx context.Context, pageID int, req *models.AutoReplyRequest) (*models.AutoReplyResponse, error) {
	// Try to get fallback responses for this page
	fallbackResponses, err := s.ruleRepo.GetFallbackResponses(ctx, &pageID)
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
func (s *ChatbotService) handleLegacyResponse(ctx context.Context, rule *models.Rule, message *models.Message, pageID int) (*models.AutoReplyResponse, error) {
	response := s.selectLegacyResponse(rule.Response)

	// Update analytics
	s.ruleRepo.UpdateRuleAnalytics(ctx, rule.ID, pageID, true, true)

	// Mark message as processed
	if message.ID != 0 {
		s.messageRepo.MarkProcessed(ctx, message.ID, response)
	}

	return &models.AutoReplyResponse{
		ShouldReply:    true,
		Response:       response,
		MatchedKeyword: strings.Join(rule.Keywords, ", "),
		MatchType:      rule.MatchType,
	}, nil
}

// selectRandomResponse selects a random response from keyword responses with weight consideration
func (s *ChatbotService) selectRandomResponse(responses []*models.RuleResponse) *models.RuleResponse {
	if len(responses) == 0 {
		return nil
	}

	// Calculate total weight
	totalWeight := 0
	for _, resp := range responses {
		totalWeight += resp.Weight
	}

	// If no weights, select randomly
	if totalWeight == 0 {
		rand.Seed(time.Now().UnixNano())
		return responses[rand.Intn(len(responses))]
	}

	// Select based on weight
	rand.Seed(time.Now().UnixNano())
	randomWeight := rand.Intn(totalWeight)
	currentWeight := 0

	for _, resp := range responses {
		currentWeight += resp.Weight
		if randomWeight < currentWeight {
			return resp
		}
	}

	// Fallback to first response
	return responses[0]
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
