package models

import (
	"time"
)

// FacebookSession represents a Facebook user session
type FacebookSession struct {
	ID             int        `json:"id" db:"id"`
	FacebookUserID string     `json:"facebook_user_id" db:"facebook_user_id"`
	AccessToken    string     `json:"access_token" db:"access_token"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	RefreshToken   *string    `json:"refresh_token,omitempty" db:"refresh_token"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// FacebookPage represents a Facebook page
type FacebookPage struct {
	ID              int        `json:"id" db:"id"`
	PageID          string     `json:"page_id" db:"page_id"`
	FacebookPageID  *string    `json:"facebook_page_id,omitempty" db:"facebook_page_id"`
	PageName        string     `json:"page_name" db:"page_name"`
	Name            *string    `json:"name,omitempty" db:"name"`
	AccessToken     *string    `json:"access_token,omitempty" db:"access_token"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	FacebookUserID  *string    `json:"facebook_user_id,omitempty" db:"facebook_user_id"`
	IsActive        bool       `json:"is_active" db:"is_active"`
	Connected       *bool      `json:"connected,omitempty" db:"connected"`
	Enabled         *bool      `json:"enabled,omitempty" db:"enabled"`
	WebhookVerified *bool      `json:"webhook_verified,omitempty" db:"webhook_verified"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// FacebookWebhookSubscription represents a webhook subscription for a Facebook page
type FacebookWebhookSubscription struct {
	ID             int        `json:"id" db:"id"`
	FacebookPageID string     `json:"facebook_page_id" db:"facebook_page_id"`
	SubscriptionID *string    `json:"subscription_id,omitempty" db:"subscription_id"`
	WebhookURL     string     `json:"webhook_url" db:"webhook_url"`
	VerifyToken    string     `json:"verify_token" db:"verify_token"`
	Fields         []string   `json:"fields" db:"fields"`
	Active         bool       `json:"active" db:"active"`
	LastVerifiedAt *time.Time `json:"last_verified_at,omitempty" db:"last_verified_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// FacebookOAuthRequest represents an OAuth request
type FacebookOAuthRequest struct {
	RedirectURI string `json:"redirect_uri"`
	State       string `json:"state,omitempty"`
}

// FacebookOAuthResponse represents an OAuth response
type FacebookOAuthResponse struct {
	AuthURL string `json:"auth_url"`
	State   string `json:"state"`
}

// FacebookAuthStatus represents the authentication status
type FacebookAuthStatus struct {
	IsAuthenticated bool           `json:"is_authenticated"`
	UserID          string         `json:"user_id,omitempty"`
	Pages           []FacebookPage `json:"pages,omitempty"`
}

// FacebookPageInfo represents page information from Facebook API
type FacebookPageInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	AccessToken string `json:"access_token"`
	Category    string `json:"category,omitempty"`
}

// FacebookPagesResponse represents the response from Facebook Pages API
type FacebookPagesResponse struct {
	Data []FacebookPageInfo `json:"data"`
}
