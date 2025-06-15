package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"ruaychatbot/backend/internal/ChatBotCore/repository"
)

// RuleEngine handles keyword matching and response generation
type RuleEngine struct {
	keywordRepo      repository.KeywordRepository
	responseRepo     repository.ResponseTemplateRepository
	conversationRepo repository.ConversationRepository

	// Cache for better performance
	rulesCache      map[string][]*Rule
	lastCacheUpdate time.Time
	cacheTimeout    time.Duration
}

// Rule represents a complete rule with keywords, responses, and page configuration
type Rule struct {
	ID             int             `json:"id"`
	Name           string          `json:"name"`
	Keywords       []string        `json:"keywords"`
	Responses      []*Response     `json:"responses"`
	FallbackRules  []*FallbackRule `json:"fallback_rules,omitempty"`
	Enabled        bool            `json:"enabled"`
	SelectedPages  []string        `json:"selected_pages"`
	HideAfterReply bool            `json:"hide_after_reply"`
	SendToInbox    bool            `json:"send_to_inbox"`
	InboxMessage   string          `json:"inbox_message,omitempty"`
	InboxImage     string          `json:"inbox_image,omitempty"`
	Priority       int             `json:"priority"`
	MatchType      string          `json:"match_type"` // "exact", "contains", "starts_with", "ends_with"
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// Response represents a response template
type Response struct {
	Text     string `json:"text"`
	ImageURL string `json:"image_url,omitempty"`
	Type     string `json:"type"` // "text", "image", "template"
}

// FallbackRule represents a fallback response when no keywords match
type FallbackRule struct {
	ID            int         `json:"id"`
	Name          string      `json:"name"`
	Responses     []*Response `json:"responses"`
	Enabled       bool        `json:"enabled"`
	SelectedPages []string    `json:"selected_pages"`
	Priority      int         `json:"priority"`
}

// MatchResult represents the result of keyword matching
type MatchResult struct {
	Matched        bool          `json:"matched"`
	Rule           *Rule         `json:"rule,omitempty"`
	FallbackRule   *FallbackRule `json:"fallback_rule,omitempty"`
	MatchedKeyword string        `json:"matched_keyword,omitempty"`
	Response       *Response     `json:"response,omitempty"`
	MatchType      string        `json:"match_type"` // "keyword", "fallback"
	ProcessingTime time.Duration `json:"processing_time"`
}

// NewRuleEngine creates a new rule engine
func NewRuleEngine(repos *repository.Repositories) *RuleEngine {
	return &RuleEngine{
		keywordRepo:      repos.Keyword,
		responseRepo:     repos.ResponseTemplate,
		conversationRepo: repos.Conversation,
		rulesCache:       make(map[string][]*Rule),
		cacheTimeout:     5 * time.Minute, // Cache for 5 minutes
	}
}

// MatchRules finds matching rules for a given page and content
func (re *RuleEngine) MatchRules(ctx context.Context, pageID, content string) (*MatchResult, error) {
	startTime := time.Now()

	// Get rules for this page
	rules, err := re.getRulesForPage(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rules for page %s: %w", pageID, err)
	}

	content = strings.ToLower(strings.TrimSpace(content))

	// Try to match keywords first
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}

		for _, keyword := range rule.Keywords {
			if matched := re.matchKeyword(content, keyword, rule.MatchType); matched {
				response := re.selectRandomResponse(rule.Responses)

				return &MatchResult{
					Matched:        true,
					Rule:           rule,
					MatchedKeyword: keyword,
					Response:       response,
					MatchType:      "keyword",
					ProcessingTime: time.Since(startTime),
				}, nil
			}
		}
	}

	// Try fallback rules if no keyword matches
	fallbackRules, err := re.getFallbackRulesForPage(ctx, pageID)
	if err != nil {
		log.Printf("⚠️ Failed to get fallback rules: %v", err)
		return re.createNoMatchResult(startTime), nil
	}

	for _, fallbackRule := range fallbackRules {
		if !fallbackRule.Enabled {
			continue
		}

		response := re.selectRandomResponse(fallbackRule.Responses)
		return &MatchResult{
			Matched:        true,
			FallbackRule:   fallbackRule,
			Response:       response,
			MatchType:      "fallback",
			ProcessingTime: time.Since(startTime),
		}, nil
	}

	return re.createNoMatchResult(startTime), nil
}

// getRulesForPage gets all rules applicable to a specific page
func (re *RuleEngine) getRulesForPage(ctx context.Context, pageID string) ([]*Rule, error) {
	// Check cache first
	if rules, exists := re.rulesCache[pageID]; exists && time.Since(re.lastCacheUpdate) < re.cacheTimeout {
		return rules, nil
	}

	// Fetch from database
	keywords, err := re.keywordRepo.GetActive(ctx)
	if err != nil {
		return nil, err
	}

	// Convert keywords to rules (simplified for Phase 1)
	var rules []*Rule
	for _, keyword := range keywords {
		rule := &Rule{
			ID:             keyword.ID,
			Name:           fmt.Sprintf("Rule for %s", keyword.Keyword),
			Keywords:       []string{keyword.Keyword},
			Responses:      []*Response{{Text: keyword.Response, Type: "text"}},
			Enabled:        keyword.IsActive,
			SelectedPages:  []string{pageID}, // For now, apply to all pages
			HideAfterReply: false,
			SendToInbox:    false,
			Priority:       keyword.Priority,
			MatchType:      keyword.MatchType,
			CreatedAt:      keyword.CreatedAt,
			UpdatedAt:      keyword.UpdatedAt,
		}
		rules = append(rules, rule)
	}

	// Update cache
	re.rulesCache[pageID] = rules
	re.lastCacheUpdate = time.Now()

	return rules, nil
}

// getFallbackRulesForPage gets fallback rules for a specific page
func (re *RuleEngine) getFallbackRulesForPage(ctx context.Context, pageID string) ([]*FallbackRule, error) {
	// For Phase 1, return a default fallback rule
	return []*FallbackRule{
		{
			ID:   0,
			Name: "Default Fallback",
			Responses: []*Response{
				{Text: "ขอบคุณสำหรับข้อความครับ เราจะตอบกลับโดยเร็วที่สุด", Type: "text"},
			},
			Enabled:       true,
			SelectedPages: []string{pageID},
			Priority:      0,
		},
	}, nil
}

// matchKeyword checks if content matches a keyword based on match type
func (re *RuleEngine) matchKeyword(content, keyword, matchType string) bool {
	keyword = strings.ToLower(strings.TrimSpace(keyword))

	switch matchType {
	case "exact":
		return content == keyword
	case "contains":
		return strings.Contains(content, keyword)
	case "starts_with":
		return strings.HasPrefix(content, keyword)
	case "ends_with":
		return strings.HasSuffix(content, keyword)
	default:
		// Default to contains
		return strings.Contains(content, keyword)
	}
}

// selectRandomResponse selects a random response from available responses
func (re *RuleEngine) selectRandomResponse(responses []*Response) *Response {
	if len(responses) == 0 {
		return &Response{Text: "ขอบคุณสำหรับข้อความ", Type: "text"}
	}

	if len(responses) == 1 {
		return responses[0]
	}

	return responses[rand.Intn(len(responses))]
}

// createNoMatchResult creates a result for when no rules match
func (re *RuleEngine) createNoMatchResult(startTime time.Time) *MatchResult {
	return &MatchResult{
		Matched:        false,
		MatchType:      "none",
		ProcessingTime: time.Since(startTime),
	}
}

// ClearCache clears the rules cache
func (re *RuleEngine) ClearCache() {
	re.rulesCache = make(map[string][]*Rule)
	re.lastCacheUpdate = time.Time{}
	log.Println("🧹 Rule engine cache cleared")
}

// GetRuleStats returns statistics about rules
func (re *RuleEngine) GetRuleStats(ctx context.Context) (map[string]interface{}, error) {
	count, err := re.keywordRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_rules":       count,
		"cached_pages":      len(re.rulesCache),
		"last_cache_update": re.lastCacheUpdate,
		"cache_timeout":     re.cacheTimeout,
	}, nil
}

// ValidateRule validates a rule configuration
func (re *RuleEngine) ValidateRule(rule *Rule) error {
	if rule.Name == "" {
		return fmt.Errorf("rule name is required")
	}

	if len(rule.Keywords) == 0 {
		return fmt.Errorf("at least one keyword is required")
	}

	if len(rule.Responses) == 0 {
		return fmt.Errorf("at least one response is required")
	}

	for _, keyword := range rule.Keywords {
		if strings.TrimSpace(keyword) == "" {
			return fmt.Errorf("empty keyword not allowed")
		}
	}

	for _, response := range rule.Responses {
		if strings.TrimSpace(response.Text) == "" && response.ImageURL == "" {
			return fmt.Errorf("response must have text or image")
		}
	}

	return nil
}

// GetRulesForPage returns all rules configured for a specific page (public method)
func (re *RuleEngine) GetRulesForPage(ctx context.Context, pageID string) ([]*Rule, error) {
	return re.getRulesForPage(ctx, pageID)
}
