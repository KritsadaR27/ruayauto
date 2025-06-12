package repository

import (
	"context"
	"database/sql"
	"fmt"

	"ruaymanagement/backend/internal/ChatBotCore/models"
)

type messageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, message *models.Message) error {
	query := `
		INSERT INTO chatbot_mvp.messages (
			conversation_id, facebook_message_id, sender_type, message_type, 
			message_text, message_data, matched_keyword, response_template, processing_time_ms, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		message.ConversationID,
		message.FacebookMessageID,
		message.SenderType,
		message.MessageType,
		message.MessageText,
		message.MessageData,
		message.MatchedKeyword,
		message.ResponseTemplate,
		message.ProcessingTimeMs,
	).Scan(&message.ID, &message.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	return nil
}

func (r *messageRepository) GetByID(ctx context.Context, id int) (*models.Message, error) {
	message := &models.Message{}
	query := `
		SELECT id, conversation_id, facebook_message_id, sender_type, message_type,
			   message_text, message_data, matched_keyword, response_template, processing_time_ms, created_at
		FROM chatbot_mvp.messages
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&message.ID,
		&message.ConversationID,
		&message.FacebookMessageID,
		&message.SenderType,
		&message.MessageType,
		&message.MessageText,
		&message.MessageData,
		&message.MatchedKeyword,
		&message.ResponseTemplate,
		&message.ProcessingTimeMs,
		&message.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return message, nil
}

func (r *messageRepository) GetByConversationID(ctx context.Context, conversationID int, offset, limit int) ([]models.Message, error) {
	query := `
		SELECT id, conversation_id, facebook_message_id, sender_type, message_type,
			   message_text, message_data, matched_keyword, response_template, processing_time_ms, created_at
		FROM chatbot_mvp.messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, conversationID, limit, offset)
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
			&message.FacebookMessageID,
			&message.SenderType,
			&message.MessageType,
			&message.MessageText,
			&message.MessageData,
			&message.MatchedKeyword,
			&message.ResponseTemplate,
			&message.ProcessingTimeMs,
			&message.CreatedAt,
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
		SELECT id, conversation_id, facebook_message_id, sender_type, message_type,
			   message_text, message_data, matched_keyword, response_template, processing_time_ms, created_at
		FROM chatbot_mvp.messages
		WHERE conversation_id = $1 AND sender_type = $2
		ORDER BY created_at ASC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, query, conversationID, senderType, limit, offset)
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
			&message.FacebookMessageID,
			&message.SenderType,
			&message.MessageType,
			&message.MessageText,
			&message.MessageData,
			&message.MatchedKeyword,
			&message.ResponseTemplate,
			&message.ProcessingTimeMs,
			&message.CreatedAt,
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
	query := `SELECT COUNT(*) FROM chatbot_mvp.messages WHERE conversation_id = $1`

	err := r.db.QueryRowContext(ctx, query, conversationID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count messages: %w", err)
	}

	return count, nil
}

func (r *messageRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM chatbot_mvp.messages WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}
