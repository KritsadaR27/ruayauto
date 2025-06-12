package repository

import (
"context"
"ruayautomsg/internal/models"
)

// KeywordRepository defines the interface for keyword data operations
type KeywordRepository interface {
	// GetAll retrieves all active keywords
	GetAll(ctx context.Context) ([]models.Keyword, error)
	
	// Create creates a new keyword
	Create(ctx context.Context, keyword models.Keyword) error
	
	// Update updates an existing keyword by ID
	Update(ctx context.Context, id int, keyword models.Keyword) error
	
	// Delete soft deletes a keyword by ID
	Delete(ctx context.Context, id int) error
	
	// GetByID retrieves a keyword by its ID
	GetByID(ctx context.Context, id int) (*models.Keyword, error)
	
	// GetByKeyword retrieves a keyword by its keyword string
	GetByKeyword(ctx context.Context, keyword string) (*models.Keyword, error)
	
	// UpsertKeyword creates or updates a keyword
	UpsertKeyword(ctx context.Context, keyword, response string) error
	
	// SearchKeywords searches keywords by partial match
	SearchKeywords(ctx context.Context, searchTerm string) ([]models.Keyword, error)
}
