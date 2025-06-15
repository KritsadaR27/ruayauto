package repository

import (
	"ruaychatbot/backend/internal/ChatBotCore/database"
)

// Repositories holds all repository instances
type Repositories struct {
	Conversation     ConversationRepository
	Message          MessageRepository
	Keyword          KeywordRepository
	ResponseTemplate ResponseTemplateRepository
	Analytics        AnalyticsRepository
}

// NewRepositories creates a new Repositories instance with all repositories
func NewRepositories(db *database.DB) *Repositories {
	return &Repositories{
		Conversation:     NewConversationRepository(db),
		Message:          NewMessageRepository(db),
		Keyword:          NewKeywordRepository(db),
		ResponseTemplate: NewResponseTemplateRepository(db),
		Analytics:        NewAnalyticsRepository(db),
	}
}
