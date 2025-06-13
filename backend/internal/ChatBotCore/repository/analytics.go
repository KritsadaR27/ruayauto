package repository

import (
	"context"
	"time"

	"ruaymanagement/backend/internal/ChatBotCore/database"
	"ruaymanagement/backend/internal/ChatBotCore/models"
)

type analyticsRepository struct {
	db *database.DB
}

func NewAnalyticsRepository(db *database.DB) AnalyticsRepository {
	return &analyticsRepository{db: db}
}

// GetConversationStats retrieves conversation statistics for a date range
func (r *analyticsRepository) GetConversationStats(ctx context.Context, startDate, endDate string) ([]models.ConversationStats, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT 
			DATE(c.created_at) as date,
			COUNT(DISTINCT c.id) as total_conversations,
			COUNT(m.id) as total_messages,
			COUNT(CASE WHEN m.sender_type = 'bot' THEN 1 END) as bot_responses,
			COALESCE(AVG(EXTRACT(EPOCH FROM (m.created_at - LAG(m.created_at) OVER (ORDER BY m.created_at))) * 1000), 0) as avg_response_time_ms
		FROM chatbot_mvp.conversations c
		LEFT JOIN chatbot_mvp.messages m ON c.id = m.conversation_id
		WHERE DATE(c.created_at) BETWEEN $1 AND $2
		GROUP BY DATE(c.created_at)
		ORDER BY DATE(c.created_at)
	`, startDate, endDate)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []models.ConversationStats	
	for rows.Next() {
		var stat models.ConversationStats
		err := rows.Scan(
			&stat.Date,
			&stat.TotalConversations,
			&stat.TotalMessages,
			&stat.BotResponses,
			&stat.AvgResponseTimeMs,
		)
		if err != nil {
			return nil, err
		}
		stat.CreatedAt = time.Now()
		stats = append(stats, stat)
	}
	
	return stats, rows.Err()
}

// CreateDailyStat creates a new daily statistic record
func (r *analyticsRepository) CreateDailyStat(ctx context.Context, stat *models.ConversationStats) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO chatbot_mvp.conversation_stats (date, total_conversations, total_messages, bot_responses, avg_response_time_ms, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, stat.Date, stat.TotalConversations, stat.TotalMessages, stat.BotResponses, stat.AvgResponseTimeMs, stat.CreatedAt)
	
	return err
}

// UpdateDailyStat updates an existing daily statistic record
func (r *analyticsRepository) UpdateDailyStat(ctx context.Context, date string, stat *models.ConversationStats) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE chatbot_mvp.conversation_stats 
		SET total_conversations = $2, total_messages = $3, bot_responses = $4, avg_response_time_ms = $5
		WHERE date = $1
	`, date, stat.TotalConversations, stat.TotalMessages, stat.BotResponses, stat.AvgResponseTimeMs)
	
	return err
}

// GetTotalConversations returns the total number of conversations
func (r *analyticsRepository) GetTotalConversations(ctx context.Context) (int, error) {
	var count int
	err := r.db.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM chatbot_mvp.conversations
	`).Scan(&count)
	
	return count, err
}

// GetTotalMessages returns the total number of messages
func (r *analyticsRepository) GetTotalMessages(ctx context.Context) (int, error) {
	var count int
	err := r.db.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM chatbot_mvp.messages
	`).Scan(&count)
	
	return count, err
}

// GetActiveConversations returns the number of active conversations
func (r *analyticsRepository) GetActiveConversations(ctx context.Context) (int, error) {
	var count int
	err := r.db.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM chatbot_mvp.conversations WHERE status = 'active'
	`).Scan(&count)
	
	return count, err
}

// GetAverageResponseTime calculates average response time for the last N days
func (r *analyticsRepository) GetAverageResponseTime(ctx context.Context, days int) (float64, error) {
	var avgTime float64
	err := r.db.Pool.QueryRow(ctx, `
		SELECT COALESCE(AVG(
			EXTRACT(EPOCH FROM (bot_msg.created_at - user_msg.created_at))
		), 0) as avg_response_time_seconds
		FROM chatbot_mvp.messages user_msg
		JOIN chatbot_mvp.messages bot_msg ON bot_msg.conversation_id = user_msg.conversation_id
		WHERE user_msg.created_at >= NOW() - INTERVAL '%d days'
			AND user_msg.sender_type = 'user'
			AND bot_msg.sender_type = 'bot'
			AND bot_msg.created_at > user_msg.created_at
			AND bot_msg.id = (
				SELECT MIN(id) FROM chatbot_mvp.messages
				WHERE conversation_id = user_msg.conversation_id
					AND sender_type = 'bot'
					AND created_at > user_msg.created_at
			)
	`, days).Scan(&avgTime)
	
	return avgTime, err
}

// GetTopKeywords returns the most frequently used keywords
func (r *analyticsRepository) GetTopKeywords(ctx context.Context, limit int, days int) ([]map[string]interface{}, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT 
			k.keyword_text,
			COUNT(m.id) as usage_count
		FROM chatbot_mvp.keywords k
		LEFT JOIN chatbot_mvp.messages m ON m.matched_keyword = k.keyword_text
			AND m.created_at >= NOW() - INTERVAL '%d days'
		WHERE k.is_active = true
		GROUP BY k.keyword_text
		ORDER BY usage_count DESC
		LIMIT $1
	`, days, limit)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var keyword string
		var count int
		
		err := rows.Scan(&keyword, &count)
		if err != nil {
			return nil, err
		}
		
		results = append(results, map[string]interface{}{
			"keyword": keyword,
			"count":   count,
		})
	}
	
	return results, rows.Err()
}

// GetResponseTypeDistribution returns distribution of response types
func (r *analyticsRepository) GetResponseTypeDistribution(ctx context.Context, days int) (map[string]int, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT 
			response_type,
			COUNT(*) as count
		FROM chatbot_mvp.messages
		WHERE sender_type = 'bot' 
			AND created_at >= NOW() - INTERVAL '%d days'
		GROUP BY response_type
	`, days)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make(map[string]int)
	for rows.Next() {
		var responseType string
		var count int
		
		err := rows.Scan(&responseType, &count)
		if err != nil {
			return nil, err
		}
		
		results[responseType] = count
	}
	
	return results, rows.Err()
}

// GetHourlyMessageDistribution returns hourly message distribution
func (r *analyticsRepository) GetHourlyMessageDistribution(ctx context.Context, days int) ([]map[string]interface{}, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT 
			EXTRACT(HOUR FROM created_at) as hour,
			COUNT(*) as message_count
		FROM chatbot_mvp.messages
		WHERE created_at >= NOW() - INTERVAL '%d days'
		GROUP BY EXTRACT(HOUR FROM created_at)
		ORDER BY hour
	`, days)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var hour int
		var count int
		
		err := rows.Scan(&hour, &count)
		if err != nil {
			return nil, err
		}
		
		results = append(results, map[string]interface{}{
			"hour":  hour,
			"count": count,
		})
	}
	
	return results, rows.Err()
}
