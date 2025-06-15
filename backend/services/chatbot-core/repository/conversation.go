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

// FindOrCreate finds an existing conversation or creates a new one
func (r *ConversationRepository) FindOrCreate(ctx context.Context, pageID, userID, commentID, postID string) (*models.Conversation, error) {
	// Try to find existing conversation
	var conv models.Conversation
	query := `
		SELECT id, page_id, user_id, comment_id, post_id, status, last_message, created_at, updated_at
		FROM conversations 
		WHERE page_id = $1 AND user_id = $2 AND status = 'active'
		ORDER BY last_message DESC
		LIMIT 1
	`

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

	if err == nil {
		// Update last message time
		r.UpdateLastMessage(ctx, conv.ID)
		return &conv, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// Create new conversation
	conv = models.Conversation{
		PageID:    pageID,
		UserID:    userID,
		CommentID: commentID,
		PostID:    postID,
		Status:    "active",
	}

	err = r.Create(ctx, &conv)
	if err != nil {
		return nil, err
	}

	return &conv, nil
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

// UpdateLastMessage updates the last message timestamp
func (r *ConversationRepository) UpdateLastMessage(ctx context.Context, conversationID int) error {
	query := `UPDATE conversations SET last_message = $1, updated_at = $2 WHERE id = $3`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, now, now, conversationID)
	return err
}
