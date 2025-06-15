package repository

import (
	"context"
	"database/sql"
	"strings"

	"chatbot-core/models"
)

// KeywordRepository handles keyword-related database operations
type KeywordRepository struct {
	db *sql.DB
}

// NewKeywordRepository creates a new keyword repository
func NewKeywordRepository(db *sql.DB) *KeywordRepository {
	return &KeywordRepository{db: db}
}

// GetActive retrieves all active keywords
func (r *KeywordRepository) GetActive(ctx context.Context) ([]*models.Keyword, error) {
	query := `
		SELECT id, keyword, response, is_active, priority, match_type, created_at, updated_at
		FROM keywords 
		WHERE is_active = true
		ORDER BY priority DESC, created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keywords []*models.Keyword
	for rows.Next() {
		keyword := &models.Keyword{}
		err := rows.Scan(
			&keyword.ID,
			&keyword.Keyword,
			&keyword.Response,
			&keyword.IsActive,
			&keyword.Priority,
			&keyword.MatchType,
			&keyword.CreatedAt,
			&keyword.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		keywords = append(keywords, keyword)
	}

	return keywords, rows.Err()
}

// FindMatching finds keywords that match the given content
func (r *KeywordRepository) FindMatching(ctx context.Context, content string) ([]*models.Keyword, error) {
	keywords, err := r.GetActive(ctx)
	if err != nil {
		return nil, err
	}

	content = strings.ToLower(strings.TrimSpace(content))
	var matching []*models.Keyword

	for _, keyword := range keywords {
		keywordText := strings.ToLower(strings.TrimSpace(keyword.Keyword))

		switch keyword.MatchType {
		case "exact":
			if content == keywordText {
				matching = append(matching, keyword)
			}
		case "contains":
			if strings.Contains(content, keywordText) {
				matching = append(matching, keyword)
			}
		case "starts_with":
			if strings.HasPrefix(content, keywordText) {
				matching = append(matching, keyword)
			}
		case "ends_with":
			if strings.HasSuffix(content, keywordText) {
				matching = append(matching, keyword)
			}
		default:
			// Default to contains
			if strings.Contains(content, keywordText) {
				matching = append(matching, keyword)
			}
		}
	}

	return matching, nil
}
