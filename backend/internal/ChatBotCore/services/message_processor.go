package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"ruaychatbot/backend/internal/ChatBotCore/channels"
	"ruaychatbot/backend/internal/ChatBotCore/repository"
)

// MessageProcessor handles incoming messages from any channel
type MessageProcessor struct {
	channels    map[string]channels.Channel
	ruleEngine  *RuleEngine
	pageManager *PageManager

	// Services
	conversationService *ConversationService
	messageService      *MessageService
	analytics           *AnalyticsService
}

// NewMessageProcessor creates a new message processor
func NewMessageProcessor(repos *repository.Repositories) *MessageProcessor {
	mp := &MessageProcessor{
		channels:            make(map[string]channels.Channel),
		ruleEngine:          NewRuleEngine(repos),
		pageManager:         NewPageManager(repos),
		conversationService: NewConversationService(repos.Conversation),
		messageService:      NewMessageService(repos.Message),
		analytics:           NewAnalyticsService(repos.Analytics),
	}

	return mp
}

// RegisterChannel registers a new channel
func (mp *MessageProcessor) RegisterChannel(channel channels.Channel, config channels.ChannelConfig) error {
	platform := channel.GetPlatform()
	mp.channels[platform] = channel

	// Store channel configuration in database
	return mp.pageManager.RegisterChannel(config)
}

// ProcessWebhook processes a webhook from any channel
func (mp *MessageProcessor) ProcessWebhook(ctx context.Context, platform string, payload []byte, signature string) error {
	// Get channel adapter
	channel, ok := mp.channels[platform]
	if !ok {
		return fmt.Errorf("unsupported platform: %s", platform)
	}

	// Validate webhook
	if err := channel.ValidateWebhook(payload, signature); err != nil {
		return fmt.Errorf("webhook validation failed: %w", err)
	}

	// Parse messages
	messages, err := channel.ParseWebhook(payload)
	if err != nil {
		return fmt.Errorf("failed to parse webhook: %w", err)
	}

	// Process each message
	for _, msg := range messages {
		if err := mp.ProcessMessage(ctx, msg); err != nil {
			// Log error but continue processing other messages
			fmt.Printf("Error processing message %s: %v\n", msg.MessageID, err)
			continue
		}
	}

	return nil
}

// ProcessMessage processes a single incoming message
func (mp *MessageProcessor) ProcessMessage(ctx context.Context, msg *channels.IncomingMessage) error {
	startTime := time.Now()

	// 1. Check if page is enabled
	page, err := mp.pageManager.GetPageByPlatformID(msg.Platform, msg.PageID)
	if err != nil {
		return fmt.Errorf("failed to get page: %w", err)
	}

	if !page.IsEnabled {
		return nil // Skip disabled pages
	}

	// 2. Get or create conversation
	conversation, err := mp.conversationService.GetOrCreateConversation(
		ctx, page.ID, msg.UserID, msg.MessageType, msg.PostID,
	)
	if err != nil {
		return fmt.Errorf("failed to get conversation: %w", err)
	}

	// 3. Store incoming message
	userMessage, err := mp.messageService.CreateUserMessage(
		ctx, conversation.ID, &msg.MessageID, msg.MessageType, &msg.Content, "",
	)
	if err != nil {
		return fmt.Errorf("failed to store message: %w", err)
	}

	// 4. Match rules
	rules, err := mp.ruleEngine.MatchRules(ctx, page.ID, msg.Content)
	if err != nil {
		return fmt.Errorf("failed to match rules: %w", err)
	}

	// 5. Process matched rules
	if len(rules) > 0 {
		// Select rule (prioritized or random)
		selectedRule := mp.selectRule(rules)

		// Generate response
		response, err := mp.generateResponse(selectedRule)
		if err != nil {
			return fmt.Errorf("failed to generate response: %w", err)
		}

		// Send reply
		channel := mp.channels[msg.Platform]
		reply := &channels.OutgoingReply{
			Platform:  msg.Platform,
			PageID:    msg.PageID,
			UserID:    msg.UserID,
			ReplyToID: msg.ContentID,
			Text:      response.Text,
			ImageURL:  response.ImageURL,
			Options: map[string]interface{}{
				"post_id": msg.PostID,
			},
		}

		result, err := channel.SendReply(ctx, reply)
		if err != nil {
			return fmt.Errorf("failed to send reply: %w", err)
		}

		// Store bot response
		processingTime := int(time.Since(startTime).Milliseconds())
		botMessage, err := mp.messageService.CreateBotMessage(
			ctx, conversation.ID, "text", &response.Text, "", &processingTime,
		)
		if err != nil {
			// Log but don't fail
			fmt.Printf("Warning: failed to store bot message: %v\n", err)
		}

		// Post-processing actions
		if selectedRule.HideAfterReply && channel.HideContent != nil {
			if err := channel.HideContent(ctx, msg.ContentID); err != nil {
				fmt.Printf("Warning: failed to hide content: %v\n", err)
			}
		}

		if selectedRule.SendToInbox && channel.SendToInbox != nil {
			inboxMsg := &channels.InboxMessage{
				Platform:     msg.Platform,
				PageID:       msg.PageID,
				UserID:       msg.UserID,
				Text:         selectedRule.InboxMessage,
				ImageURL:     selectedRule.InboxImageURL,
				SourceType:   "auto_reply",
				SourceRuleID: &selectedRule.ID,
			}

			if err := channel.SendToInbox(ctx, inboxMsg); err != nil {
				fmt.Printf("Warning: failed to send inbox message: %v\n", err)
			}
		}

		// Record analytics
		mp.analytics.RecordRuleExecution(ctx, selectedRule.ID, page.ID, result.Success)

		return nil
	}

	// No rules matched - check for fallback
	fallbackRules, err := mp.ruleEngine.GetFallbackRules(ctx, page.ID)
	if err != nil {
		return fmt.Errorf("failed to get fallback rules: %w", err)
	}

	if len(fallbackRules) > 0 {
		// Process fallback rule
		selectedRule := mp.selectRule(fallbackRules)
		// ... similar processing as above
	}

	return nil
}

// selectRule selects a rule from multiple matches (priority-based with randomization)
func (mp *MessageProcessor) selectRule(rules []*RuleMatch) *RuleMatch {
	if len(rules) == 1 {
		return rules[0]
	}

	// Group by priority
	maxPriority := rules[0].Priority
	highestPriorityRules := []*RuleMatch{}

	for _, rule := range rules {
		if rule.Priority > maxPriority {
			maxPriority = rule.Priority
			highestPriorityRules = []*RuleMatch{rule}
		} else if rule.Priority == maxPriority {
			highestPriorityRules = append(highestPriorityRules, rule)
		}
	}

	// Random selection among highest priority rules
	if len(highestPriorityRules) == 1 {
		return highestPriorityRules[0]
	}

	return highestPriorityRules[rand.Intn(len(highestPriorityRules))]
}

// generateResponse generates a response from a rule (with random selection)
func (mp *MessageProcessor) generateResponse(rule *RuleMatch) (*RuleResponse, error) {
	if len(rule.Responses) == 0 {
		return nil, fmt.Errorf("no responses available for rule %d", rule.ID)
	}

	// Random selection among responses
	response := rule.Responses[rand.Intn(len(rule.Responses))]
	return response, nil
}

// HealthCheck checks health of all registered channels
func (mp *MessageProcessor) HealthCheck(ctx context.Context) map[string]error {
	results := make(map[string]error)

	for platform, channel := range mp.channels {
		results[platform] = channel.HealthCheck(ctx)
	}

	return results
}

// Rule matching structures
type RuleMatch struct {
	ID              int             `json:"id"`
	Name            string          `json:"name"`
	Priority        int             `json:"priority"`
	MatchedKeywords []string        `json:"matched_keywords"`
	Responses       []*RuleResponse `json:"responses"`
	HideAfterReply  bool            `json:"hide_after_reply"`
	SendToInbox     bool            `json:"send_to_inbox"`
	InboxMessage    string          `json:"inbox_message"`
	InboxImageURL   *string         `json:"inbox_image_url"`
	ConfidenceScore float64         `json:"confidence_score"`
}

type RuleResponse struct {
	ID       int     `json:"id"`
	Text     string  `json:"text"`
	ImageURL *string `json:"image_url"`
	Priority int     `json:"priority"`
}

// Note: RuleEngine and PageManager are defined in their respective files
// This file imports and uses them

// AnalyticsService handles analytics recording
type AnalyticsService struct {
	repo repository.AnalyticsRepository
}

func NewAnalyticsService(repo repository.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

// RecordRuleExecution records rule execution for analytics
func (as *AnalyticsService) RecordRuleExecution(ctx context.Context, ruleID, pageID int, success bool) {
	// Record analytics data
}
