package repository

import (
	"context"

	"chatbot/internal/models"

	"github.com/jmoiron/sqlx"
)

type MessageRepository struct {
	db *sqlx.DB
}

func NewMessageRepository(db *sqlx.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(ctx context.Context, message *models.Message) error {
	query := `
		INSERT INTO messages (platform, user_id, page_id, content, comment_id, post_id, message_type, processed, response, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		message.Platform, message.UserID, message.PageID, message.Content,
		message.CommentID, message.PostID, message.MessageType, message.Processed, message.Response,
	).Scan(&message.ID, &message.CreatedAt, &message.UpdatedAt)

	return err
}

func (r *MessageRepository) GetAll(ctx context.Context, limit int) ([]*models.Message, error) {
	query := `
		SELECT id, platform, user_id, page_id, content, comment_id, post_id, message_type, processed, response, created_at, updated_at
		FROM messages 
		ORDER BY created_at DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}
		err := rows.Scan(
			&message.ID, &message.Platform, &message.UserID, &message.PageID,
			&message.Content, &message.CommentID, &message.PostID, &message.MessageType,
			&message.Processed, &message.Response, &message.CreatedAt, &message.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, rows.Err()
}

func (r *MessageRepository) GetByPlatform(ctx context.Context, platform string, limit int) ([]*models.Message, error) {
	query := `
		SELECT id, platform, user_id, page_id, content, comment_id, post_id, message_type, processed, response, created_at, updated_at
		FROM messages 
		WHERE platform = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, platform, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}
		err := rows.Scan(
			&message.ID, &message.Platform, &message.UserID, &message.PageID,
			&message.Content, &message.CommentID, &message.PostID, &message.MessageType,
			&message.Processed, &message.Response, &message.CreatedAt, &message.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, rows.Err()
}

func (r *MessageRepository) MarkProcessed(ctx context.Context, id int, response string) error {
	query := `
		UPDATE messages 
		SET processed = true, response = $2, updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id, response)
	return err
}
