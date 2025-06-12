package repository

import (
	"context"
	"database/sql"
	"fmt"

	"ruaymanagement/backend/internal/ChatBotCore/models"
)

// ConversationRepositoryImpl implements ConversationRepository
type ConversationRepositoryImpl struct {
	db *sql.DB
}

// NewConversationRepository creates a new conversation repository
func NewConversationRepository(db *sql.DB) ConversationRepository {
	return &ConversationRepositoryImpl{db: db}
}

// Create creates a new conversation
func (r *ConversationRepositoryImpl) Create(ctx context.Context, conv *models.Conversation) error {
	query := `
		INSERT INTO chatbot_mvp.conversations (facebook_user_id, user_name, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query, conv.FacebookUserID, conv.UserName, conv.Status).
		Scan(&conv.ID, &conv.CreatedAt, &conv.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}

	return nil
}

// GetByID retrieves a conversation by ID
func (r *ConversationRepositoryImpl) GetByID(ctx context.Context, id int) (*models.Conversation, error) {
	query := `
		SELECT id, facebook_user_id, user_name, status, last_message_at, created_at, updated_at
		FROM chatbot_mvp.conversations
		WHERE id = $1
	`

	var conv models.Conversation
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&conv.ID,
		&conv.FacebookUserID,
		&conv.UserName,
		&conv.Status,
		&conv.LastMessageAt,
		&conv.CreatedAt,
		&conv.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation by ID: %w", err)
	}

	return &conv, nil
}

// GetByFacebookUserID retrieves a conversation by Facebook user ID
func (r *ConversationRepositoryImpl) GetByFacebookUserID(ctx context.Context, facebookUserID string) (*models.Conversation, error) {
	query := `
		SELECT id, facebook_user_id, user_name, status, last_message_at, created_at, updated_at
		FROM chatbot_mvp.conversations
		WHERE facebook_user_id = $1
	`

	var conv models.Conversation
	err := r.db.QueryRowContext(ctx, query, facebookUserID).Scan(
		&conv.ID,
		&conv.FacebookUserID,
		&conv.UserName,
		&conv.Status,
		&conv.LastMessageAt,
		&conv.CreatedAt,
		&conv.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation by Facebook user ID: %w", err)
	}

	return &conv, nil
}

// Update updates an existing conversation
func (r *ConversationRepositoryImpl) Update(ctx context.Context, conv *models.Conversation) error {
	query := `
		UPDATE chatbot_mvp.conversations
		SET user_name = $2, status = $3, last_message_at = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(ctx, query, conv.ID, conv.UserName, conv.Status, conv.LastMessageAt).
		Scan(&conv.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update conversation: %w", err)
	}

	return nil
}

// Delete deletes a conversation
func (r *ConversationRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM chatbot_mvp.conversations WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("conversation not found")
	}

	return nil
}

// List retrieves conversations with pagination
func (r *ConversationRepositoryImpl) List(ctx context.Context, offset, limit int) ([]models.Conversation, error) {
	query := `
		SELECT id, facebook_user_id, user_name, status, last_message_at, created_at, updated_at
		FROM chatbot_mvp.conversations
		ORDER BY last_message_at DESC NULLS LAST, created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversations: %w", err)
	}
	defer rows.Close()

	var conversations []models.Conversation
	for rows.Next() {
		var conv models.Conversation
		err := rows.Scan(
			&conv.ID,
			&conv.FacebookUserID,
			&conv.UserName,
			&conv.Status,
			&conv.LastMessageAt,
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
		SELECT id, facebook_user_id, user_name, status, last_message_at, created_at, updated_at
		FROM chatbot_mvp.conversations
		WHERE status = $1
		ORDER BY last_message_at DESC NULLS LAST, created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversations by status: %w", err)
	}
	defer rows.Close()

	var conversations []models.Conversation
	for rows.Next() {
		var conv models.Conversation
		err := rows.Scan(
			&conv.ID,
			&conv.FacebookUserID,
			&conv.UserName,
			&conv.Status,
			&conv.LastMessageAt,
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
	query := `SELECT COUNT(*) FROM chatbot_mvp.conversations`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count conversations: %w", err)
	}

	return count, nil
}

// CountByStatus returns number of conversations by status
func (r *ConversationRepositoryImpl) CountByStatus(ctx context.Context, status string) (int, error) {
	query := `SELECT COUNT(*) FROM chatbot_mvp.conversations WHERE status = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, status).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count conversations by status: %w", err)
	}

	return count, nil
}

// UpdateLastMessageTime updates the last message time for a conversation
func (r *ConversationRepositoryImpl) UpdateLastMessageTime(ctx context.Context, id int) error {
	query := `
		UPDATE chatbot_mvp.conversations
		SET last_message_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to update last message time: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("conversation not found")
	}

	return nil
}

// SetStatus updates the status of a conversation
func (r *ConversationRepositoryImpl) SetStatus(ctx context.Context, id int, status string) error {
	query := `
		UPDATE chatbot_mvp.conversations
		SET status = $2, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to set conversation status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("conversation not found")
	}

	return nil
}
