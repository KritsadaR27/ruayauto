package models

import (
	"encoding/json"
	"time"
)

// ResponseTemplate represents a reusable response template
type ResponseTemplate struct {
	ID           int             `json:"id" db:"id"`
	Name         string          `json:"name" db:"name"`
	TemplateType string          `json:"template_type" db:"template_type"` // text, image, quick_reply, generic
	TemplateData json.RawMessage `json:"template_data" db:"template_data"`
	Description  *string         `json:"description" db:"description"`
	IsActive     bool            `json:"is_active" db:"is_active"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    *time.Time      `json:"updated_at" db:"updated_at"`
}

// ConversationStats represents analytics data
type ConversationStats struct {
	ID                    int     `json:"id" db:"id"`
	Date                  string  `json:"date" db:"date"`
	TotalConversations    int     `json:"total_conversations" db:"total_conversations"`
	TotalMessages         int     `json:"total_messages" db:"total_messages"`
	BotResponses          int     `json:"bot_responses" db:"bot_responses"`
	AvgResponseTimeMs     float64 `json:"avg_response_time_ms" db:"avg_response_time_ms"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
}

// ProcessMessageRequest represents incoming message processing request
type ProcessMessageRequest struct {
	FacebookUserID    string                 `json:"facebook_user_id"`
	FacebookMessageID *string                `json:"facebook_message_id,omitempty"`
	UserName          string                 `json:"user_name"`
	MessageText       string                 `json:"message_text"`
	MessageType       string                 `json:"message_type"` // text, image, etc.
	MessageData       map[string]interface{} `json:"message_data,omitempty"`
}

// ProcessMessageResponse represents message processing response
type ProcessMessageResponse struct {
	ShouldReply      bool                   `json:"should_reply"`
	ResponseType     string                 `json:"response_type"` // text, image, quick_reply, template
	ResponseContent  string                 `json:"response_content"`
	ResponseData     map[string]interface{} `json:"response_data,omitempty"`
	MatchedKeyword   *string                `json:"matched_keyword,omitempty"`
	ProcessingTimeMs int                    `json:"processing_time_ms"`
	ConversationID   int                    `json:"conversation_id"`
	MessageID        int                    `json:"message_id"`
}

// SendMessageRequest represents outgoing message request to FacebookConnect
type SendMessageRequest struct {
	FacebookUserID  string                 `json:"facebook_user_id"`
	MessageType     string                 `json:"message_type"`
	MessageContent  string                 `json:"message_content"`
	MessageData     map[string]interface{} `json:"message_data,omitempty"`
}

// SendMessageResponse represents response from FacebookConnect
type SendMessageResponse struct {
	Success           bool    `json:"success"`
	FacebookMessageID *string `json:"facebook_message_id,omitempty"`
	Error             *string `json:"error,omitempty"`
}

// APIResponse represents standardized API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// APIError represents API error details
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Meta represents response metadata
type Meta struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

// Analytics Models

// KeywordStat represents keyword usage statistics
type KeywordStat struct {
	KeywordID           int    `json:"keyword_id"`
	KeywordText         string `json:"keyword_text"`
	UsageCount          int    `json:"usage_count"`
	UniqueConversations int    `json:"unique_conversations"`
}

// DailyMessageStat represents daily message statistics
type DailyMessageStat struct {
	MessageDate         time.Time `json:"message_date"`
	TotalMessages       int       `json:"total_messages"`
	UserMessages        int       `json:"user_messages"`
	BotMessages         int       `json:"bot_messages"`
	UniqueConversations int       `json:"unique_conversations"`
}

// ResponseTemplateStat represents response template usage statistics
type ResponseTemplateStat struct {
	TemplateID          int    `json:"template_id"`
	TemplateName        string `json:"template_name"`
	UsageCount          int    `json:"usage_count"`
	UniqueConversations int    `json:"unique_conversations"`
}

// OverallStats represents overall system statistics
type OverallStats struct {
	TotalMessages              int      `json:"total_messages"`
	UserMessages               int      `json:"user_messages"`
	BotMessages                int      `json:"bot_messages"`
	TotalConversations         int      `json:"total_conversations"`
	ActiveKeywords             int      `json:"active_keywords"`
	ActiveTemplates            int      `json:"active_templates"`
	AverageResponseTimeSeconds *float64 `json:"average_response_time_seconds,omitempty"`
}
