package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"ruaymanagement/backend/internal/ChatBotCore/database"
	"ruaymanagement/backend/internal/ChatBotCore/models"
)

// keywordRepository implements KeywordRepository interface using PostgreSQL with pgx
type keywordRepository struct {
	db *database.DB
}

// NewKeywordRepository creates a new keyword repository instance
func NewKeywordRepository(db *database.DB) KeywordRepository {
	return &keywordRepository{
		db: db,
	}
}

// GetAll retrieves all active keywords from the database
func (r *keywordRepository) GetAll(ctx context.Context) ([]models.SimpleKeyword, error) {
	query := `
		SELECT id, keyword, response, is_active, match_type, priority, created_at, updated_at 
		FROM keywords 
		WHERE is_active = true 
		ORDER BY priority DESC, keyword ASC
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query keywords: %w", err)
	}
	defer rows.Close()

	var keywords []models.SimpleKeyword
	for rows.Next() {
		var k models.SimpleKeyword
		err := rows.Scan(&k.ID, &k.Keyword, &k.Response, &k.IsActive, &k.MatchType, &k.Priority, &k.CreatedAt, &k.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan keyword: %w", err)
		}
		keywords = append(keywords, k)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating keywords: %w", err)
	}

	return keywords, nil
}

// Create creates a new keyword in the database
func (r *keywordRepository) Create(ctx context.Context, keyword models.SimpleKeyword) error {
	query := `
		INSERT INTO keywords (keyword, response, is_active, match_type, priority, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id
	`

	now := time.Now()
	if keyword.CreatedAt.IsZero() {
		keyword.CreatedAt = now
	}
	if keyword.UpdatedAt.IsZero() {
		keyword.UpdatedAt = now
	}
	if !keyword.IsActive {
		keyword.IsActive = true
	}
	if keyword.MatchType == "" {
		keyword.MatchType = "exact"
	}

	err := r.db.Pool.QueryRow(ctx, query, 
		keyword.Keyword, 
		keyword.Response, 
		keyword.IsActive,
		keyword.MatchType,
		keyword.Priority,
		keyword.CreatedAt,
		keyword.UpdatedAt,
	).Scan(&keyword.ID)
	
	if err != nil {
		return fmt.Errorf("failed to create keyword: %w", err)
	}

	return nil
}

// Update updates an existing keyword by ID
func (r *keywordRepository) Update(ctx context.Context, id int, keyword models.SimpleKeyword) error {
	query := `
		UPDATE keywords 
		SET keyword = $2, response = $3, is_active = $4, match_type = $5, priority = $6, updated_at = $7
		WHERE id = $1 AND is_active = true
	`

	keyword.UpdatedAt = time.Now()
	
	commandTag, err := r.db.Pool.Exec(ctx, query, 
		id,
		keyword.Keyword, 
		keyword.Response, 
		keyword.IsActive,
		keyword.MatchType,
		keyword.Priority,
		keyword.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update keyword: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("keyword with id %d not found or already deleted", id)
	}

	return nil
}

// Delete soft deletes a keyword by ID
func (r *keywordRepository) Delete(ctx context.Context, id int) error {
	query := `
		UPDATE keywords 
		SET is_active = false, updated_at = $2
		WHERE id = $1 AND is_active = true
	`

	commandTag, err := r.db.Pool.Exec(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete keyword: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("keyword with id %d not found or already deleted", id)
	}

	return nil
}

// GetByID retrieves a keyword by its ID
func (r *keywordRepository) GetByID(ctx context.Context, id int) (*models.SimpleKeyword, error) {
	query := `
		SELECT id, keyword, response, is_active, match_type, priority, created_at, updated_at 
		FROM keywords 
		WHERE id = $1 AND is_active = true
	`

	var k models.SimpleKeyword
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&k.ID, &k.Keyword, &k.Response, &k.IsActive, &k.MatchType, &k.Priority, &k.CreatedAt, &k.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("keyword with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get keyword by id: %w", err)
	}

	return &k, nil
}

// GetByKeyword retrieves a keyword by its keyword string
func (r *keywordRepository) GetByKeyword(ctx context.Context, keyword string) (*models.SimpleKeyword, error) {
	query := `
		SELECT id, keyword, response, is_active, match_type, priority, created_at, updated_at 
		FROM keywords 
		WHERE keyword = $1 AND is_active = true
	`

	var k models.SimpleKeyword
	err := r.db.Pool.QueryRow(ctx, query, keyword).Scan(
		&k.ID, &k.Keyword, &k.Response, &k.IsActive, &k.MatchType, &k.Priority, &k.CreatedAt, &k.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Not found, return nil without error
		}
		return nil, fmt.Errorf("failed to get keyword by keyword: %w", err)
	}

	return &k, nil
}

// UpsertKeyword creates or updates a keyword
func (r *keywordRepository) UpsertKeyword(ctx context.Context, keyword, response string) error {
	query := `
		INSERT INTO keywords (keyword, response, is_active, created_at, updated_at) 
		VALUES ($1, $2, true, $3, $4) 
		ON CONFLICT (keyword) 
		DO UPDATE SET 
			response = EXCLUDED.response,
			is_active = true,
			updated_at = EXCLUDED.updated_at
	`

	now := time.Now()
	_, err := r.db.Pool.Exec(ctx, query, keyword, response, now, now)
	if err != nil {
		return fmt.Errorf("failed to upsert keyword: %w", err)
	}

	return nil
}

// SearchKeywords searches keywords by partial match
func (r *keywordRepository) SearchKeywords(ctx context.Context, searchTerm string) ([]models.SimpleKeyword, error) {
	query := `
		SELECT id, keyword, response, is_active, match_type, priority, created_at, updated_at 
		FROM keywords 
		WHERE is_active = true AND keyword ILIKE $1
		ORDER BY priority DESC, keyword ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, "%"+searchTerm+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search keywords: %w", err)
	}
	defer rows.Close()

	var keywords []models.SimpleKeyword
	for rows.Next() {
		var k models.SimpleKeyword
		err := rows.Scan(&k.ID, &k.Keyword, &k.Response, &k.IsActive, &k.MatchType, &k.Priority, &k.CreatedAt, &k.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan keyword during search: %w", err)
		}
		keywords = append(keywords, k)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating search results: %w", err)
	}

	return keywords, nil
}

// FindMatchingKeywords finds keywords that match the given text
func (r *keywordRepository) FindMatchingKeywords(ctx context.Context, text string) ([]models.SimpleKeyword, error) {
	if strings.TrimSpace(text) == "" {
		return []models.SimpleKeyword{}, nil
	}

	keywords, err := r.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get keywords: %w", err)
	}

	// Convert text to lowercase for case-insensitive matching
	lowerText := strings.ToLower(text)
	var matches []models.SimpleKeyword

	// Find keywords that are contained in the text
	for _, keyword := range keywords {
		if strings.Contains(lowerText, strings.ToLower(keyword.Keyword)) {
			matches = append(matches, keyword)
		}
	}

	return matches, nil
}

// SetActive sets the active status of a keyword
func (r *keywordRepository) SetActive(ctx context.Context, id int, isActive bool) error {
	query := `
		UPDATE keywords 
		SET is_active = $2, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query, id, isActive)
	if err != nil {
		return fmt.Errorf("failed to set keyword active status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("keyword with id %d not found", id)
	}

	return nil
}

// GetActive retrieves all active keywords
func (r *keywordRepository) GetActive(ctx context.Context) ([]models.SimpleKeyword, error) {
	return r.GetAll(ctx) // GetAll already filters for active keywords
}

// List retrieves keywords with pagination
func (r *keywordRepository) List(ctx context.Context, offset, limit int) ([]models.SimpleKeyword, error) {
	query := `
		SELECT id, keyword, response, is_active, match_type, priority, created_at, updated_at 
		FROM keywords 
		WHERE is_active = true 
		ORDER BY priority DESC, keyword ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query keywords: %w", err)
	}
	defer rows.Close()

	var keywords []models.SimpleKeyword
	for rows.Next() {
		var k models.SimpleKeyword
		err := rows.Scan(&k.ID, &k.Keyword, &k.Response, &k.IsActive, &k.MatchType, &k.Priority, &k.CreatedAt, &k.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan keyword: %w", err)
		}
		keywords = append(keywords, k)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating keywords: %w", err)
	}

	return keywords, nil
}

// Count returns the total number of active keywords
func (r *keywordRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM keywords WHERE is_active = true`

	var count int
	err := r.db.Pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count keywords: %w", err)
	}

	return count, nil
}
