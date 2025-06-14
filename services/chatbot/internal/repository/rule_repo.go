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

func (r *RuleRepository) GetActive(ctx context.Context) ([]*models.Rule, error) {
	query := `
		SELECT id, rule, response, is_active, priority, match_type, created_at, updated_at,
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
			&rule.ID, &rule.Rule, &rule.Response,
			&rule.IsActive, &rule.Priority, &rule.MatchType,
			&rule.CreatedAt, &rule.UpdatedAt,
			&rule.RuleName, &rule.HideAfterReply, &rule.SendToInbox,
			&rule.InboxMessage, &rule.InboxImage,
		)
		if err != nil {
			return nil, err
		}
		rule.CreatedBy = nil
		rules = append(rules, rule)
	}
	return rules, rows.Err()
}

func (r *RuleRepository) Create(ctx context.Context, rule *models.Rule) error {
	query := `
		INSERT INTO rules (rule, response, is_active, priority, match_type, rule_name, 
		                     hide_after_reply, send_to_inbox, inbox_message, inbox_image, 
		                     created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		rule.Rule, rule.Response, rule.IsActive, 
		rule.Priority, rule.MatchType, rule.RuleName,
		rule.HideAfterReply, rule.SendToInbox, rule.InboxMessage,
		rule.InboxImage,
	).Scan(&rule.ID, &rule.CreatedAt, &rule.UpdatedAt)
	rule.CreatedBy = nil
	return err
}

func (r *RuleRepository) Update(ctx context.Context, rule *models.Rule) error {
	query := `
		UPDATE rules 
		SET rule = $2, response = $3, is_active = $4, priority = $5, match_type = $6,
		    rule_name = $7, hide_after_reply = $8, send_to_inbox = $9, inbox_message = $10,
		    inbox_image = $11, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		rule.ID, rule.Rule, rule.Response, 
		rule.IsActive, rule.Priority, rule.MatchType,
		rule.RuleName, rule.HideAfterReply, rule.SendToInbox,
		rule.InboxMessage, rule.InboxImage,
	).Scan(&rule.UpdatedAt)
	return err
}

func (r *RuleRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM rules WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *RuleRepository) GetByID(ctx context.Context, id int) (*models.Rule, error) {
	query := `
		SELECT id, rule, response, is_active, priority, match_type, created_at, updated_at,
		       rule_name, hide_after_reply, send_to_inbox, inbox_message, inbox_image
		FROM rules 
		WHERE id = $1
	`
	rule := &models.Rule{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rule.ID, &rule.Rule, &rule.Response,
		&rule.IsActive, &rule.Priority, &rule.MatchType,
		&rule.CreatedAt, &rule.UpdatedAt,
		&rule.RuleName, &rule.HideAfterReply, &rule.SendToInbox,
		&rule.InboxMessage, &rule.InboxImage,
	)
	if err != nil {
		return nil, err
	}
	rule.CreatedBy = nil
	return rule, nil
}

// ...implement all other methods from keyword_repo.go, replacing 'Keyword' with 'Rule', 'keyword' with 'rule', 'keywords' with 'rules', and updating SQL queries and struct/method names accordingly...
