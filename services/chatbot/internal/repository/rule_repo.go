// ...existing code from keyword_repo.go, replacing all 'Keyword' with 'Rule', 'keyword' with 'rule', 'keywords' with 'rules',
// and updating all SQL queries and struct/method names accordingly...

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"chatbot/internal/models"
)

type RuleRepository struct {
	db *sql.DB
}

func NewRuleRepository(db *sql.DB) *RuleRepository {
	return &RuleRepository{db: db}
}

// Helper: fetch keywords for a rule
func (r *RuleRepository) getKeywordsForRule(ctx context.Context, ruleID int) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT keyword FROM rule_keywords WHERE rule_id = $1", ruleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var keywords []string
	for rows.Next() {
		var k string
		if err := rows.Scan(&k); err == nil {
			keywords = append(keywords, k)
		}
	}
	return keywords, nil
}

// Helper: insert keywords for a rule
func (r *RuleRepository) insertKeywords(ctx context.Context, ruleID int, keywords []string) error {
	for _, k := range keywords {
		_, err := r.db.ExecContext(ctx, "INSERT INTO rule_keywords (rule_id, keyword) VALUES ($1, $2)", ruleID, strings.TrimSpace(k))
		if err != nil {
			return err
		}
	}
	return nil
}

// Helper: delete all keywords for a rule
func (r *RuleRepository) deleteKeywords(ctx context.Context, ruleID int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM rule_keywords WHERE rule_id = $1", ruleID)
	return err
}

func (r *RuleRepository) GetActive(ctx context.Context) ([]*models.Rule, error) {
	query := `
		SELECT id, response, is_active, priority, match_type, created_at, updated_at,
		       rule_name, hide_after_reply, send_to_inbox, inbox_message, inbox_image
		FROM rules 
		WHERE is_active = true
		ORDER BY priority DESC, created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*models.Rule
	for rows.Next() {
		rule := &models.Rule{}
		err := rows.Scan(
			&rule.ID, &rule.Response,
			&rule.IsActive, &rule.Priority, &rule.MatchType,
			&rule.CreatedAt, &rule.UpdatedAt,
			&rule.RuleName, &rule.HideAfterReply, &rule.SendToInbox,
			&rule.InboxMessage, &rule.InboxImage,
		)
		if err != nil {
			return nil, err
		}
		rule.CreatedBy = nil
		// fetch keywords
		keywords, _ := r.getKeywordsForRule(ctx, rule.ID)
		rule.Keywords = keywords
		rules = append(rules, rule)
	}
	return rules, rows.Err()
}

func (r *RuleRepository) Create(ctx context.Context, rule *models.Rule) error {
	query := `
		INSERT INTO rules (response, is_active, priority, match_type, rule_name, 
		                     hide_after_reply, send_to_inbox, inbox_message, inbox_image, 
		                     created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		rule.Response, rule.IsActive,
		rule.Priority, rule.MatchType, rule.RuleName,
		rule.HideAfterReply, rule.SendToInbox, rule.InboxMessage,
		rule.InboxImage,
	).Scan(&rule.ID, &rule.CreatedAt, &rule.UpdatedAt)
	rule.CreatedBy = nil
	if err != nil {
		return err
	}
	// insert keywords
	if len(rule.Keywords) > 0 {
		return r.insertKeywords(ctx, rule.ID, rule.Keywords)
	}
	return nil
}

func (r *RuleRepository) Update(ctx context.Context, rule *models.Rule) error {
	query := `
		UPDATE rules 
		SET response = $2, is_active = $3, priority = $4, match_type = $5,
		    rule_name = $6, hide_after_reply = $7, send_to_inbox = $8, inbox_message = $9,
		    inbox_image = $10, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		rule.ID, rule.Response,
		rule.IsActive, rule.Priority, rule.MatchType,
		rule.RuleName, rule.HideAfterReply, rule.SendToInbox,
		rule.InboxMessage, rule.InboxImage,
	).Scan(&rule.UpdatedAt)
	if err != nil {
		return err
	}
	// update keywords
	_ = r.deleteKeywords(ctx, rule.ID)
	if len(rule.Keywords) > 0 {
		return r.insertKeywords(ctx, rule.ID, rule.Keywords)
	}
	return nil
}

func (r *RuleRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM rules WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *RuleRepository) GetByID(ctx context.Context, id int) (*models.Rule, error) {
	query := `
		SELECT id, response, is_active, priority, match_type, created_at, updated_at,
		       rule_name, hide_after_reply, send_to_inbox, inbox_message, inbox_image
		FROM rules 
		WHERE id = $1
	`
	rule := &models.Rule{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rule.ID, &rule.Response,
		&rule.IsActive, &rule.Priority, &rule.MatchType,
		&rule.CreatedAt, &rule.UpdatedAt,
		&rule.RuleName, &rule.HideAfterReply, &rule.SendToInbox,
		&rule.InboxMessage, &rule.InboxImage,
	)
	if err != nil {
		return nil, err
	}
	rule.CreatedBy = nil
	keywords, _ := r.getKeywordsForRule(ctx, rule.ID)
	rule.Keywords = keywords
	return rule, nil
}

// FindMatching finds rules that match the given content
func (r *RuleRepository) FindMatching(ctx context.Context, content string) ([]*models.Rule, error) {
	rules, err := r.GetActive(ctx)
	if err != nil {
		return nil, err
	}

	content = strings.ToLower(strings.TrimSpace(content))
	var matching []*models.Rule

	for _, rule := range rules {
	keywordLoop:
		for _, keyword := range rule.Keywords {
			keywordText := strings.ToLower(strings.TrimSpace(keyword))
			switch rule.MatchType {
			case "exact":
				if content == keywordText {
					matching = append(matching, rule)
					break keywordLoop
				}
			case "contains":
				if strings.Contains(content, keywordText) {
					matching = append(matching, rule)
					break keywordLoop
				}
			case "starts_with":
				if strings.HasPrefix(content, keywordText) {
					matching = append(matching, rule)
					break keywordLoop
				}
			case "ends_with":
				if strings.HasSuffix(content, keywordText) {
					matching = append(matching, rule)
					break keywordLoop
				}
			default:
				if strings.Contains(content, keywordText) {
					matching = append(matching, rule)
					break keywordLoop
				}
			}
		}
	}
	return matching, nil
}

// GetRuleResponses gets all responses for a rule
func (r *RuleRepository) GetRuleResponses(ctx context.Context, ruleID int) ([]*models.RuleResponse, error) {
	query := `
		SELECT id, rule_id, response_text, response_type, media_url, weight, is_active, created_at
		FROM rule_responses
		WHERE rule_id = $1 AND is_active = true
		ORDER BY weight DESC, created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, ruleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responses []*models.RuleResponse
	for rows.Next() {
		response := &models.RuleResponse{}
		err := rows.Scan(
			&response.ID, &response.RuleID, &response.ResponseText,
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

// CreateRuleResponse creates a new response for a rule
func (r *RuleRepository) CreateRuleResponse(ctx context.Context, response *models.RuleResponse) error {
	query := `
		INSERT INTO rule_responses (rule_id, response_text, response_type, media_url, weight, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		response.RuleID, response.ResponseText, response.ResponseType,
		response.MediaURL, response.Weight, response.IsActive,
	).Scan(&response.ID, &response.CreatedAt)
	return err
}

// UpdateRuleResponse updates a rule response
func (r *RuleRepository) UpdateRuleResponse(ctx context.Context, response *models.RuleResponse) error {
	query := `
		UPDATE rule_responses 
		SET response_text = $2, response_type = $3, media_url = $4, weight = $5, is_active = $6
		WHERE id = $1
	`
	_, err := r.db.ExecContext(
		ctx, query,
		response.ID, response.ResponseText, response.ResponseType,
		response.MediaURL, response.Weight, response.IsActive,
	)
	return err
}

// DeleteRuleResponse deletes a rule response
func (r *RuleRepository) DeleteRuleResponse(ctx context.Context, responseID int) error {
	query := `DELETE FROM rule_responses WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, responseID)
	return err
}

// GetRuleResponseByID gets a single rule response by ID
func (r *RuleRepository) GetRuleResponseByID(ctx context.Context, responseID int) (*models.RuleResponse, error) {
	query := `
		SELECT id, rule_id, response_text, response_type, media_url, weight, is_active, created_at
		FROM rule_responses
		WHERE id = $1
	`
	response := &models.RuleResponse{}
	err := r.db.QueryRowContext(ctx, query, responseID).Scan(
		&response.ID, &response.RuleID, &response.ResponseText,
		&response.ResponseType, &response.MediaURL, &response.Weight,
		&response.IsActive, &response.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// GetFallbackResponses gets fallback responses for a page
func (r *RuleRepository) GetFallbackResponses(ctx context.Context, pageID *int) ([]*models.FallbackResponse, error) {
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

// CheckRateLimit checks if user is rate limited for a rule
func (r *RuleRepository) CheckRateLimit(ctx context.Context, pageID int, userID string, ruleID int, limitMinutes int) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM rate_limits 
		WHERE page_id = $1 AND user_id = $2 AND rule_id = $3
		  AND last_response_at > NOW() - INTERVAL '%d minutes'
	`
	var count int
	formattedQuery := fmt.Sprintf(query, limitMinutes)
	err := r.db.QueryRowContext(ctx, formattedQuery, pageID, userID, ruleID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// UpdateRateLimit updates or creates rate limit record
func (r *RuleRepository) UpdateRateLimit(ctx context.Context, rateLimit *models.RateLimit) error {
	query := `
		INSERT INTO rate_limits (page_id, user_id, rule_id, last_response_at, response_count, created_at)
		VALUES ($1, $2, $3, NOW(), $4, NOW())
		ON CONFLICT (page_id, user_id, rule_id)
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
func (r *RuleRepository) LogIncomingMessage(ctx context.Context, msg *models.IncomingMessage) error {
	query := `
		INSERT INTO incoming_messages (page_id, user_id, message_text, comment_id, post_id, 
										matched_rule_id, response_sent, hidden_comment, 
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
func (r *RuleRepository) UpdateRuleAnalytics(ctx context.Context, ruleID, pageID int, triggered, responded bool) error {
	query := `
		INSERT INTO rule_analytics (rule_id, page_id, trigger_count, response_count, last_triggered_at, date, created_at)
		VALUES ($1, $2, $3, $4, NOW(), CURRENT_DATE, NOW())
		ON CONFLICT (rule_id, page_id, date)
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
	_, err := r.db.ExecContext(ctx, query, ruleID, pageID, triggerCount, responseCount)
	return err
}

// UpdateIncomingMessage updates an incoming message with processing results
func (r *RuleRepository) UpdateIncomingMessage(ctx context.Context, msg *models.IncomingMessage) error {
	query := `
		UPDATE incoming_messages 
		SET matched_rule_id = $2, response_sent = $3, hidden_comment = $4, sent_to_inbox = $5
		WHERE id = $1
	`
	_, err := r.db.ExecContext(
		ctx, query,
		msg.ID, msg.MatchedKeywordID, msg.ResponseSent, msg.HiddenComment, msg.SentToInbox,
	)
	return err
}
