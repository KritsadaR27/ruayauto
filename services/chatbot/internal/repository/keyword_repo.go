package repository

import (
	"context"
	"database/sql"
	"strings"

	"chatbot/internal/models"
)

type KeywordRepository struct {
	db *sql.DB
}

func NewKeywordRepository(db *sql.DB) *KeywordRepository {
	return &KeywordRepository{db: db}
}

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
			&keyword.ID, &keyword.Keyword, &keyword.Response,
			&keyword.IsActive, &keyword.Priority, &keyword.MatchType,
			&keyword.CreatedAt, &keyword.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		keywords = append(keywords, keyword)
	}
	return keywords, rows.Err()
}

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
			if strings.Contains(content, keywordText) {
				matching = append(matching, keyword)
			}
		}
	}
	return matching, nil
}

func (r *KeywordRepository) Create(ctx context.Context, keyword *models.Keyword) error {
	query := `
		INSERT INTO keywords (keyword, response, is_active, priority, match_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	
	err := r.db.QueryRowContext(
		ctx, query,
		keyword.Keyword, keyword.Response, keyword.IsActive, 
		keyword.Priority, keyword.MatchType,
	).Scan(&keyword.ID, &keyword.CreatedAt, &keyword.UpdatedAt)
	
	return err
}

func (r *KeywordRepository) Update(ctx context.Context, keyword *models.Keyword) error {
	query := `
		UPDATE keywords 
		SET keyword = $2, response = $3, is_active = $4, priority = $5, match_type = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`
	
	err := r.db.QueryRowContext(
		ctx, query,
		keyword.ID, keyword.Keyword, keyword.Response, 
		keyword.IsActive, keyword.Priority, keyword.MatchType,
	).Scan(&keyword.UpdatedAt)
	
	return err
}

func (r *KeywordRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM keywords WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *KeywordRepository) DeleteAll(ctx context.Context) error {
	query := `DELETE FROM keywords`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *KeywordRepository) GetByID(ctx context.Context, id int) (*models.Keyword, error) {
	query := `
		SELECT id, keyword, response, is_active, priority, match_type, created_at, updated_at
		FROM keywords 
		WHERE id = $1
	`
	
	keyword := &models.Keyword{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&keyword.ID, &keyword.Keyword, &keyword.Response,
		&keyword.IsActive, &keyword.Priority, &keyword.MatchType,
		&keyword.CreatedAt, &keyword.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return keyword, nil
}
