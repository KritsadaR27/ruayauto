package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"ruaymanagement/backend/internal/ChatBotCore/database"
	"ruaymanagement/backend/internal/ChatBotCore/models"
)

type messageRepository struct {
	db *database.DB
}

func NewMessageRepository(db *database.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, message *models.Message) error {
	query := `
		INSERT INTO messages (
			conversation_id, message_id, sender_type, message_type, 
			content, metadata, created_at, processed_at, response_time_ms
		) VALUES ($1, $2, $3, $4, $5, $6, NOW(), $7, $8)
		RETURNING id, created_at
	`

	err := r.db.Pool.QueryRow(
		ctx, query,
		message.ConversationID,
		message.MessageID,
		message.SenderType,
		message.MessageType,
		message.Content,
		message.Metadata,
		message.ProcessedAt,
		message.ResponseTimeMs,
	).Scan(&message.ID, &message.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	return nil
}

func (r *messageRepository) GetByID(ctx context.Context, id int) (*models.Message, error) {
	message := &models.Message{}
	query := `
		SELECT id, conversation_id, message_id, sender_type, message_type,
			   content, metadata, created_at, processed_at, response_time_ms
		FROM messages
		WHERE id = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&message.ID,
		&message.ConversationID,
		&message.MessageID,
		&message.SenderType,
		&message.MessageType,
		&message.Content,
		&message.Metadata,
		&message.CreatedAt,
		&message.ProcessedAt,
		&message.ResponseTimeMs,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return message, nil
}

func (r *messageRepository) GetByConversationID(ctx context.Context, conversationID int, offset, limit int) ([]models.Message, error) {
	query := `
		SELECT id, conversation_id, message_id, sender_type, message_type,
			   content, metadata, created_at, processed_at, response_time_ms
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, query, conversationID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages by conversation: %w", err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		message := models.Message{}
		err := rows.Scan(
			&message.ID,
			&message.ConversationID,
			&message.MessageID,
			&message.SenderType,
			&message.MessageType,
			&message.Content,
			&message.Metadata,
			&message.CreatedAt,
			&message.ProcessedAt,
			&message.ResponseTimeMs,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating messages: %w", err)
	}

	return messages, nil
}

func (r *messageRepository) GetBySenderType(ctx context.Context, conversationID int, senderType string, offset, limit int) ([]models.Message, error) {
	query := `
		SELECT id, conversation_id, message_id, sender_type, message_type,
			   content, metadata, created_at, processed_at, response_time_ms
		FROM messages
		WHERE conversation_id = $1 AND sender_type = $2
		ORDER BY created_at ASC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Pool.Query(ctx, query, conversationID, senderType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages by sender type: %w", err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		message := models.Message{}
		err := rows.Scan(
			&message.ID,
			&message.ConversationID,
			&message.MessageID,
			&message.SenderType,
			&message.MessageType,
			&message.Content,
			&message.Metadata,
			&message.CreatedAt,
			&message.ProcessedAt,
			&message.ResponseTimeMs,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating messages: %w", err)
	}

	return messages, nil
}

func (r *messageRepository) CountByConversationID(ctx context.Context, conversationID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM messages WHERE conversation_id = $1`

	err := r.db.Pool.QueryRow(ctx, query, conversationID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count messages: %w", err)
	}

	return count, nil
}

func (r *messageRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM messages WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}
