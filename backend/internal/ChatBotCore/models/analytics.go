package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ruaymanagement/backend/internal/ChatBotCore/database"
)

// MessageAnalytics represents a message interaction record
type MessageAnalytics struct {
	ID                  int       `json:"id" db:"id"`
	FacebookMessageID   *string   `json:"facebook_message_id" db:"facebook_message_id"`
	SenderID            string    `json:"sender_id" db:"sender_id"`
	RecipientID         string    `json:"recipient_id" db:"recipient_id"`
	MessageText         *string   `json:"message_text" db:"message_text"`
	MatchedKeyword      *string   `json:"matched_keyword" db:"matched_keyword"`
	ResponseSent        *string   `json:"response_sent" db:"response_sent"`
	ResponseTimeMs      *int      `json:"response_time_ms" db:"response_time_ms"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}

// WebhookLog represents a webhook request/response log
type WebhookLog struct {
	ID               int             `json:"id" db:"id"`
	WebhookType      string          `json:"webhook_type" db:"webhook_type"`
	RequestBody      json.RawMessage `json:"request_body" db:"request_body"`
	ResponseStatus   *int            `json:"response_status" db:"response_status"`
	ResponseBody     *string         `json:"response_body" db:"response_body"`
	ProcessingTimeMs *int            `json:"processing_time_ms" db:"processing_time_ms"`
	ErrorMessage     *string         `json:"error_message" db:"error_message"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
}

// AnalyticsRepository handles database operations for analytics
type AnalyticsRepository struct {
	db *database.DB
}

// NewAnalyticsRepository creates a new analytics repository
func NewAnalyticsRepository(db *database.DB) *AnalyticsRepository {
	return &AnalyticsRepository{
		db: db,
	}
}

// LogMessage logs a message interaction
func (r *AnalyticsRepository) LogMessage(ctx context.Context, analytics *MessageAnalytics) error {
	query := `
		INSERT INTO message_analytics 
		(facebook_message_id, sender_id, recipient_id, message_text, matched_keyword, response_sent, response_time_ms) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id, created_at
	`

	err := r.db.Pool.QueryRow(ctx, query,
		analytics.FacebookMessageID,
		analytics.SenderID,
		analytics.RecipientID,
		analytics.MessageText,
		analytics.MatchedKeyword,
		analytics.ResponseSent,
		analytics.ResponseTimeMs,
	).Scan(&analytics.ID, &analytics.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to log message analytics: %w", err)
	}

	return nil
}

// LogWebhook logs a webhook request/response
func (r *AnalyticsRepository) LogWebhook(ctx context.Context, log *WebhookLog) error {
	query := `
		INSERT INTO webhook_logs 
		(webhook_type, request_body, response_status, response_body, processing_time_ms, error_message) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id, created_at
	`

	err := r.db.Pool.QueryRow(ctx, query,
		log.WebhookType,
		log.RequestBody,
		log.ResponseStatus,
		log.ResponseBody,
		log.ProcessingTimeMs,
		log.ErrorMessage,
	).Scan(&log.ID, &log.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to log webhook: %w", err)
	}

	return nil
}

// GetMessageStats returns message statistics for a time period
func (r *AnalyticsRepository) GetMessageStats(ctx context.Context, from, to time.Time) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_messages,
			COUNT(DISTINCT sender_id) as unique_senders,
			COUNT(CASE WHEN matched_keyword IS NOT NULL THEN 1 END) as matched_messages,
			COUNT(CASE WHEN response_sent IS NOT NULL THEN 1 END) as responses_sent,
			AVG(response_time_ms) as avg_response_time_ms
		FROM message_analytics 
		WHERE created_at BETWEEN $1 AND $2
	`

	var stats struct {
		TotalMessages       int     `db:"total_messages"`
		UniqueSenders       int     `db:"unique_senders"`
		MatchedMessages     int     `db:"matched_messages"`
		ResponsesSent       int     `db:"responses_sent"`
		AvgResponseTimeMs   *float64 `db:"avg_response_time_ms"`
	}

	err := r.db.Pool.QueryRow(ctx, query, from, to).Scan(
		&stats.TotalMessages,
		&stats.UniqueSenders,
		&stats.MatchedMessages,
		&stats.ResponsesSent,
		&stats.AvgResponseTimeMs,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get message stats: %w", err)
	}

	result := map[string]interface{}{
		"total_messages":        stats.TotalMessages,
		"unique_senders":        stats.UniqueSenders,
		"matched_messages":      stats.MatchedMessages,
		"responses_sent":        stats.ResponsesSent,
		"avg_response_time_ms":  stats.AvgResponseTimeMs,
	}

	return result, nil
}

// GetTopKeywords returns the most used keywords
func (r *AnalyticsRepository) GetTopKeywords(ctx context.Context, limit int, from, to time.Time) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			matched_keyword,
			COUNT(*) as usage_count
		FROM message_analytics 
		WHERE matched_keyword IS NOT NULL 
		AND created_at BETWEEN $1 AND $2
		GROUP BY matched_keyword 
		ORDER BY usage_count DESC 
		LIMIT $3
	`

	rows, err := r.db.Pool.Query(ctx, query, from, to, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top keywords: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var keyword string
		var count int
		err := rows.Scan(&keyword, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan keyword stats: %w", err)
		}
		results = append(results, map[string]interface{}{
			"keyword": keyword,
			"count":   count,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating keyword stats: %w", err)
	}

	return results, nil
}

// GetRecentWebhookLogs returns recent webhook logs
func (r *AnalyticsRepository) GetRecentWebhookLogs(ctx context.Context, limit int) ([]WebhookLog, error) {
	query := `
		SELECT id, webhook_type, request_body, response_status, response_body, 
			   processing_time_ms, error_message, created_at
		FROM webhook_logs 
		ORDER BY created_at DESC 
		LIMIT $1
	`

	rows, err := r.db.Pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get webhook logs: %w", err)
	}
	defer rows.Close()

	var logs []WebhookLog
	for rows.Next() {
		var log WebhookLog
		err := rows.Scan(
			&log.ID,
			&log.WebhookType,
			&log.RequestBody,
			&log.ResponseStatus,
			&log.ResponseBody,
			&log.ProcessingTimeMs,
			&log.ErrorMessage,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan webhook log: %w", err)
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating webhook logs: %w", err)
	}

	return logs, nil
}
