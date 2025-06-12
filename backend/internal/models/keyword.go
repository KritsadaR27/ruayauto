package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"ruayautomsg/internal/database"
)

// Keyword represents a keyword-response pair in the database
type Keyword struct {
	ID        int       `json:"id" db:"id"`
	Keyword   string    `json:"keyword" db:"keyword"`
	Response  string    `json:"response" db:"response"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// KeywordRepository handles database operations for keywords
type KeywordRepository struct {
	db *database.DB
}

// NewKeywordRepository creates a new keyword repository
func NewKeywordRepository(db *database.DB) *KeywordRepository {
	return &KeywordRepository{
		db: db,
	}
}

// GetAll retrieves all active keywords from the database
func (r *KeywordRepository) GetAll(ctx context.Context) ([]Keyword, error) {
	query := `
		SELECT id, keyword, response, is_active, created_at, updated_at 
		FROM keywords 
		WHERE is_active = true 
		ORDER BY keyword ASC
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query keywords: %w", err)
	}
	defer rows.Close()

	var keywords []Keyword
	for rows.Next() {
		var k Keyword
		err := rows.Scan(&k.ID, &k.Keyword, &k.Response, &k.IsActive, &k.CreatedAt, &k.UpdatedAt)
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

// GetByKeyword retrieves a keyword by its keyword string
func (r *KeywordRepository) GetByKeyword(ctx context.Context, keyword string) (*Keyword, error) {
	query := `
		SELECT id, keyword, response, is_active, created_at, updated_at 
		FROM keywords 
		WHERE keyword = $1 AND is_active = true
	`

	var k Keyword
	err := r.db.Pool.QueryRow(ctx, query, keyword).Scan(
		&k.ID, &k.Keyword, &k.Response, &k.IsActive, &k.CreatedAt, &k.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to get keyword: %w", err)
	}

	return &k, nil
}

// Create creates a new keyword in the database
func (r *KeywordRepository) Create(ctx context.Context, keyword *Keyword) error {
	query := `
		INSERT INTO keywords (keyword, response, is_active) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at
	`

	err := r.db.Pool.QueryRow(ctx, query, keyword.Keyword, keyword.Response, keyword.IsActive).Scan(
		&keyword.ID, &keyword.CreatedAt, &keyword.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create keyword: %w", err)
	}

	return nil
}

// Update updates an existing keyword in the database
func (r *KeywordRepository) Update(ctx context.Context, keyword *Keyword) error {
	query := `
		UPDATE keywords 
		SET keyword = $1, response = $2, is_active = $3 
		WHERE id = $4 
		RETURNING updated_at
	`

	err := r.db.Pool.QueryRow(ctx, query, keyword.Keyword, keyword.Response, keyword.IsActive, keyword.ID).Scan(
		&keyword.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update keyword: %w", err)
	}

	return nil
}

// Delete soft deletes a keyword by setting is_active to false
func (r *KeywordRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE keywords SET is_active = false WHERE id = $1`

	err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete keyword: %w", err)
	}

	return nil
}

// UpsertKeyword creates or updates a keyword
func (r *KeywordRepository) UpsertKeyword(ctx context.Context, keyword, response string) error {
	query := `
		INSERT INTO keywords (keyword, response, is_active) 
		VALUES ($1, $2, true) 
		ON CONFLICT (keyword) 
		DO UPDATE SET 
			response = EXCLUDED.response,
			is_active = true,
			updated_at = CURRENT_TIMESTAMP
	`

	err := r.db.Exec(ctx, query, keyword, response)
	if err != nil {
		return fmt.Errorf("failed to upsert keyword: %w", err)
	}

	return nil
}

// SearchKeywords searches keywords by partial match
func (r *KeywordRepository) SearchKeywords(ctx context.Context, searchTerm string) ([]Keyword, error) {
	query := `
		SELECT id, keyword, response, is_active, created_at, updated_at 
		FROM keywords 
		WHERE is_active = true AND keyword ILIKE $1
		ORDER BY keyword ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, "%"+searchTerm+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search keywords: %w", err)
	}
	defer rows.Close()

	var keywords []Keyword
	for rows.Next() {
		var k Keyword
		err := rows.Scan(&k.ID, &k.Keyword, &k.Response, &k.IsActive, &k.CreatedAt, &k.UpdatedAt)
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
