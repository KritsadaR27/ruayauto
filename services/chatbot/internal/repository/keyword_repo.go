package repository

import (
	"context"
	"database/sql"
	"fmt"
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
		SELECT id, keyword, response, is_active, priority, match_type, created_at, updated_at,
		       rule_name, hide_after_reply, send_to_inbox, inbox_message, inbox_image
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
			&keyword.RuleName, &keyword.HideAfterReply, &keyword.SendToInbox,
			&keyword.InboxMessage, &keyword.InboxImage,
		)
		if err != nil {
			return nil, err
		}
		// Set created_by to nil since we don't have users table yet
		keyword.CreatedBy = nil
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
		INSERT INTO keywords (keyword, response, is_active, priority, match_type, rule_name, 
		                     hide_after_reply, send_to_inbox, inbox_message, inbox_image, 
		                     created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		keyword.Keyword, keyword.Response, keyword.IsActive,
		keyword.Priority, keyword.MatchType, keyword.RuleName,
		keyword.HideAfterReply, keyword.SendToInbox, keyword.InboxMessage,
		keyword.InboxImage,
	).Scan(&keyword.ID, &keyword.CreatedAt, &keyword.UpdatedAt)

	// Set created_by to nil since we don't have users table yet
	keyword.CreatedBy = nil
	return err
}

func (r *KeywordRepository) Update(ctx context.Context, keyword *models.Keyword) error {
	query := `
		UPDATE keywords 
		SET keyword = $2, response = $3, is_active = $4, priority = $5, match_type = $6,
		    rule_name = $7, hide_after_reply = $8, send_to_inbox = $9, inbox_message = $10,
		    inbox_image = $11, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		keyword.ID, keyword.Keyword, keyword.Response,
		keyword.IsActive, keyword.Priority, keyword.MatchType,
		keyword.RuleName, keyword.HideAfterReply, keyword.SendToInbox,
		keyword.InboxMessage, keyword.InboxImage,
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
		SELECT id, keyword, response, is_active, priority, match_type, created_at, updated_at,
		       rule_name, hide_after_reply, send_to_inbox, inbox_message, inbox_image
		FROM keywords 
		WHERE id = $1
	`

	keyword := &models.Keyword{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&keyword.ID, &keyword.Keyword, &keyword.Response,
		&keyword.IsActive, &keyword.Priority, &keyword.MatchType,
		&keyword.CreatedAt, &keyword.UpdatedAt,
		&keyword.RuleName, &keyword.HideAfterReply, &keyword.SendToInbox,
		&keyword.InboxMessage, &keyword.InboxImage,
	)

	if err != nil {
		return nil, err
	}

	// Set created_by to nil since we don't have users table yet
	keyword.CreatedBy = nil
	return keyword, nil
}

// New methods for enhanced rules system

// GetKeywordWithResponses gets keyword with all its responses using the view
func (r *KeywordRepository) GetKeywordWithResponses(ctx context.Context, keywordID int) (*models.KeywordWithResponses, error) {
	query := `
		SELECT id, keyword, default_message, rule_name, hide_after_reply, send_to_inbox,
		       inbox_message, inbox_image, is_active, priority, created_by, created_at, updated_at,
		       responses, page_ids
		FROM keywords_with_responses
		WHERE id = $1
	`

	kwr := &models.KeywordWithResponses{}
	var responsesJSON, pageIDsJSON string

	err := r.db.QueryRowContext(ctx, query, keywordID).Scan(
		&kwr.ID, &kwr.Keyword, &kwr.DefaultMessage, &kwr.RuleName,
		&kwr.HideAfterReply, &kwr.SendToInbox, &kwr.InboxMessage, &kwr.InboxImage,
		&kwr.IsActive, &kwr.Priority, &kwr.CreatedBy, &kwr.CreatedAt, &kwr.UpdatedAt,
		&responsesJSON, &pageIDsJSON,
	)

	if err != nil {
		return nil, err
	}

	// Parse JSON responses and page IDs would go here
	// For now, return the basic structure
	return kwr, nil
}

// CreateKeywordResponse creates a new response for a keyword
func (r *KeywordRepository) CreateKeywordResponse(ctx context.Context, response *models.KeywordResponse) error {
	query := `
		INSERT INTO keyword_responses (keyword_id, response_text, response_type, media_url, weight, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		response.KeywordID, response.ResponseText, response.ResponseType,
		response.MediaURL, response.Weight, response.IsActive,
	).Scan(&response.ID, &response.CreatedAt)

	return err
}

// GetKeywordResponses gets all responses for a keyword
func (r *KeywordRepository) GetKeywordResponses(ctx context.Context, keywordID int) ([]*models.KeywordResponse, error) {
	query := `
		SELECT id, keyword_id, response_text, response_type, media_url, weight, is_active, created_at
		FROM keyword_responses
		WHERE keyword_id = $1 AND is_active = true
		ORDER BY weight DESC, created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, keywordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responses []*models.KeywordResponse
	for rows.Next() {
		response := &models.KeywordResponse{}
		err := rows.Scan(
			&response.ID, &response.KeywordID, &response.ResponseText,
			&response.ResponseType, &response.MediaURL, &response.Weight,
			&response.IsActive, &response.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}
	return responses, rows.Err()
}

// GetFallbackResponses gets fallback responses for a page
func (r *KeywordRepository) GetFallbackResponses(ctx context.Context, pageID *int) ([]*models.FallbackResponse, error) {
	var query string
	var args []interface{}

	if pageID != nil {
		query = `
			SELECT id, page_id, response_text, response_type, media_url, weight, is_active, created_at
			FROM fallback_responses
			WHERE (page_id = $1 OR page_id IS NULL) AND is_active = true
			ORDER BY weight DESC, created_at ASC
		`
		args = append(args, *pageID)
	} else {
		query = `
			SELECT id, page_id, response_text, response_type, media_url, weight, is_active, created_at
			FROM fallback_responses
			WHERE page_id IS NULL AND is_active = true
			ORDER BY weight DESC, created_at ASC
		`
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responses []*models.FallbackResponse
	for rows.Next() {
		response := &models.FallbackResponse{}
		err := rows.Scan(
			&response.ID, &response.PageID, &response.ResponseText,
			&response.ResponseType, &response.MediaURL, &response.Weight,
			&response.IsActive, &response.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}
	return responses, rows.Err()
}

// CheckRateLimit checks if user is rate limited for a keyword
func (r *KeywordRepository) CheckRateLimit(ctx context.Context, pageID int, userID string, keywordID int, limitMinutes int) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM rate_limits 
		WHERE page_id = $1 AND user_id = $2 AND keyword_id = $3
		  AND last_response_at > NOW() - INTERVAL '%d minutes'
	`

	var count int
	formattedQuery := fmt.Sprintf(query, limitMinutes)
	err := r.db.QueryRowContext(ctx, formattedQuery, pageID, userID, keywordID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UpdateRateLimit updates or creates rate limit record
func (r *KeywordRepository) UpdateRateLimit(ctx context.Context, rateLimit *models.RateLimit) error {
	query := `
		INSERT INTO rate_limits (page_id, user_id, keyword_id, last_response_at, response_count, created_at)
		VALUES ($1, $2, $3, NOW(), $4, NOW())
		ON CONFLICT (page_id, user_id, keyword_id)
		DO UPDATE SET 
			last_response_at = NOW(),
			response_count = rate_limits.response_count + 1
		RETURNING id, last_response_at, response_count, created_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		rateLimit.PageID, rateLimit.UserID, rateLimit.KeywordID, rateLimit.ResponseCount,
	).Scan(&rateLimit.ID, &rateLimit.LastResponseAt, &rateLimit.ResponseCount, &rateLimit.CreatedAt)

	return err
}

// LogIncomingMessage logs incoming message for analytics
func (r *KeywordRepository) LogIncomingMessage(ctx context.Context, msg *models.IncomingMessage) error {
	query := `
		INSERT INTO incoming_messages (page_id, user_id, message_text, comment_id, post_id, 
		                              matched_keyword_id, response_sent, hidden_comment, 
		                              sent_to_inbox, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		msg.PageID, msg.UserID, msg.MessageText, msg.CommentID, msg.PostID,
		msg.MatchedKeywordID, msg.ResponseSent, msg.HiddenComment, msg.SentToInbox,
	).Scan(&msg.ID, &msg.CreatedAt)

	return err
}

// UpdateRuleAnalytics updates analytics for a rule
func (r *KeywordRepository) UpdateRuleAnalytics(ctx context.Context, keywordID, pageID int, triggered, responded bool) error {
	query := `
		INSERT INTO rule_analytics (keyword_id, page_id, trigger_count, response_count, last_triggered_at, date, created_at)
		VALUES ($1, $2, $3, $4, NOW(), CURRENT_DATE, NOW())
		ON CONFLICT (keyword_id, page_id, date)
		DO UPDATE SET
			trigger_count = rule_analytics.trigger_count + $3,
			response_count = rule_analytics.response_count + $4,
			last_triggered_at = CASE WHEN $3 > 0 THEN NOW() ELSE rule_analytics.last_triggered_at END
	`

	triggerCount := 0
	responseCount := 0
	if triggered {
		triggerCount = 1
	}
	if responded {
		responseCount = 1
	}

	_, err := r.db.ExecContext(ctx, query, keywordID, pageID, triggerCount, responseCount)
	return err
}

// UpdateIncomingMessage updates an incoming message with processing results
func (r *KeywordRepository) UpdateIncomingMessage(ctx context.Context, msg *models.IncomingMessage) error {
	query := `
		UPDATE incoming_messages 
		SET matched_keyword_id = $2, response_sent = $3, hidden_comment = $4, sent_to_inbox = $5
		WHERE id = $1
	`

	_, err := r.db.ExecContext(
		ctx, query,
		msg.ID, msg.MatchedKeywordID, msg.ResponseSent, msg.HiddenComment, msg.SentToInbox,
	)

	return err
}
