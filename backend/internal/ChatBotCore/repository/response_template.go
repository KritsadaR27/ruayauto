package repository

import (
	"context"
	"fmt"

	"ruaychatbot/backend/internal/ChatBotCore/database"
	"ruaychatbot/backend/internal/ChatBotCore/models"

	"github.com/jackc/pgx/v5"
)

type responseTemplateRepository struct {
	db *database.DB
}

func NewResponseTemplateRepository(db *database.DB) ResponseTemplateRepository {
	return &responseTemplateRepository{db: db}
}

func (r *responseTemplateRepository) Create(ctx context.Context, template *models.ResponseTemplate) error {
	query := `
		INSERT INTO chatbot_mvp.response_templates (
			template_name, template_type, template_data, description, is_active, created_at
		) VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at
	`

	err := r.db.Pool.QueryRow(
		ctx, query,
		template.Name,
		template.TemplateType,
		template.TemplateData,
		template.Description,
		template.IsActive,
	).Scan(&template.ID, &template.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create response template: %w", err)
	}

	return nil
}

func (r *responseTemplateRepository) GetByID(ctx context.Context, id int) (*models.ResponseTemplate, error) {
	template := &models.ResponseTemplate{}
	query := `
		SELECT id, template_name, template_type, template_data, description, is_active, 
			   created_at, updated_at
		FROM chatbot_mvp.response_templates
		WHERE id = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&template.ID,
		&template.Name,
		&template.TemplateType,
		&template.TemplateData,
		&template.Description,
		&template.IsActive,
		&template.CreatedAt,
		&template.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("response template not found")
		}
		return nil, fmt.Errorf("failed to get response template: %w", err)
	}

	return template, nil
}

// GetByName retrieves a response template by name
func (r *responseTemplateRepository) GetByName(ctx context.Context, name string) (*models.ResponseTemplate, error) {
	template := &models.ResponseTemplate{}
	query := `
		SELECT id, template_name, template_type, template_data, description, is_active, 
			   created_at, updated_at
		FROM chatbot_mvp.response_templates
		WHERE template_name = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, name).Scan(
		&template.ID,
		&template.Name,
		&template.TemplateType,
		&template.TemplateData,
		&template.Description,
		&template.IsActive,
		&template.CreatedAt,
		&template.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("response template with name '%s' not found", name)
		}
		return nil, fmt.Errorf("failed to get response template by name: %w", err)
	}

	return template, nil
}

func (r *responseTemplateRepository) List(ctx context.Context, offset, limit int) ([]models.ResponseTemplate, error) {
	query := `
		SELECT id, template_name, template_type, template_data, description, is_active, 
			   created_at, updated_at
		FROM chatbot_mvp.response_templates
		ORDER BY template_name ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list response templates: %w", err)
	}
	defer rows.Close()

	var templates []models.ResponseTemplate
	for rows.Next() {
		template := models.ResponseTemplate{}
		err := rows.Scan(
			&template.ID,
			&template.Name,
			&template.TemplateType,
			&template.TemplateData,
			&template.Description,
			&template.IsActive,
			&template.CreatedAt,
			&template.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan response template: %w", err)
		}
		templates = append(templates, template)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating response templates: %w", err)
	}

	return templates, nil
}

// GetActive retrieves all active response templates
func (r *responseTemplateRepository) GetActive(ctx context.Context) ([]models.ResponseTemplate, error) {
	query := `
		SELECT id, template_name, template_type, template_data, description, is_active, 
			   created_at, updated_at
		FROM chatbot_mvp.response_templates
		WHERE is_active = true
		ORDER BY created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active response templates: %w", err)
	}
	defer rows.Close()

	var templates []models.ResponseTemplate
	for rows.Next() {
		template := models.ResponseTemplate{}
		err := rows.Scan(
			&template.ID,
			&template.Name,
			&template.TemplateType,
			&template.TemplateData,
			&template.Description,
			&template.IsActive,
			&template.CreatedAt,
			&template.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan response template: %w", err)
		}
		templates = append(templates, template)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating response templates: %w", err)
	}

	return templates, nil
}

// GetByType retrieves response templates by type
func (r *responseTemplateRepository) GetByType(ctx context.Context, templateType string) ([]models.ResponseTemplate, error) {
	query := `
		SELECT id, template_name, template_type, template_data, description, is_active, 
			   created_at, updated_at
		FROM chatbot_mvp.response_templates
		WHERE template_type = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, templateType)
	if err != nil {
		return nil, fmt.Errorf("failed to get response templates by type: %w", err)
	}
	defer rows.Close()

	var templates []models.ResponseTemplate
	for rows.Next() {
		template := models.ResponseTemplate{}
		err := rows.Scan(
			&template.ID,
			&template.Name,
			&template.TemplateType,
			&template.TemplateData,
			&template.Description,
			&template.IsActive,
			&template.CreatedAt,
			&template.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan response template: %w", err)
		}
		templates = append(templates, template)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating response templates: %w", err)
	}

	return templates, nil
}

func (r *responseTemplateRepository) Update(ctx context.Context, template *models.ResponseTemplate) error {
	query := `
		UPDATE chatbot_mvp.response_templates
		SET template_name = $2, template_type = $3, template_data = $4, 
			description = $5, is_active = $6, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(
		ctx, query,
		template.ID,
		template.Name,
		template.TemplateType,
		template.TemplateData,
		template.Description,
		template.IsActive,
	)

	if err != nil {
		return fmt.Errorf("failed to update response template: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("response template not found")
	}

	return nil
}

func (r *responseTemplateRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM chatbot_mvp.response_templates WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete response template: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("response template not found")
	}

	return nil
}

func (r *responseTemplateRepository) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM chatbot_mvp.response_templates`

	err := r.db.Pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count response templates: %w", err)
	}

	return count, nil
}

func (r *responseTemplateRepository) SetActive(ctx context.Context, id int, isActive bool) error {
	query := `
		UPDATE chatbot_mvp.response_templates
		SET is_active = $2, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query, id, isActive)
	if err != nil {
		return fmt.Errorf("failed to set response template active status: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("response template not found")
	}

	return nil
}

// GetAll retrieves all response templates
func (r *responseTemplateRepository) GetAll(ctx context.Context) ([]models.ResponseTemplate, error) {
	query := `
		SELECT id, template_name, template_type, template_data, description, is_active, 
			   created_at, updated_at
		FROM chatbot_mvp.response_templates
		ORDER BY created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all response templates: %w", err)
	}
	defer rows.Close()

	var templates []models.ResponseTemplate
	for rows.Next() {
		template := models.ResponseTemplate{}
		err := rows.Scan(
			&template.ID,
			&template.Name,
			&template.TemplateType,
			&template.TemplateData,
			&template.Description,
			&template.IsActive,
			&template.CreatedAt,
			&template.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan response template: %w", err)
		}
		templates = append(templates, template)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating response templates: %w", err)
	}

	return templates, nil
}

// ToggleActive toggles the active status of a response template
func (r *responseTemplateRepository) ToggleActive(ctx context.Context, id int, isActive bool) error {
	query := `
		UPDATE chatbot_mvp.response_templates
		SET is_active = $2, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query, id, isActive)
	if err != nil {
		return fmt.Errorf("failed to toggle response template active status: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("response template not found")
	}

	return nil
}
