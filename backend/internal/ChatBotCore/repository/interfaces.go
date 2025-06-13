package repository

import (
	"context"
	"ruaymanagement/backend/internal/ChatBotCore/models"
)

// ConversationRepository handles conversation data operations
type ConversationRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, conv *models.Conversation) error
	GetByID(ctx context.Context, id int) (*models.Conversation, error)
	GetByFacebookUserID(ctx context.Context, facebookUserID string) (*models.Conversation, error)
	Update(ctx context.Context, conv *models.Conversation) error
	Delete(ctx context.Context, id int) error

	// List and search operations
	List(ctx context.Context, offset, limit int) ([]models.Conversation, error)
	ListByStatus(ctx context.Context, status string, offset, limit int) ([]models.Conversation, error)
	Count(ctx context.Context) (int, error)
	CountByStatus(ctx context.Context, status string) (int, error)

	// Business operations
	UpdateLastMessageTime(ctx context.Context, id int) error
	SetStatus(ctx context.Context, id int, status string) error
}

// MessageRepository handles message data operations
type MessageRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, msg *models.Message) error
	GetByID(ctx context.Context, id int) (*models.Message, error)
	Delete(ctx context.Context, id int) error

	// Conversation-related operations
	GetByConversationID(ctx context.Context, conversationID int, offset, limit int) ([]models.Message, error)
	CountByConversationID(ctx context.Context, conversationID int) (int, error)

	// Sender type operations
	GetBySenderType(ctx context.Context, conversationID int, senderType string, offset, limit int) ([]models.Message, error)
}

// KeywordRepository handles keyword data operations
type KeywordRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, keyword models.SimpleKeyword) error
	GetByID(ctx context.Context, id int) (*models.SimpleKeyword, error)
	Update(ctx context.Context, id int, keyword models.SimpleKeyword) error
	Delete(ctx context.Context, id int) error

	// List and search operations
	GetAll(ctx context.Context) ([]models.SimpleKeyword, error)
	GetActive(ctx context.Context) ([]models.SimpleKeyword, error)
	List(ctx context.Context, offset, limit int) ([]models.SimpleKeyword, error)
	Count(ctx context.Context) (int, error)

	// Specific lookup operations
	GetByKeyword(ctx context.Context, keyword string) (*models.SimpleKeyword, error)
	UpsertKeyword(ctx context.Context, keyword, response string) error
	SearchKeywords(ctx context.Context, searchTerm string) ([]models.SimpleKeyword, error)

	// Matching operations
	FindMatchingKeywords(ctx context.Context, text string) ([]models.SimpleKeyword, error)
	
	// Status operations
	SetActive(ctx context.Context, id int, isActive bool) error
}

// ResponseTemplateRepository handles response template operations
type ResponseTemplateRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, template *models.ResponseTemplate) error
	GetByID(ctx context.Context, id int) (*models.ResponseTemplate, error)
	Update(ctx context.Context, template *models.ResponseTemplate) error
	Delete(ctx context.Context, id int) error

	// List and search operations
	GetAll(ctx context.Context) ([]models.ResponseTemplate, error)
	GetActive(ctx context.Context) ([]models.ResponseTemplate, error)
	GetByType(ctx context.Context, templateType string) ([]models.ResponseTemplate, error)
	GetByName(ctx context.Context, name string) (*models.ResponseTemplate, error)

	// Business operations
	ToggleActive(ctx context.Context, id int, isActive bool) error
}

// AnalyticsRepository handles analytics and statistics
type AnalyticsRepository interface {
	// Conversation statistics
	GetConversationStats(ctx context.Context, startDate, endDate string) ([]models.ConversationStats, error)
	CreateDailyStat(ctx context.Context, stat *models.ConversationStats) error
	UpdateDailyStat(ctx context.Context, date string, stat *models.ConversationStats) error

	// Real-time analytics
	GetTotalConversations(ctx context.Context) (int, error)
	GetTotalMessages(ctx context.Context) (int, error)
	GetActiveConversations(ctx context.Context) (int, error)
	GetAverageResponseTime(ctx context.Context, days int) (float64, error)

	// Performance metrics
	GetTopKeywords(ctx context.Context, limit int, days int) ([]map[string]interface{}, error)
	GetResponseTypeDistribution(ctx context.Context, days int) (map[string]int, error)
	GetHourlyMessageDistribution(ctx context.Context, days int) ([]map[string]interface{}, error)
}

// Analytics repository interfaces and methods are already defined above
