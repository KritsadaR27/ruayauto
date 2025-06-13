package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"ruaymanagement/backend/internal/ChatBotCore/repository"
)

// AutoReplyService provides business logic for automatic replies
type AutoReplyService struct {
	keywordRepo repository.KeywordRepository
	rng         *rand.Rand
}

// AutoReplyConfig holds configuration for auto-reply behavior
type AutoReplyConfig struct {
	EnableDefault      bool     `json:"enableDefault"`
	DefaultResponses   []string `json:"defaultResponses"`
	CaseSensitive      bool     `json:"caseSensitive"`
	EnableFuzzyMatch   bool     `json:"enableFuzzyMatch"`
	MinMatchThreshold  float64  `json:"minMatchThreshold"`
}

// MatchResult represents the result of keyword matching
type MatchResult struct {
	Matched          bool    `json:"matched"`
	MatchedKeyword   string  `json:"matchedKeyword,omitempty"`
	Response         string  `json:"response,omitempty"`
	ConfidenceScore  float64 `json:"confidenceScore,omitempty"`
	MatchType        string  `json:"matchType,omitempty"` // "exact", "partial", "fuzzy", "default"
}

// NewAutoReplyService creates a new auto-reply service
func NewAutoReplyService(keywordRepo repository.KeywordRepository) *AutoReplyService {
	return &AutoReplyService{
		keywordRepo: keywordRepo,
		rng:         rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// MatchKeyword finds a matching keyword for the given text and returns a response
func (s *AutoReplyService) MatchKeyword(ctx context.Context, text string, config *AutoReplyConfig) (*MatchResult, error) {
	// Get all active keywords from repository
	keywords, err := s.keywordRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get keywords: %w", err)
	}

	// Prepare text for matching
	searchText := text
	if !config.CaseSensitive {
		searchText = strings.ToLower(text)
	}

	// Try exact keyword matching first
	for _, keyword := range keywords {
		matchKeyword := keyword.Keyword
		if !config.CaseSensitive {
			matchKeyword = strings.ToLower(keyword.Keyword)
		}

		if strings.Contains(searchText, matchKeyword) {
			return &MatchResult{
				Matched:         true,
				MatchedKeyword:  keyword.Keyword,
				Response:        keyword.Response,
				ConfidenceScore: 1.0,
				MatchType:       "exact",
			}, nil
		}
	}

	// Return default response if enabled
	if config.EnableDefault && len(config.DefaultResponses) > 0 {
		response := s.getRandomResponse(config.DefaultResponses)
		return &MatchResult{
			Matched:         true,
			Response:        response,
			ConfidenceScore: 0.1,
			MatchType:       "default",
		}, nil
	}

	// No match found
	return &MatchResult{
		Matched: false,
	}, nil
}

// getRandomResponse selects a random response from a list
func (s *AutoReplyService) getRandomResponse(responses []string) string {
	if len(responses) == 0 {
		return ""
	}
	if len(responses) == 1 {
		return responses[0]
	}
	return responses[s.rng.Intn(len(responses))]
}

// PostCommentReply posts a reply to a Facebook comment
func (s *AutoReplyService) PostCommentReply(ctx context.Context, commentID, reply, pageToken string) error {
	url := fmt.Sprintf("https://graph.facebook.com/v19.0/%s/comments", commentID)
	
	body := map[string]string{
		"message":      reply,
		"access_token": pageToken,
	}
	
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to post comment: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Facebook Graph API error: status %d", resp.StatusCode)
	}

	return nil
}

// GetAutoReplyStats returns statistics about auto-reply performance
func (s *AutoReplyService) GetAutoReplyStats(ctx context.Context, fromTime, toTime time.Time) (map[string]interface{}, error) {
	keywords, err := s.keywordRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get keywords for stats: %w", err)
	}

	return map[string]interface{}{
		"total_keywords":  len(keywords),
		"active_keywords": len(keywords),
		"period_start":    fromTime,
		"period_end":      toTime,
		"generated_at":    time.Now(),
	}, nil
}

// ValidateAutoReplyConfig validates the configuration
func (s *AutoReplyService) ValidateAutoReplyConfig(config *AutoReplyConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	if config.EnableDefault && len(config.DefaultResponses) == 0 {
		return fmt.Errorf("defaultResponses cannot be empty when enableDefault is true")
	}

	return nil
}

// GetDefaultConfig returns a default configuration
func (s *AutoReplyService) GetDefaultConfig() *AutoReplyConfig {
	return &AutoReplyConfig{
		EnableDefault:     true,
		DefaultResponses:  []string{"ขอบคุณที่แสดงความคิดเห็นค่ะ"},
		CaseSensitive:     false,
		EnableFuzzyMatch:  false,
		MinMatchThreshold: 0.6,
	}
}

// ProcessAutoReply processes a message and returns an auto-reply if appropriate
func (s *AutoReplyService) ProcessAutoReply(ctx context.Context, messageText string, config *AutoReplyConfig) (*MatchResult, error) {
	if err := s.ValidateAutoReplyConfig(config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	if strings.TrimSpace(messageText) == "" {
		return &MatchResult{Matched: false}, nil
	}

	return s.MatchKeyword(ctx, messageText, config)
}
