package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"ruaymanagement/backend/internal/ChatBotCore/database"
	"ruaymanagement/backend/internal/ChatBotCore/models"
)

// ConversationRepositoryImpl implements ConversationRepository
type ConversationRepositoryImpl struct {
	db *database.DB
}

// NewConversationRepository creates a new conversation repository
func NewConversationRepository(db *database.DB) ConversationRepository {
	return &ConversationRepositoryImpl{db: db}
}

// Create creates a new conversation
func (r *ConversationRepositoryImpl) Create(ctx context.Context, conv *models.Conversation) error {
	query := `
		INSERT INTO conversations (page_id, facebook_user_id, facebook_user_name, source_type, post_id, status, started_at, last_message_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	err := r.db.Pool.QueryRow(ctx, query, 
		conv.PageID, 
		conv.FacebookUserID, 
		conv.FacebookUserName, 
		conv.SourceType, 
		conv.PostID, 
		conv.Status, 
		conv.StartedAt, 
		conv.LastMessageAt,
	).Scan(&conv.ID, &conv.CreatedAt, &conv.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}

	return nil
}

// GetByID retrieves a conversation by ID
func (r *ConversationRepositoryImpl) GetByID(ctx context.Context, id int) (*models.Conversation, error) {
	query := `
		SELECT id, page_id, facebook_user_id, facebook_user_name, source_type, post_id, started_at, last_message_at, status, created_at, updated_at
		FROM conversations
		WHERE id = $1
	`

	var conv models.Conversation
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&conv.ID,
		&conv.PageID,
		&conv.FacebookUserID,
		&conv.FacebookUserName,
		&conv.SourceType,
		&conv.PostID,
		&conv.StartedAt,
		&conv.LastMessageAt,
		&conv.Status,
		&conv.CreatedAt,
		&conv.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation by ID: %w", err)
	}

	return &conv, nil
}

// GetByFacebookUserID retrieves a conversation by Facebook user ID
func (r *ConversationRepositoryImpl) GetByFacebookUserID(ctx context.Context, facebookUserID string) (*models.Conversation, error) {
	query := `
		SELECT id, page_id, facebook_user_id, facebook_user_name, source_type, post_id, started_at, last_message_at, status, created_at, updated_at
		FROM conversations
		WHERE facebook_user_id = $1
	`

	var conv models.Conversation
	err := r.db.Pool.QueryRow(ctx, query, facebookUserID).Scan(
		&conv.ID,
		&conv.PageID,
		&conv.FacebookUserID,
		&conv.FacebookUserName,
		&conv.SourceType,
		&conv.PostID,
		&conv.StartedAt,
		&conv.LastMessageAt,
		&conv.Status,
		&conv.CreatedAt,
		&conv.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation by Facebook user ID: %w", err)
	}

	return &conv, nil
}

// Update updates an existing conversation
func (r *ConversationRepositoryImpl) Update(ctx context.Context, conv *models.Conversation) error {
	query := `
		UPDATE conversations
		SET page_id = $2, facebook_user_name = $3, source_type = $4, post_id = $5, status = $6, last_message_at = $7, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.Pool.QueryRow(ctx, query, 
		conv.ID, 
		conv.PageID, 
		conv.FacebookUserName, 
		conv.SourceType, 
		conv.PostID, 
		conv.Status, 
		conv.LastMessageAt,
	).Scan(&conv.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update conversation: %w", err)
	}

	return nil
}

// Delete deletes a conversation
func (r *ConversationRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM conversations WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("conversation not found")
	}

	return nil
}

// List retrieves conversations with pagination
func (r *ConversationRepositoryImpl) List(ctx context.Context, offset, limit int) ([]models.Conversation, error) {
	query := `
		SELECT id, page_id, facebook_user_id, facebook_user_name, source_type, post_id, started_at, last_message_at, status, created_at, updated_at
		FROM conversations
		ORDER BY last_message_at DESC NULLS LAST, created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversations: %w", err)
	}
	defer rows.Close()

	var conversations []models.Conversation
	for rows.Next() {
		var conv models.Conversation
		err := rows.Scan(
			&conv.ID,
			&conv.PageID,
			&conv.FacebookUserID,
			&conv.FacebookUserName,
			&conv.SourceType,
			&conv.PostID,
			&conv.StartedAt,
			&conv.LastMessageAt,
			&conv.Status,
			&conv.CreatedAt,
			&conv.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}
		conversations = append(conversations, conv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating conversations: %w", err)
	}

	return conversations, nil
}

// ListByStatus retrieves conversations by status with pagination
func (r *ConversationRepositoryImpl) ListByStatus(ctx context.Context, status string, offset, limit int) ([]models.Conversation, error) {
	query := `
		SELECT id, page_id, facebook_user_id, facebook_user_name, source_type, post_id, started_at, last_message_at, status, created_at, updated_at
		FROM conversations
		WHERE status = $1
		ORDER BY last_message_at DESC NULLS LAST, created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, query, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversations by status: %w", err)
	}
	defer rows.Close()

	var conversations []models.Conversation
	for rows.Next() {
		var conv models.Conversation
		err := rows.Scan(
			&conv.ID,
			&conv.PageID,
			&conv.FacebookUserID,
			&conv.FacebookUserName,
			&conv.SourceType,
			&conv.PostID,
			&conv.StartedAt,
			&conv.LastMessageAt,
			&conv.Status,
			&conv.CreatedAt,
			&conv.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}
		conversations = append(conversations, conv)
	}

	return conversations, nil
}

// Count returns total number of conversations
func (r *ConversationRepositoryImpl) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM conversations`

	var count int
	err := r.db.Pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count conversations: %w", err)
	}

	return count, nil
}

// CountByStatus returns number of conversations by status
func (r *ConversationRepositoryImpl) CountByStatus(ctx context.Context, status string) (int, error) {
	query := `SELECT COUNT(*) FROM conversations WHERE status = $1`

	var count int
	err := r.db.Pool.QueryRow(ctx, query, status).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count conversations by status: %w", err)
	}

	return count, nil
}

// UpdateLastMessageTime updates the last message time for a conversation
func (r *ConversationRepositoryImpl) UpdateLastMessageTime(ctx context.Context, id int) error {
	query := `
		UPDATE conversations
		SET last_message_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to update last message time: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("conversation not found")
	}

	return nil
}

// SetStatus updates the status of a conversation
func (r *ConversationRepositoryImpl) SetStatus(ctx context.Context, id int, status string) error {
	query := `
		UPDATE conversations
		SET status = $2, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to set conversation status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("conversation not found")
	}

	return nil
}
