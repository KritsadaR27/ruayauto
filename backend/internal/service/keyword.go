package service

import (
"context"
"fmt"
"strings"

"ruayautomsg/internal/models"
"ruayautomsg/internal/repository"
)

// KeywordService provides business logic for keyword operations
type KeywordService struct {
	repo repository.KeywordRepository
}

// NewKeywordService creates a new keyword service
func NewKeywordService(repo repository.KeywordRepository) *KeywordService {
	return &KeywordService{
		repo: repo,
	}
}

// GetRepository returns the underlying repository (for integration with other services)
func (s *KeywordService) GetRepository() repository.KeywordRepository {
	return s.repo
}

// GetAllKeywords retrieves all active keywords
func (s *KeywordService) GetAllKeywords(ctx context.Context) ([]models.Keyword, error) {
	return s.repo.GetAll(ctx)
}

// CreateKeyword creates a new keyword with validation
func (s *KeywordService) CreateKeyword(ctx context.Context, keyword, response string) error {
	// Validate input
	if err := s.validateKeywordInput(keyword, response); err != nil {
		return err
	}

	// Check if keyword already exists
	existing, err := s.repo.GetByKeyword(ctx, keyword)
	if err != nil {
		return fmt.Errorf("failed to check existing keyword: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("keyword '%s' already exists", keyword)
	}

	// Create new keyword
	newKeyword := models.Keyword{
		Keyword:  strings.TrimSpace(keyword),
		Response: strings.TrimSpace(response),
		IsActive: true,
	}

	return s.repo.Create(ctx, newKeyword)
}

// UpdateKeyword updates an existing keyword
func (s *KeywordService) UpdateKeyword(ctx context.Context, id int, keyword, response string) error {
	// Validate input
	if err := s.validateKeywordInput(keyword, response); err != nil {
		return err
	}

	// Get existing keyword
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get keyword: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("keyword with ID %d not found", id)
	}

	// Update fields
	existing.Keyword = strings.TrimSpace(keyword)
	existing.Response = strings.TrimSpace(response)

return s.repo.Update(ctx, existing.ID, *existing)
}

// DeleteKeyword soft deletes a keyword
func (s *KeywordService) DeleteKeyword(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

// GetKeywordByID retrieves a keyword by ID
func (s *KeywordService) GetKeywordByID(ctx context.Context, id int) (*models.Keyword, error) {
	return s.repo.GetByID(ctx, id)
}

// SearchKeywords searches for keywords containing the search term
func (s *KeywordService) SearchKeywords(ctx context.Context, searchTerm string) ([]models.Keyword, error) {
	return s.repo.SearchKeywords(ctx, searchTerm)
}

// UpsertKeyword creates or updates a keyword
func (s *KeywordService) UpsertKeyword(ctx context.Context, keyword, response string) error {
	return s.repo.UpsertKeyword(ctx, keyword, response)
}

// validateKeywordInput validates keyword and response input
func (s *KeywordService) validateKeywordInput(keyword, response string) error {
	keyword = strings.TrimSpace(keyword)
	response = strings.TrimSpace(response)

	if keyword == "" {
		return fmt.Errorf("keyword cannot be empty")
	}
	if response == "" {
		return fmt.Errorf("response cannot be empty")
	}
	if len(keyword) > 255 {
		return fmt.Errorf("keyword too long: maximum 255 characters, got %d", len(keyword))
	}
	if len(response) > 10000 {
		return fmt.Errorf("response too long: maximum 10000 characters, got %d", len(response))
	}

	// Check for invalid characters in keyword
	if strings.ContainsAny(keyword, "\n\r\t") {
		return fmt.Errorf("keyword cannot contain newlines or tabs")
	}

	return nil
}

// FindMatchingKeyword finds a keyword that matches the given text (legacy compatibility)
func (s *KeywordService) FindMatchingKeyword(ctx context.Context, text string) (*models.Keyword, error) {
	if strings.TrimSpace(text) == "" {
		return nil, nil
	}

	keywords, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get keywords: %w", err)
	}

	// Convert text to lowercase for case-insensitive matching
	lowerText := strings.ToLower(text)

	// First try exact substring matching
	for _, keyword := range keywords {
		if strings.Contains(lowerText, strings.ToLower(keyword.Keyword)) {
			return &keyword, nil
		}
	}

	// No match found
	return nil, nil
}

// GetKeywordPairs returns keywords in the legacy format for backward compatibility
func (s *KeywordService) GetKeywordPairs(ctx context.Context) ([]map[string]interface{}, error) {
	keywords, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get keywords: %w", err)
	}

	pairs := make([]map[string]interface{}, 0, len(keywords))
	for _, keyword := range keywords {
		pairs = append(pairs, map[string]interface{}{
"keywords":  []string{keyword.Keyword},
"responses": []string{keyword.Response},
})
	}

	return pairs, nil
}

// BulkUpdateKeywords updates multiple keywords from legacy format
func (s *KeywordService) BulkUpdateKeywords(ctx context.Context, pairs []struct {
Keywords  []string `json:"keywords"`
Responses []string `json:"responses"`
}) error {
	// Validate input
	for i, pair := range pairs {
		if len(pair.Keywords) == 0 {
			return fmt.Errorf("pair %d: keywords cannot be empty", i)
		}
		if len(pair.Responses) == 0 {
			return fmt.Errorf("pair %d: responses cannot be empty", i)
		}
	}

	// Get existing keywords
	existing, err := s.repo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get existing keywords: %w", err)
	}

	// Create a map of existing keywords for quick lookup
	existingMap := make(map[string]models.Keyword)
	for _, kw := range existing {
		existingMap[kw.Keyword] = kw
	}

	// Process each pair
	for _, pair := range pairs {
		for _, keyword := range pair.Keywords {
			keyword = strings.TrimSpace(keyword)
			if keyword == "" {
				continue
			}

			// Use first response for simplicity
			response := pair.Responses[0]
			if len(pair.Responses) > 1 {
				// For multiple responses, join them or pick randomly
				// For now, just use the first one
				response = pair.Responses[0]
			}

			if existing, exists := existingMap[keyword]; exists {
				// Update existing keyword
				existing.Response = response
				if err := s.repo.Update(ctx, existing.ID, existing); err != nil {
					return fmt.Errorf("failed to update keyword '%s': %w", keyword, err)
				}
			} else {
				// Create new keyword
				newKeyword := models.Keyword{
					Keyword:  keyword,
					Response: response,
					IsActive: true,
				}
				if err := s.repo.Create(ctx, newKeyword); err != nil {
					return fmt.Errorf("failed to create keyword '%s': %w", keyword, err)
				}
			}
		}
	}

	return nil
}
