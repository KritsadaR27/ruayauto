package repository

import (
	"context"
	"database/sql"
	"time"

	"chatbot-core/models"
)

// ConversationRepository handles conversation-related database operations
type ConversationRepository struct {
	db *sql.DB
}

// NewConversationRepository creates a new conversation repository
func NewConversationRepository(db *sql.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

// Create creates a new conversation
func (r *ConversationRepository) Create(ctx context.Context, conv *models.Conversation) error {
	query := `
		INSERT INTO conversations (page_id, user_id, comment_id, post_id, status, last_message, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	now := time.Now()
	conv.CreatedAt = now
	conv.UpdatedAt = now
	conv.LastMessage = now

	return r.db.QueryRowContext(ctx, query,
		conv.PageID,
		conv.UserID,
		conv.CommentID,
		conv.PostID,
		conv.Status,
		conv.LastMessage,
		conv.CreatedAt,
		conv.UpdatedAt,
	).Scan(&conv.ID)
}

// FindByPageAndUser finds conversation by page and user
func (r *ConversationRepository) FindByPageAndUser(ctx context.Context, pageID, userID string) (*models.Conversation, error) {
	query := `
		SELECT id, page_id, user_id, comment_id, post_id, status, last_message, created_at, updated_at
		FROM conversations 
		WHERE page_id = $1 AND user_id = $2 AND status = 'active'
		ORDER BY last_message DESC
		LIMIT 1
	`

	conv := &models.Conversation{}
	err := r.db.QueryRowContext(ctx, query, pageID, userID).Scan(
		&conv.ID,
		&conv.PageID,
		&conv.UserID,
		&conv.CommentID,
		&conv.PostID,
		&conv.Status,
		&conv.LastMessage,
		&conv.CreatedAt,
		&conv.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return conv, nil
}
