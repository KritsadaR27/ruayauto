package repository

import (
"context"
"fmt"
"time"

"github.com/jackc/pgx/v5"
"ruayautomsg/internal/database"
"ruayautomsg/internal/models"
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
func (r *keywordRepository) GetAll(ctx context.Context) ([]models.Keyword, error) {
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

	var keywords []models.Keyword
	for rows.Next() {
		var k models.Keyword
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

// Create creates a new keyword in the database
func (r *keywordRepository) Create(ctx context.Context, keyword models.Keyword) error {
	query := `
		INSERT INTO keywords (keyword, response, is_active, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5) 
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

	err := r.db.Pool.QueryRow(ctx, query, 
keyword.Keyword, 
keyword.Response, 
keyword.IsActive,
keyword.CreatedAt,
keyword.UpdatedAt,
).Scan(&keyword.ID)
	
	if err != nil {
		return fmt.Errorf("failed to create keyword: %w", err)
	}

	return nil
}

// Update updates an existing keyword by ID
func (r *keywordRepository) Update(ctx context.Context, id int, keyword models.Keyword) error {
	query := `
		UPDATE keywords 
		SET keyword = $2, response = $3, is_active = $4, updated_at = $5
		WHERE id = $1 AND is_active = true
	`

	keyword.UpdatedAt = time.Now()
	
	commandTag, err := r.db.Pool.Exec(ctx, query, 
id,
keyword.Keyword, 
keyword.Response, 
keyword.IsActive,
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
func (r *keywordRepository) GetByID(ctx context.Context, id int) (*models.Keyword, error) {
	query := `
		SELECT id, keyword, response, is_active, created_at, updated_at 
		FROM keywords 
		WHERE id = $1 AND is_active = true
	`

	var k models.Keyword
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
&k.ID, &k.Keyword, &k.Response, &k.IsActive, &k.CreatedAt, &k.UpdatedAt,
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
func (r *keywordRepository) GetByKeyword(ctx context.Context, keyword string) (*models.Keyword, error) {
	query := `
		SELECT id, keyword, response, is_active, created_at, updated_at 
		FROM keywords 
		WHERE keyword = $1 AND is_active = true
	`

	var k models.Keyword
	err := r.db.Pool.QueryRow(ctx, query, keyword).Scan(
&k.ID, &k.Keyword, &k.Response, &k.IsActive, &k.CreatedAt, &k.UpdatedAt,
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
func (r *keywordRepository) SearchKeywords(ctx context.Context, searchTerm string) ([]models.Keyword, error) {
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

	var keywords []models.Keyword
	for rows.Next() {
		var k models.Keyword
		err := rows.Scan(&k.ID, &k.Keyword, &k.Response, &k.IsActive, &k.CreatedAt, &k.UpdatedAt)
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
