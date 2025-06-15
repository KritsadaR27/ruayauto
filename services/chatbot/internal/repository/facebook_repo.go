package repository

import (
	"database/sql"

	"chatbot/internal/models"

	"github.com/jmoiron/sqlx"
)

type FacebookRepository interface {
	// Session management
	CreateSession(session *models.FacebookSession) error
	GetSession(facebookUserID string) (*models.FacebookSession, error)
	UpdateSession(session *models.FacebookSession) error
	DeleteSession(facebookUserID string) error

	// Page management
	CreatePage(page *models.FacebookPage) error
	GetPage(pageID string) (*models.FacebookPage, error)
	GetPagesByUser(facebookUserID string) ([]models.FacebookPage, error)
	UpdatePage(page *models.FacebookPage) error
	DeletePage(pageID string) error
	SetPageConnectionStatus(facebookPageID string, connected bool) error

	// Webhook subscriptions
	CreateWebhookSubscription(subscription *models.FacebookWebhookSubscription) error
	GetWebhookSubscription(facebookPageID string) (*models.FacebookWebhookSubscription, error)
	UpdateWebhookSubscription(subscription *models.FacebookWebhookSubscription) error
	DeleteWebhookSubscription(facebookPageID string) error
}

type facebookRepository struct {
	db *sqlx.DB
}

func NewFacebookRepository(db *sqlx.DB) FacebookRepository {
	return &facebookRepository{db: db}
}

// Session management
func (r *facebookRepository) CreateSession(session *models.FacebookSession) error {
	query := `
		INSERT INTO facebook_sessions (facebook_user_id, access_token, expires_at, refresh_token)
		VALUES (:facebook_user_id, :access_token, :expires_at, :refresh_token)
	`
	_, err := r.db.NamedExec(query, session)
	return err
}

func (r *facebookRepository) GetSession(facebookUserID string) (*models.FacebookSession, error) {
	var session models.FacebookSession
	query := `
		SELECT id, facebook_user_id, access_token, expires_at, refresh_token, created_at, updated_at
		FROM facebook_sessions 
		WHERE facebook_user_id = $1
	`
	err := r.db.Get(&session, query, facebookUserID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &session, err
}

func (r *facebookRepository) UpdateSession(session *models.FacebookSession) error {
	query := `
		UPDATE facebook_sessions 
		SET access_token = :access_token, expires_at = :expires_at, refresh_token = :refresh_token,
		    updated_at = CURRENT_TIMESTAMP
		WHERE facebook_user_id = :facebook_user_id
	`
	_, err := r.db.NamedExec(query, session)
	return err
}

func (r *facebookRepository) DeleteSession(facebookUserID string) error {
	query := `DELETE FROM facebook_sessions WHERE facebook_user_id = $1`
	_, err := r.db.Exec(query, facebookUserID)
	return err
}

// Page management
func (r *facebookRepository) CreatePage(page *models.FacebookPage) error {
	query := `
		INSERT INTO facebook_pages (page_id, facebook_page_id, page_name, name, access_token, expires_at, 
		                           facebook_user_id, is_active, connected, enabled, webhook_verified)
		VALUES (:page_id, :facebook_page_id, :page_name, :name, :access_token, :expires_at,
		        :facebook_user_id, :is_active, :connected, :enabled, :webhook_verified)
		ON CONFLICT (page_id) DO UPDATE SET
		    facebook_page_id = EXCLUDED.facebook_page_id,
		    page_name = EXCLUDED.page_name,
		    name = EXCLUDED.name,
		    access_token = EXCLUDED.access_token,
		    expires_at = EXCLUDED.expires_at,
		    facebook_user_id = EXCLUDED.facebook_user_id,
		    is_active = EXCLUDED.is_active,
		    connected = EXCLUDED.connected,
		    enabled = EXCLUDED.enabled,
		    webhook_verified = EXCLUDED.webhook_verified,
		    updated_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.NamedExec(query, page)
	return err
}

func (r *facebookRepository) GetPage(pageID string) (*models.FacebookPage, error) {
	var page models.FacebookPage
	query := `
		SELECT id, page_id, facebook_page_id, page_name, name, access_token, expires_at,
		       facebook_user_id, is_active, connected, enabled, webhook_verified, created_at, updated_at
		FROM facebook_pages 
		WHERE page_id = $1
	`
	err := r.db.Get(&page, query, pageID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &page, err
}

func (r *facebookRepository) GetPagesByUser(facebookUserID string) ([]models.FacebookPage, error) {
	var pages []models.FacebookPage
	query := `
		SELECT id, page_id, facebook_page_id, page_name, name, access_token, expires_at,
		       facebook_user_id, is_active, connected, enabled, webhook_verified, created_at, updated_at
		FROM facebook_pages 
		WHERE facebook_user_id = $1
		ORDER BY created_at DESC
	`
	err := r.db.Select(&pages, query, facebookUserID)
	return pages, err
}

func (r *facebookRepository) UpdatePage(page *models.FacebookPage) error {
	query := `
		UPDATE facebook_pages 
		SET facebook_page_id = :facebook_page_id, page_name = :page_name, name = :name,
		    access_token = :access_token, expires_at = :expires_at, facebook_user_id = :facebook_user_id,
		    is_active = :is_active, connected = :connected, enabled = :enabled,
		    webhook_verified = :webhook_verified, updated_at = CURRENT_TIMESTAMP
		WHERE page_id = :page_id
	`
	_, err := r.db.NamedExec(query, page)
	return err
}

func (r *facebookRepository) DeletePage(pageID string) error {
	query := `DELETE FROM facebook_pages WHERE page_id = $1`
	_, err := r.db.Exec(query, pageID)
	return err
}

func (r *facebookRepository) SetPageConnectionStatus(facebookPageID string, connected bool) error {
	query := `
		UPDATE facebook_pages 
		SET connected = $1, updated_at = CURRENT_TIMESTAMP
		WHERE facebook_page_id = $2 OR page_id = $2
	`
	_, err := r.db.Exec(query, connected, facebookPageID)
	return err
}

// Webhook subscriptions
func (r *facebookRepository) CreateWebhookSubscription(subscription *models.FacebookWebhookSubscription) error {
	query := `
		INSERT INTO facebook_webhook_subscriptions (facebook_page_id, subscription_id, webhook_url, 
		                                          verify_token, fields, active)
		VALUES (:facebook_page_id, :subscription_id, :webhook_url, :verify_token, :fields, :active)
		ON CONFLICT (facebook_page_id) DO UPDATE SET
		    subscription_id = EXCLUDED.subscription_id,
		    webhook_url = EXCLUDED.webhook_url,
		    verify_token = EXCLUDED.verify_token,
		    fields = EXCLUDED.fields,
		    active = EXCLUDED.active,
		    updated_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.NamedExec(query, subscription)
	return err
}

func (r *facebookRepository) GetWebhookSubscription(facebookPageID string) (*models.FacebookWebhookSubscription, error) {
	var subscription models.FacebookWebhookSubscription
	query := `
		SELECT id, facebook_page_id, subscription_id, webhook_url, verify_token, fields,
		       active, last_verified_at, created_at, updated_at
		FROM facebook_webhook_subscriptions 
		WHERE facebook_page_id = $1
	`
	err := r.db.Get(&subscription, query, facebookPageID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &subscription, err
}

func (r *facebookRepository) UpdateWebhookSubscription(subscription *models.FacebookWebhookSubscription) error {
	query := `
		UPDATE facebook_webhook_subscriptions 
		SET subscription_id = :subscription_id, webhook_url = :webhook_url, verify_token = :verify_token,
		    fields = :fields, active = :active, last_verified_at = :last_verified_at,
		    updated_at = CURRENT_TIMESTAMP
		WHERE facebook_page_id = :facebook_page_id
	`
	_, err := r.db.NamedExec(query, subscription)
	return err
}

func (r *facebookRepository) DeleteWebhookSubscription(facebookPageID string) error {
	query := `DELETE FROM facebook_webhook_subscriptions WHERE facebook_page_id = $1`
	_, err := r.db.Exec(query, facebookPageID)
	return err
}
