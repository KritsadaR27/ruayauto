package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"ruaymanagement/backend/internal/ChatBotCore/models"
)

type keywordRepository struct {
	db *sql.DB
}

func NewKeywordRepository(db *sql.DB) KeywordRepository {
	return &keywordRepository{db: db}
}

func (r *keywordRepository) Create(ctx context.Context, keyword *models.Keyword) error {
	query := `
		INSERT INTO chatbot_mvp.keywords (
			keyword_text, match_type, response_type, response_content, response_data, priority, is_active, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		keyword.Keyword,
		keyword.MatchType,
		keyword.ResponseType,
		keyword.ResponseContent,
		keyword.ResponseData,
		keyword.Priority,
		keyword.IsActive,
	).Scan(&keyword.ID, &keyword.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create keyword: %w", err)
	}

	return nil
}

func (r *keywordRepository) GetByID(ctx context.Context, id int) (*models.Keyword, error) {
	keyword := &models.Keyword{}
	query := `
		SELECT id, keyword_text, match_type, response_type, response_content, response_data, priority, 
			   is_active, created_at, updated_at
		FROM chatbot_mvp.keywords
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&keyword.ID,
		&keyword.Keyword,
		&keyword.MatchType,
		&keyword.ResponseType,
		&keyword.ResponseContent,
		&keyword.ResponseData,
		&keyword.Priority,
		&keyword.IsActive,
		&keyword.CreatedAt,
		&keyword.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("keyword not found")
		}
		return nil, fmt.Errorf("failed to get keyword: %w", err)
	}

	return keyword, nil
}

func (r *keywordRepository) List(ctx context.Context, offset, limit int) ([]models.Keyword, error) {
	query := `
		SELECT id, keyword_text, match_type, response_type, response_content, response_data, priority, 
			   is_active, created_at, updated_at
		FROM chatbot_mvp.keywords
		ORDER BY priority DESC, id ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list keywords: %w", err)
	}
	defer rows.Close()

	var keywords []models.Keyword
	for rows.Next() {
		keyword := models.Keyword{}
		err := rows.Scan(
			&keyword.ID,
			&keyword.Keyword,
			&keyword.MatchType,
			&keyword.ResponseType,
			&keyword.ResponseContent,
			&keyword.ResponseData,
			&keyword.Priority,
			&keyword.IsActive,
			&keyword.CreatedAt,
			&keyword.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan keyword: %w", err)
		}
		keywords = append(keywords, keyword)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating keywords: %w", err)
	}

	return keywords, nil
}

func (r *keywordRepository) GetActive(ctx context.Context) ([]models.Keyword, error) {
	query := `
		SELECT id, keyword_text, match_type, response_type, response_content, response_data, priority, 
			   is_active, created_at, updated_at
		FROM chatbot_mvp.keywords
		WHERE is_active = true
		ORDER BY priority DESC, id ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active keywords: %w", err)
	}
	defer rows.Close()

	var keywords []models.Keyword
	for rows.Next() {
		keyword := models.Keyword{}
		err := rows.Scan(
			&keyword.ID,
			&keyword.Keyword,
			&keyword.MatchType,
			&keyword.ResponseType,
			&keyword.ResponseContent,
			&keyword.ResponseData,
			&keyword.Priority,
			&keyword.IsActive,
			&keyword.CreatedAt,
			&keyword.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan keyword: %w", err)
		}
		keywords = append(keywords, keyword)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating keywords: %w", err)
	}

	return keywords, nil
}

func (r *keywordRepository) Update(ctx context.Context, keyword *models.Keyword) error {
	query := `
		UPDATE chatbot_mvp.keywords
		SET keyword_text = $2, match_type = $3, response_type = $4, response_content = $5, 
			response_data = $6, priority = $7, is_active = $8, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(
		ctx, query,
		keyword.ID,
		keyword.Keyword,
		keyword.MatchType,
		keyword.ResponseType,
		keyword.ResponseContent,
		keyword.ResponseData,
		keyword.Priority,
		keyword.IsActive,
	)

	if err != nil {
		return fmt.Errorf("failed to update keyword: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("keyword not found")
	}

	return nil
}

func (r *keywordRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM chatbot_mvp.keywords WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete keyword: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("keyword not found")
	}

	return nil
}

func (r *keywordRepository) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM chatbot_mvp.keywords`

	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count keywords: %w", err)
	}

	return count, nil
}

// FindMatchingKeywords finds keywords that match the given text
func (r *keywordRepository) FindMatchingKeywords(ctx context.Context, messageText string) ([]models.Keyword, error) {
	// Get all active keywords ordered by priority
	keywords, err := r.GetActive(ctx)
	if err != nil {
		return nil, err
	}

	var matchedKeywords []models.Keyword
	messageTextLower := strings.ToLower(messageText)

	for _, keyword := range keywords {
		matched := false
		keywordTextLower := strings.ToLower(keyword.Keyword)

		switch keyword.MatchType {
		case "exact":
			matched = messageTextLower == keywordTextLower
		case "contains":
			matched = strings.Contains(messageTextLower, keywordTextLower)
		case "starts_with":
			matched = strings.HasPrefix(messageTextLower, keywordTextLower)
		case "ends_with":
			matched = strings.HasSuffix(messageTextLower, keywordTextLower)
		case "word_boundary":
			// Simple word boundary check - could be enhanced with regex
			words := strings.Fields(messageTextLower)
			for _, word := range words {
				if word == keywordTextLower {
					matched = true
					break
				}
			}
		}

		if matched {
			matchedKeywords = append(matchedKeywords, keyword)
		}
	}

	return matchedKeywords, nil
}

func (r *keywordRepository) SetActive(ctx context.Context, id int, isActive bool) error {
	query := `
		UPDATE chatbot_mvp.keywords
		SET is_active = $2, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, isActive)
	if err != nil {
		return fmt.Errorf("failed to set keyword active status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("keyword not found")
	}

	return nil
}
