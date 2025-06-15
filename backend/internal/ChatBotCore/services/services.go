package service

import "ruaychatbot/backend/internal/ChatBotCore/repository"

// Repositories struct to match handler expectations
type Repositories struct {
	Conversation     repository.ConversationRepository
	Message          repository.MessageRepository
	Keyword          repository.KeywordRepository
	ResponseTemplate repository.ResponseTemplateRepository
	Analytics        repository.AnalyticsRepository
}

// NewServiceRepositories creates service-level repositories wrapper
func NewServiceRepositories(repos *repository.Repositories) *Repositories {
	return &Repositories{
		Conversation:     repos.Conversation,
		Message:          repos.Message,
		Keyword:          repos.Keyword,
		ResponseTemplate: repos.ResponseTemplate,
		Analytics:        repos.Analytics,
	}
}
